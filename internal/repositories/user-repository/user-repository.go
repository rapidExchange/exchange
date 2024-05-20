package userrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"rapidEx/internal/domain/user"
	"rapidEx/internal/storage"

	"github.com/go-sql-driver/mysql"
)

const (
	SetUserQuery    = "INSERT INTO users(uuid, email, password_hash) VALUES(?, ?, ?)"
	GetUserQuery    = "SELECT * FROM users WHERE email=?"
	UpdateUserQuery = "UPDATE users SET password_hash=? WHERE email=?"
	DeleteUserQuery = "DELETE FROM users WHERE email=?"

	UpdateUserBalanceSheetQuery = `IF EXISTS (SELECT ticker FROM balance WHERE(email=? AND ticker=?))
	BEGIN
		UPDATE balance SET quantity=? WHERE(email=? AND ticker=?)
	END
	ELSE
	BEGIN
		INSERT INTO balance(email, ticker, quantity) VALUES(?, ?, ?)`

	GetTickerQuery = "SELECT ticker, quantity FROM balance WHERE email = ?"
)

type Repository interface {
	Set(ctx context.Context, user *user.User) error
	User(ctx context.Context, email string) (*user.User, error)
	Update(ctx context.Context, user *user.User) error
	Delete(ctx context.Context, email string) error
}

type mysqlRepository struct {
	mc *sql.DB
}

func (m *mysqlRepository) Set(ctx context.Context, user *user.User) error {
	const op = "userRepository.Set"
	err := m.mc.Ping()
	if err != nil {
		return err
	}

	stmt, err := m.mc.Prepare(SetUserQuery)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.ExecContext(ctx, user.UUID.String(), user.Email, user.PasswordHash)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return storage.ErrUserExists
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (m *mysqlRepository) User(ctx context.Context, email string) (*user.User, error) {
	const op = "userRepository.User"
	err := m.mc.Ping()
	if err != nil {
		return nil, err
	}
	row := m.mc.QueryRowContext(context.Background(), GetUserQuery, email)
	if row.Err() != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			return nil, errors.New("userNotFound")
		}
		return nil, err
	}
	u := &user.User{}
	err = row.Scan(&u.UUID, &u.Email, &u.PasswordHash)
	if err != nil {
		return nil, err
	}
	balanceSheet, err := m.userBalanceSheet(u.Email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	u.Balance = balanceSheet

	return u, nil
}

func (m *mysqlRepository) userBalanceSheet(email string) (map[string]float64, error) {
	rows, err := m.mc.QueryContext(context.Background(), GetTickerQuery, email)
	if err != nil {
		return nil, err
	}

	var ticker string
	var quantity float64

	balance := make(map[string]float64)
	for rows.Next() {
		err := rows.Scan(&ticker, &quantity)
		if err != nil {
			return nil, err
		}
		balance[ticker] = quantity
	}
	return balance, nil
}

func (m *mysqlRepository) Update(ctx context.Context, user *user.User) error {
	err := m.mc.Ping()
	if err != nil {
		return err
	}

	_, err = m.mc.ExecContext(context.Background(),
		UpdateUserQuery,
		user.PasswordHash, user.Email)
	for ticker, quantity := range user.Balance {
		_, err = m.mc.ExecContext(context.Background(), UpdateUserBalanceSheetQuery,
			user.Email, ticker, quantity, user.Email, ticker, user.Email, ticker, quantity)
	}
	return err
}

func (m *mysqlRepository) Delete(ctx context.Context, email string) error {
	err := m.mc.Ping()
	if err != nil {
		return err
	}

	_, err = m.mc.Exec(DeleteUserQuery, email)
	return err
}

func NewUserRepository(mc *sql.DB) Repository {
	return &mysqlRepository{mc: mc}
}
