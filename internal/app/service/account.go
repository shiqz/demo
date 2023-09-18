package service

import (
	"context"
	"demo/internal/domain"
	"demo/internal/domain/entity"
	"demo/internal/domain/types"
)

type AccountService struct {
	repo domain.AccountRepository
}

func NewAccountService(repo domain.AccountRepository) domain.AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) Login(ctx context.Context, email, pass string) (*entity.Session, error) {
	ag, err := s.repo.GetAccountByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if err = ag.Account.ValidPassword(pass); err != nil {
		return nil, err
	}
	return entity.NewSession(types.AdminSession, ag.Account.AdminID), nil
}

func (s *AccountService) Create(ctx context.Context, account *domain.AccountAggregate) error {
	return nil
}

func (s *AccountService) UpdatePassByEmail(ctx context.Context, email, pass string) error {
	return nil
}
