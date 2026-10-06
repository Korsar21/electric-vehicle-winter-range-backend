package auth

import (
	"sync"

	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/model"
)

const fixedUserID uint = 1

var (
	currentUser     *model.User
	currentUserOnce sync.Once
)

func CurrentUser() *model.User {
	currentUserOnce.Do(func() {
		currentUser = &model.User{
			ID: fixedUserID,
		}
	})

	return currentUser
}

func CurrentUserID() uint {
	return CurrentUser().ID
}
