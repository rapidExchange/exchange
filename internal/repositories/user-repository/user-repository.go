package userrepository

import (
	"context"
	"database/sql"
	"errors"
	"rapidEx/internal/domain/user"
)

type Repository interface {
	Create(ctx context.Context, user *user.User) error
	Get(ctx context.Context, email string) (*user.User, error)
	Update(ctx context.Context, user *user.User) error
	Delete(ctx context.Context, email string) error
}

type mysqlRepository struct {
	mc *sql.DB
}

func (m *mysqlRepository) Create(ctx context.Context, user *user.User) error {
	err := m.mc.Ping()
	if err != nil {
		return err
	}

	_, err = m.mc.Exec("INSERT INTO users(uuid, email, password_hash) VALUES(?, ?, ?)", user.UUID.String(),
		user.Email, user.PasswordHash)
	return err
}

// TODO: refactor
func (m *mysqlRepository) Get(ctx context.Context, email string) (*user.User, error) {
	err := m.mc.Ping()
	if err != nil {
		return nil, err
	}

	row := m.mc.QueryRowContext(context.Background(), "SELECT * FROM users WHERE email=?", email)
	if err != nil {
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
	rows, err := m.mc.QueryContext(context.Background(), "SELECT ticker, quantity FROM balance WHERE email = ?", email)
	if err != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			return nil, nil
		}
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
	u.Balance = balance
	return u, nil
}

func (m *mysqlRepository) Update(ctx context.Context, user *user.User) error {
	err := m.mc.Ping()
	if err != nil {
		return err
	}

	_, err = m.mc.ExecContext(context.Background(),
		"UPDATE users SET password_hash=? WHERE email=?",
		user.PasswordHash, user.Email)
	for ticker, quantity := range user.Balance {
		_, err = m.mc.ExecContext(context.Background(),
			`IF EXISTS (SELECT ticker FROM balance WHERE(email=? AND ticker=?))
	BEGIN
		UPDATE balance SET quantity=? WHERE(email=? AND ticker=?)
	END
	ELSE
	BEGIN
		INSERT INTO balance(email, ticker, quantity) VALUES(?, ?, ?)`,
			user.Email, ticker, quantity, user.Email, ticker, user.Email, ticker, quantity)
	}
	return err
}

func (m *mysqlRepository) Delete(ctx context.Context, email string) error {
	err := m.mc.Ping()
	if err != nil {
		return err
	}

	_, err = m.mc.Exec("DELETE FROM users WHERE email=?", email)
	return err
}

func NewUserRepository(mc *sql.DB) Repository {
	return &mysqlRepository{mc: mc}
}
