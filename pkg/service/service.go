package service

import (
	"nedorezov/pkg/db"
	accountI "nedorezov/pkg/repo/account/interface"

	"context"
	interfaces "nedorezov/pkg/service/interface"
	"strconv"
)

type service struct {
	rAccount accountI.AccountRepository
}

func NewService(
	accountRepository accountI.AccountRepository,
) interfaces.ServiceUseCase {
	return &service{
		rAccount: accountRepository,
	}
}

func (s *service) Migrate(ctx context.Context) error {
	if err := s.rAccount.Migrate(ctx); err != nil {
		return err
	}

	return nil
}

func (s *service) checkLogin(ctx context.Context, userId string) (int64, error) {
	userIdInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return 0, db.ErrAuthorize
	}
	_, err = s.rAccount.Balance(ctx, userIdInt)
	if err != nil {
		return 0, db.ErrAuthorize
	}
	return userIdInt, nil
}

func (s *service) checkIdParam(id string) (int64, error) {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil || idInt <= 0 {
		return 0, db.ErrParamNotFound
	}
	return idInt, nil
}
