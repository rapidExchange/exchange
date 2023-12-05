package test

import (
	"context"
	"errors"
	"testing"

	"rapidEx/internal/domain/user"
	mysqlconnect "rapidEx/internal/mysql-connect"
)

func TestUser(t *testing.T) {
	mc, err := mysqlconnect.SetMysqlConnection()
	if err != nil {
		t.Error(err)
		return
	}
	userRepository := user.NewRepository(mc)
	u1 := user.New("a@gmail.com", "awfag4319285ygq2h0")
	err = userRepository.Create(context.Background(), u1)
	if err != nil {
		t.Error(err)
		return
	}
	u2, err := userRepository.Get(context.Background(), u1.Email)
	if err != nil {
		t.Error(err)
		return
	}
	if u1.Email != u2.Email || u1.PasswordHash != u2.PasswordHash || u1.UUID != u2.UUID {
		t.Error(errors.New("users are not equal"))
		return
	}
	err = userRepository.Delete(context.Background(), u1.Email)
	if err != nil {
		t.Error(err)
		return
	}
}