package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetUserID mengambil user_id dari gin.Context.
// Panggil fungsi ini di handler setelah melewati AuthMiddleware.
func GetUserID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, errors.New("user_id not found in context")
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user_id is not a valid UUID")
	}

	return uid, nil
}
