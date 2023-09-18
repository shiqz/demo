package commands

import (
	"context"
	"demo/internal/domain"
)

type UserCommandService struct {
}

func (c *UserCommandService) Create(ctx context.Context, ug *domain.UserAggregate) error {
	return nil
}
