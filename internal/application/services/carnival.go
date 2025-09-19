package services

import (
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"haoma/internal/config"
	"haoma/internal/domain/leaderboard"
	"haoma/internal/domain/player"
	"haoma/internal/domain/question"
	"haoma/internal/domain/session"
)

type CarnivalService struct {
	sessionRepo     SessionRepository
	questionRepo    QuestionRepository
	playerRepo      PlayerRepository
	leaderboardRepo LeaderboardRepository
}

type SessionRepository interface {
	Save(s *session.Session) error
	FindByID(id uuid.UUID) (*session.Session, error)
	Update(s *session.Session) error
}

type QuestionRepository interface {
	GetCategories() ([]question.Category, error)
	GetQuestionsByCategory(categoryID uuid.UUID, limit int) ([]question.Question, error)
	GetUnusedFunQuestionsForSession(sessionID uuid.UUID, limit int) ([]question.Question, error)
	FindByID(id uuid.UUID) (*question.Question, error)
}

type PlayerRepository interface {
	Save(p *player.Player) error
	FindByID(id uuid.UUID) (*player.Player, error)
	FindByEmail(email string) (*player.Player, error)
	SaveAttempt(a *player.Attempt) error
	GetAttemptsBySessionAndCategory(sessionID, categoryID uuid.UUID) ([]player.Attempt, error)
	HasAnsweredQuestion(sessionID, questionID uuid.UUID) (bool, error)
}

type LeaderboardRepository interface {
	AddEntry(entry *leaderboard.Entry) error
	UpsertEntry(entry *leaderboard.Entry) error
	GetTop10() ([]leaderboard.Entry, error)
}

func NewCarnivalService(
	sessionRepo SessionRepository,
	questionRepo QuestionRepository,
	playerRepo PlayerRepository,
	leaderboardRepo LeaderboardRepository,
) *CarnivalService {
	return &CarnivalService{
		sessionRepo:     sessionRepo,
		questionRepo:    questionRepo,
		playerRepo:      playerRepo,
		leaderboardRepo: leaderboardRepo,
	}
}

func (c *CarnivalService) CreatePlayer(p *player.Player) error {
	_, err := c.playerRepo.FindByEmail(p.Email)
	if err == nil {
		return errors.New("player already exists")
	}

	return c.playerRepo.Save(p)
}

func (c *CarnivalService) GetPlayerByID(id uuid.UUID) (*player.Player, error) {
	return c.playerRepo.FindByID(id)
}

func (c *CarnivalService) GetPlayerByEmail(email string) (*player.Player, error) {
	return c.playerRepo.FindByEmail(email)
}

func (c *CarnivalService) CreateSession(playerID uuid.UUID) (*session.Session, error) {
	_, err := c.playerRepo.FindByID(playerID)
	if err != nil {
		return nil, errors.New("player not found")
	}

	randomCategories, err := c.generateRandomCategoryAssignment()
	if err != nil {
		return nil, err
	}

	s := &session.Session{
		ID:             uuid.New(),
		PlayerID:       playerID,
		StartedAt:      time.Now(),
		CurrentNode:    config.DEFAULT_NODE_START,
		Score:          session.Score{},
		Categories:     session.StringSlice(randomCategories),
		NodeStartTimes: make(session.IntMap),
	}

	if err := c.sessionRepo.Save(s); err != nil {
		return nil, err
	}

	return s, nil
}

func (c *CarnivalService) ScanNodeQR(playerID uuid.UUID, nodeCode string, sessionID *uuid.UUID) (*uuid.UUID, *question.Node, error) {
	_, err := c.playerRepo.FindByID(playerID)
	if err != nil {
		return nil, nil, errors.New("player not found")
	}

	nodeNumber, err := c.parseNodeCode(nodeCode)
	if err != nil {
		return nil, nil, errors.New("node not found")
	}

	var s *session.Session

	if sessionID != nil {
		s, err = c.sessionRepo.FindByID(*sessionID)
		if err != nil {
			return nil, nil, errors.New("session not found")
		}
		if !s.IsActive() {
			return nil, nil, errors.New("session expired")
		}
	} else {
		s, err = c.CreateSession(playerID)
		if err != nil {
			return nil, nil, err
		}
	}

	if nodeNumber > len([]string(s.Categories)) {
		return nil, nil, errors.New("invalid node number for session")
	}

	assignedCategoryName := ([]string(s.Categories))[nodeNumber-1] // Arrays are 0-indexed, nodes are 1-indexed

	categories, err := c.questionRepo.GetCategories()
	if err != nil {
		return nil, nil, err
	}

	var nodeCategory *question.Category
	for _, cat := range categories {
		if cat.Name == assignedCategoryName {
			nodeCategory = &cat
			break
		}
	}

	if nodeCategory == nil {
		return nil, nil, errors.New("assigned category not found")
	}

	node, err := c.generateNodeFromCategory(nodeNumber, nodeCategory.ID, s.ID)
	if err != nil {
		return nil, nil, err
	}

	if _, exists := s.NodeStartTimes[nodeNumber]; !exists {
		if s.NodeStartTimes == nil {
			s.NodeStartTimes = make(session.IntMap)
		}
		s.NodeStartTimes[nodeNumber] = time.Now().Unix()
	}

	s.CurrentNode = nodeNumber
	if err := c.sessionRepo.Update(s); err != nil {
		return nil, nil, err
	}

	return &s.ID, node, nil
}

