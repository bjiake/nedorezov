package account

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"nedorezov/pkg/db"
	"nedorezov/pkg/domain/account"
	interfaces "nedorezov/pkg/repo/account/interface"

	"github.com/jackc/pgconn"
)

type accountDataBase struct {
	db *sql.DB
}

func NewAccountDataBase(db *sql.DB) interfaces.AccountRepository {
	return &accountDataBase{
		db: db,
	}
}

func (r *accountDataBase) Migrate(ctx context.Context) error {
	accQuery := `
    CREATE TABLE IF NOT EXISTS account(
       	id SERIAL PRIMARY KEY,
		firstName text not NULL,
		lastName text not NULL,
		balance double precision,
		cardNumber text not NULL,
		password text not NULL
    );
    `
	_, err := r.db.ExecContext(ctx, accQuery)
	if err != nil {
		message := db.ErrMigrate.Error() + " account"
		log.Printf("%q: %s\n", message, err.Error())
		return db.ErrMigrate
	}

	return err
}

func (r *accountDataBase) Registration(ctx context.Context, newAccount account.Registration) (*account.Info, error) {
	// Check if an account with the same email already exists
	var existingCount int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM account WHERE cardNumber = $1", newAccount.CardNumber).Scan(&existingCount)
	if err != nil {
		return nil, err
	}

	if existingCount > 0 {
		return nil, db.ErrDuplicate
	}

	var id int64

	err = r.db.QueryRowContext(ctx,
		"INSERT INTO account(firstName, lastName, balance, cardNumber, password) values($1, $2, $3, $4, $5) RETURNING id",
		newAccount.FirstName, newAccount.LastName, 0.0, newAccount.CardNumber, newAccount.Password).Scan(&id)
	// Check if a user with the same email already exists
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, db.ErrDuplicate
			}
		}
		return nil, err
	}

	// Add the new account
	requestAccount := &account.Info{
		ID:         id,
		FirstName:  newAccount.FirstName,
		LastName:   newAccount.LastName,
		Balance:    0.0,
		CardNumber: newAccount.CardNumber,
	}

	return requestAccount, nil
}

func (r *accountDataBase) Login(ctx context.Context, acc account.Login) (int64, error) {
	var id int64
	row := r.db.QueryRowContext(ctx, "SELECT id FROM account WHERE cardNumber = $1 and password = $2", acc.CardNumber, acc.Password)

	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, db.ErrNotExist
		}
		return 0, err
	}
	return id, nil
}

func (r *accountDataBase) Put(ctx context.Context, id int64, updateAcc account.Registration) (*account.Info, error) {
	res, err := r.db.ExecContext(ctx, "UPDATE account SET firstName = $1, lastName = $2, cardNumber = $3, Password = $4 WHERE account.id = $5",
		updateAcc.FirstName, updateAcc.LastName, updateAcc.CardNumber, updateAcc.Password, id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, db.ErrDuplicate
			}
		}
		return nil, err
	}

	result := &account.Info{
		ID:         id,
		FirstName:  updateAcc.FirstName,
		LastName:   updateAcc.LastName,
		Balance:    0,
		CardNumber: updateAcc.CardNumber,
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, db.ErrUpdateFailed
	}

	return result, nil
}

func (r *accountDataBase) ChangeBalance(ctx context.Context, id int64, balance float64) (*float64, error) {
	res, err := r.db.ExecContext(ctx, "UPDATE account SET balance = $1 WHERE account.id = $2", balance, id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, db.ErrDuplicate
			}
		}
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, db.ErrUpdateFailed
	}

	return &balance, nil
}

func (r *accountDataBase) Balance(ctx context.Context, id int64) (*account.Info, error) {
	row := r.db.QueryRowContext(ctx, "SELECT * FROM account WHERE id = $1", id)

	var currAccount account.Account
	if err := row.Scan(&currAccount.ID, &currAccount.FirstName, &currAccount.LastName, &currAccount.Balance, &currAccount.CardNumber, &currAccount.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, db.ErrNotExist
		}
		return nil, err
	}
	var result = account.Info{
		ID:         id,
		FirstName:  currAccount.FirstName,
		LastName:   currAccount.LastName,
		Balance:    currAccount.Balance,
		CardNumber: currAccount.CardNumber,
	}

	return &result, nil
}
