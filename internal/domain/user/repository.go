package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	Get(ctx context.Context, uuid uuid.UUID) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, uuid uuid.UUID) error
}

type mysqlRepository struct {
	mc *sql.DB
}

func (m *mysqlRepository) Create(ctx context.Context, user *User) error {
	err := m.mc.Ping()
	if err != nil {
		return err
	}

	_, err = m.mc.Exec("INSERT INTO users(uuid, email, password_hash, orders_quantity) VALUES(?, ?, ?, ?)", user.UUID.String(),
		user.Email, user.PasswordHash, user.OrdersQuantity)
	return err
}

func (m *mysqlRepository) Get(ctx context.Context, uuid uuid.UUID) (*User, error) {
	err := m.mc.Ping()
	if err != nil {
		return nil, err
	}

	row := m.mc.QueryRow("SELECT * FROM users WHERE uuid=?", uuid.String())
	if row.Err() == sql.ErrNoRows {
		return nil, errors.New("no user with provided uuid")
	}
	u := &User{}
	err = row.Scan(&u.UUID, &u.Email, &u.PasswordHash, &u.OrdersQuantity)
	if err != nil {
		return nil, err
	}
	return u, err
}

func (m *mysqlRepository) Update(ctx context.Context, user *User) error {
	err := m.mc.Ping()
	if err != nil {
		return err
	}

	_, err = m.mc.Exec("UPDATE USERS SET email=?, password_hash=?, orders_quantity=?", user.Email, user.PasswordHash, user.OrdersQuantity)
	return err
}

func (m *mysqlRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	err := m.mc.Ping()
	if err != nil {
		return err
	}

	_, err = m.mc.Exec("DELETE FROM users WHERE uuid=?", uuid.String())
	return err
}

func NewRepository(mc *sql.DB) Repository {
	return &mysqlRepository{mc: mc}
}