type AnswerResult struct {
	IsCorrect               bool
	NodeCompleted           bool
	QuestionsAnsweredInNode int
	CurrentScore            int
}

func (c *CarnivalService) SubmitAnswer(sessionID, questionID uuid.UUID, answer string) (*AnswerResult, error) {
	s, err := c.sessionRepo.FindByID(sessionID)
	if err != nil {
		return nil, err
	}

	if !s.IsActive() {
		return nil, errors.New("session expired or finished")
	}

	q, err := c.questionRepo.FindByID(questionID)
	if err != nil {
		return nil, err
	}

	hasAnswered, err := c.playerRepo.HasAnsweredQuestion(sessionID, questionID)
	if err != nil {
		return nil, err
	}
	if hasAnswered {
		return nil, errors.New("question already answered")
	}

	isCorrect := q.ValidateAnswer(answer)

	attempt := player.NewAttempt(sessionID, questionID, answer, isCorrect)
	if err := c.playerRepo.SaveAttempt(attempt); err != nil {
		return nil, err
	}

	s.Score.Total++
	if isCorrect {
		s.Score.Correct++
	}

	currentNodeCategoryID, err := c.getCurrentNodeCategoryID(sessionID, s.CurrentNode)
	if err != nil {
		return nil, err
	}

	questionsInCurrentNode, err := c.countQuestionsAnsweredInCategory(sessionID, currentNodeCategoryID)
	if err != nil {
		return nil, err
	}

	nodeCompleted := questionsInCurrentNode >= config.QUESTIONS_TO_COMPLETE_NODE

	result := &AnswerResult{
		IsCorrect:               isCorrect,
		NodeCompleted:           nodeCompleted,
		QuestionsAnsweredInNode: questionsInCurrentNode,
	}

	if nodeCompleted {
		if err := c.addNodeTimePenalty(s, s.CurrentNode); err != nil {
			return nil, err
		}

		currentScore := s.CalculateScore()
		result.CurrentScore = currentScore.Final

		if err := c.updateLeaderboardAfterNode(s, result); err != nil {
			return nil, err
		}
	}

	if err := c.sessionRepo.Update(s); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *CarnivalService) GetLeaderboard() ([]leaderboard.Entry, error) {
	return c.leaderboardRepo.GetTop10()
}

func (c *CarnivalService) parseNodeCode(nodeCode string) (int, error) {
	// QR codes format: "NODE_XXX" where XXX is a unique identifier
	// Examples: "NODE_001", "NODE_002", "NODE_003", etc.

	if nodeNumber, exists := config.NodeCodes[nodeCode]; exists {
		return nodeNumber, nil
	}

	return config.DEFAULT_RANK, errors.New("invalid node code")
}

func (c *CarnivalService) generateNodeFromCategory(nodeNumber int, categoryID uuid.UUID, sessionID uuid.UUID) (*question.Node, error) {
	if nodeNumber < config.MIN_NODE_NUMBER || nodeNumber > config.MAX_NODE_NUMBER {
		return nil, errors.New("invalid node number")
	}

	categoryQuestions, err := c.getUniqueQuestionsFromCategory(categoryID, config.CATEGORY_QUESTIONS_PER_NODE)
	if err != nil {
		return nil, err
	}

	funQuestions, err := c.getUnusedFunQuestionsForSession(sessionID, config.FUN_QUESTIONS_PER_NODE)
	if err != nil {
		return nil, err
	}

	allQuestions := append(categoryQuestions, funQuestions...)

	rand.Shuffle(len(allQuestions), func(i, j int) {
		allQuestions[i], allQuestions[j] = allQuestions[j], allQuestions[i]
	})

	return &question.Node{
		Number:     nodeNumber,
		CategoryID: categoryID,
		Questions:  allQuestions,
	}, nil
}

