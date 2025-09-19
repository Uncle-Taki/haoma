package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func JWTMiddleware(jwtService *JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header - Bearer token required"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		claims, err := jwtService.ValidatePlayerToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("player_id", claims.PlayerID)
		c.Set("player_name", claims.PlayerName)
		c.Set("player_email", claims.PlayerEmail)
		c.Set("player_claims", claims)

		c.Next()
	}
}

func SessionMiddleware(jwtService *JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken := c.GetHeader("X-Session-Token")
		if sessionToken == "" {
			sessionToken = c.Query("session_token")
		}

		if sessionToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing session token"})
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateSessionToken(sessionToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session token"})
			c.Abort()
			return
		}

		// Verify session belongs to authenticated player
		playerID, exists := c.Get("player_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Player not authenticated"})
			c.Abort()
			return
		}

		if claims.PlayerID != playerID.(uuid.UUID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Session does not belong to authenticated player"})
			c.Abort()
			return
		}

		c.Set("session_id", claims.SessionID)
		c.Set("session_claims", claims)

		c.Next()
	}
}

func OptionalSessionMiddleware(jwtService *JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken := c.GetHeader("X-Session-Token")
		if sessionToken == "" {
			sessionToken = c.Query("session_token")
		}

		if sessionToken == "" {
			c.Next()
			return
		}

		claims, err := jwtService.ValidateSessionToken(sessionToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session token"})
			c.Abort()
			return
		}

		// Verify session belongs to authenticated player
		playerID, exists := c.Get("player_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Player not authenticated"})
			c.Abort()
			return
		}

		if claims.PlayerID != playerID.(uuid.UUID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Session does not belong to authenticated player"})
			c.Abort()
			return
		}

		c.Set("session_id", claims.SessionID)
		c.Set("session_claims", claims)

		c.Next()
	}
}
