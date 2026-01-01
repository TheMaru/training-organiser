package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	val := ctx.Value(UserIDKey)
	if val == nil {
		return uuid.Nil, errors.New("No user id in context")
	}

	userID, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("User id is not a uuid")
	}

	return userID, nil
}