func (c *CarnivalService) getCurrentNodeCategoryID(sessionID uuid.UUID, nodeNumber int) (uuid.UUID, error) {
	s, err := c.sessionRepo.FindByID(sessionID)
	if err != nil {
		return uuid.Nil, err
	}

	if nodeNumber < 1 || nodeNumber > len([]string(s.Categories)) {
		return uuid.Nil, errors.New("invalid node number for session")
	}

	assignedCategoryName := ([]string(s.Categories))[nodeNumber-1] // Arrays are 0-indexed, nodes are 1-indexed

	categories, err := c.questionRepo.GetCategories()
	if err != nil {
		return uuid.Nil, err
	}

	for _, cat := range categories {
		if cat.Name == assignedCategoryName {
			return cat.ID, nil
		}
	}

	return uuid.Nil, errors.New("assigned category not found")
}

func (c *CarnivalService) countQuestionsAnsweredInCategory(sessionID, categoryID uuid.UUID) (int, error) {
	attempts, err := c.playerRepo.GetAttemptsBySessionAndCategory(sessionID, categoryID)
	if err != nil {
		return 0, err
	}
	return len(attempts), nil
}

func (c *CarnivalService) generateRandomCategoryAssignment() ([]string, error) {
	allCategories, err := c.questionRepo.GetCategories()
	if err != nil {
		return nil, err
	}

	var generalCategories []string
	for _, cat := range allCategories {
		if cat.Name != config.FUN_CATEGORY_NAME {
			generalCategories = append(generalCategories, cat.Name)
		}
	}

	if len(generalCategories) < config.MIN_REQUIRED_GENERAL_CATEGORIES {
		return nil, errors.New("insufficient general categories available")
	}

	if len(generalCategories) == config.MIN_REQUIRED_GENERAL_CATEGORIES {
		rand.Shuffle(len(generalCategories), func(i, j int) {
			generalCategories[i], generalCategories[j] = generalCategories[j], generalCategories[i]
		})
		return generalCategories, nil
	}

	// If more than required, randomly select the required number
	rand.Shuffle(len(generalCategories), func(i, j int) {
		generalCategories[i], generalCategories[j] = generalCategories[j], generalCategories[i]
	})

	return generalCategories[:config.MAX_CARNIVAL_NODES], nil
}

func (c *CarnivalService) getUniqueQuestionsFromCategory(categoryID uuid.UUID, limit int) ([]question.Question, error) {
	questions, err := c.questionRepo.GetQuestionsByCategory(categoryID, limit*config.QUESTION_FETCH_MULTIPLIER) // Get multiple to ensure uniqueness
	if err != nil {
		return nil, err
	}

	if len(questions) < limit {
		return nil, errors.New("insufficient questions in category")
	}

	// Shuffle and select the required number
	rand.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})

	return questions[:limit], nil
}

func (c *CarnivalService) getUnusedFunQuestionsForSession(sessionID uuid.UUID, limit int) ([]question.Question, error) {
	return c.questionRepo.GetUnusedFunQuestionsForSession(sessionID, limit)
}

func (c *CarnivalService) addNodeTimePenalty(s *session.Session, nodeNumber int) error {
	startTimeUnix, exists := s.NodeStartTimes[nodeNumber]
	if !exists {
		return nil
	}

	startTime := time.Unix(startTimeUnix, 0)
	elapsedSeconds := int(time.Since(startTime).Seconds())

	nodePenalty := elapsedSeconds / config.TIME_PENALTY_INTERVAL_SECONDS

	s.Score.TimePenalty += nodePenalty

	return nil
}

func (c *CarnivalService) updateLeaderboardAfterNode(s *session.Session, result *AnswerResult) error {
	s.Score = s.CalculateScore()
	result.CurrentScore = s.Score.Final

	p, err := c.playerRepo.FindByID(s.PlayerID)
	if err != nil {
		return err
	}

	currentTime := time.Since(s.StartedAt)
	entry := leaderboard.NewEntry(p.ID, p.Name, s.ID, s.Score.Final, currentTime)
	if err := c.leaderboardRepo.UpsertEntry(entry); err != nil {
		return err
	}

	return nil
}
