package commands

import (
	"context"
	"demo/internal/domain/entity"
)

type UserCommandService struct {
}

func (c *UserCommandService) Create(ctx context.Context, user *entity.User) error {
	return nil
}
