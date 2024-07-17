package service

import (
	"context"
	"log"
	"nedorezov/pkg/db"
	"nedorezov/pkg/domain/account"
	"strings"
)

func (s *service) Registration(ctx context.Context, userId string, newAccount account.Registration) (*account.Info, error) {
	if userId != "" {
		return nil, db.ErrAuthorize
	}
	//check valid data
	if !isValidAccountRegister(newAccount) {
		log.Println("invalid data")
		return nil, db.ErrValidate
	}

	result, err := s.rAccount.Registration(ctx, newAccount)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) Login(ctx context.Context, acc account.Login) (int64, error) {
	id, err := s.rAccount.Login(ctx, acc)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *service) PutAccount(ctx context.Context, id string, updateAcc account.Registration) (*account.Info, error) {
	idInt, err := s.checkIdParam(id)
	if err != nil {
		return nil, err
	}

	// Валидация обновляемых данных
	if !isValidAccountRegister(updateAcc) {
		return nil, db.ErrValidate
	}

	result, err := s.rAccount.Put(ctx, idInt, updateAcc)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) GetBalance(ctx context.Context, id string) (*float64, error) {
	idInt, err := s.checkIdParam(id)
	if err != nil {
		return nil, err
	}
	currAccount, err := s.rAccount.Balance(ctx, idInt)
	if err != nil {
		return nil, err
	}
	return &currAccount.Balance, nil
}

func (s *service) Deposit(ctx context.Context, id string, amount float64) (*float64, error) {
	if !isValidAmount(amount) {
		return nil, db.ErrNotValidAmount
	}

	idInt, err := s.checkIdParam(id)
	if err != nil {
		return nil, err
	}

	currAccount, err := s.rAccount.Balance(ctx, idInt)
	if err != nil {
		return nil, err
	}
	currAccount.Balance += amount

	balance, err := s.rAccount.ChangeBalance(ctx, idInt, currAccount.Balance)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (s *service) WithDraw(ctx context.Context, id string, amount float64) (*float64, error) {
	if !isValidAmount(amount) {
		return nil, db.ErrNotValidAmount
	}

	idInt, err := s.checkIdParam(id)
	if err != nil {
		return nil, err
	}

	currAccount, err := s.rAccount.Balance(ctx, idInt)
	if err != nil {
		return nil, err
	}
	currAccount.Balance -= amount

	balance, err := s.rAccount.ChangeBalance(ctx, idInt, currAccount.Balance)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func isValidAmount(amount float64) bool {
	return amount >= 0.0
}

func isValidAccountRegister(newAccountRegistration account.Registration) bool {
	if newAccountRegistration.FirstName == "" || strings.TrimSpace(newAccountRegistration.FirstName) == "" {
		return false
	}

	if newAccountRegistration.LastName == "" || strings.TrimSpace(newAccountRegistration.LastName) == "" {
		return false
	}

	if newAccountRegistration.CardNumber == "" || strings.TrimSpace(newAccountRegistration.CardNumber) == "" {
		return false
	}

	if newAccountRegistration.Password == "" || strings.TrimSpace(newAccountRegistration.Password) == "" {
		return false
	}

	return true
}
