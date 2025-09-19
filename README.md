# Haoma — Black-Box Carnival 🎪

> Persian god meets cyber trials. A 7-node security quiz for ELECOMP 1404.

In the shadow-realm where code meets chaos, **Haoma** opens his eternal carnival. Seven tents shimmer in digital twilight—each a trial, each a doorway. Enter as student, emerge as guardian of the cyber-realm.

## Architecture 🏗️

Built with **Go + Gin** following **Hexagonal/DDD** principles:

- **Domain Layer**: Pure business logic (Session, Question, Player, Leaderboard)
- **Application Layer**: Use cases and services
- **Infrastructure Layer**: Database, HTTP, Excel seeding
- **Adapters Layer**: HTTP handlers, repository implementations

## Quick Start 🚀

### Prerequisites
- Go 1.23+
- Make (optional but recommended)

### Setup & Run
```bash
# Start the carnival
make dev-env-build

# fill the database with excel files(when docker is running)
make seed-excel
```

### Alternative Setup
```bash
# Manual setup
go mod tidy
go run cmd/server/main.go
```

## API Endpoints 🎯

The carnival speaks through these mystical endpoints:

### **Authentication**
- `POST /api/v1/auth/signup` — Create player account
- `POST /api/v1/auth/login` — Authenticate & get JWT token
- `GET /api/v1/auth/profile` — Get player profile

### **Game Flow**  
- `POST /api/v1/sessions/start` — Begin the journey
- `POST /api/v1/nodes/scan` — Scan QR codes at physical locations
- `POST /api/v1/sessions/{id}/answer` — Answer riddles
- `GET /api/v1/leaderboard` — View champions

**Key Features:**
- 🔐 **JWT Authentication** - Secure player verification
- 🚫 **Duplicate Prevention** - Each question answerable only once  
- 📊 **Real-time Leaderboard** - Updates after each node completion
- 🎯 **Location-based** - Physical QR codes at carnival stations

## Explore 🗺️

- **Swagger UI**: http://localhost:8080/docs
- **Health Check**: http://localhost:8080/health  
- **API Base**: http://localhost:8080/api/v1

## Development 🛠️

### Available Commands
```bash
make help         # Show all available targets
make deps         # Install dependencies
make build        # Build the carnival
make run          # Start server
make dev          # Hot reload development
make test         # Run tests
make test-coverage # Test with coverage report
make lint         # Code quality checks
make fmt          # Format code
make seed         # Create sample data
make seed-excel   # Load from Excel files
make swagger      # Generate API docs
make clean        # Clean artifacts
```

### Project Structure
```
haoma/
├── cmd/server/           # Application entry point
├── internal/
│   ├── domain/          # Business entities & rules
│   ├── application/     # Use cases & services  
│   ├── infrastructure/  # External concerns
│   └── adapters/        # Interface adapters
├── api/                 # API specifications
├── data/                # Excel seed files
└── Makefile            # Development commands
```

### Testing
```bash
make test              # Run all tests
make test-coverage     # Generate coverage report
```

### Data Seeding

**From Excel Files** (preferred):
1. Place `SCENARIOS.xlsx` and `questions.xlsx` in `data/` folder
2. Run: `make seed-excel`


Excel format:
- **SCENARIOS.xlsx**: Category definitions with PhDT marking
- **questions.xlsx**: Questions with options A-D and correct answers

## Game Rules ⚖️

The carnival follows ancient laws:

- **7 Unique Categories**: Each session selects 7 from 8 available
- **5 Questions per Node**: 4 from category + 1 from "Fun" pool  
- **PhDT Special**: Phishing questions use only A/B options
- **Per-Node Timing**: Time penalty calculated separately for each node
- **Scoring**: `(correct × 100) - accumulated_time_penalties`
- **Real-time Competition**: Leaderboard updates after each node completion
- **One Chance Rule**: Each question can only be answered once per session
- **Time Limit**: 2 hours maximum per session
- **Physical Movement**: Must scan QR codes at actual carnival locations

## Etymology 📜

**Haoma** (هوما) derives from Persian Zoroastrian mythology—the divine bird of fortune and the sacred plant of immortality. In our digital realm, it represents the transformative journey from student to cyber-guardian.

## Credits 🙏

*Inspired by Zoroastrian lore and the eternal dance of challenge and wisdom.*

Built with:
- [Go](https://golang.org/) - The language of gophers
- [Gin](https://gin-gonic.com/) - HTTP web framework
- [GORM](https://gorm.io/) - ORM library
- [Excelize](https://github.com/xuri/excelize) - Excel document processing
- [Swagger](https://swagger.io/) - API documentation

---

*The carnival awaits. Will you answer its call?* ✨
