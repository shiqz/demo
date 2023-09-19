package commands

import "demo/internal/domain"

type UserHandler struct {
	userService    domain.UserService
	sessionService domain.SessionService
}
