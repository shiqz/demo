package service

import (
	"context"
	"database/sql"
	"demo/internal/app/errs"
	"demo/internal/domain"
	"demo/internal/domain/entity"
	"demo/internal/domain/types"
	"github.com/pkg/errors"
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.EcInvalidUser
		}
		return nil, err
	}
	if err = ag.Account.ValidPassword(pass); err != nil {
		return nil, err
	}
	return entity.NewSession(types.AdminSession, ag.Account.AdminID), nil
}

func (s *AccountService) Create(ctx context.Context, account *domain.AccountAggregate) error {
	ag, err := s.repo.GetAccountByEmail(ctx, account.Account.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if ag != nil {
		return errs.EcUserHasBeenExist
	}
	return s.repo.Save(ctx, account)
}

func (s *AccountService) UpdatePassByEmail(ctx context.Context, email, pass string) error {
	return nil
}
