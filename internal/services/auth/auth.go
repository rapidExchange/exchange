package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"rapidEx/internal/domain/user"
	jwt "rapidEx/internal/lib"
	"rapidEx/internal/storage"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredenitals = errors.New("invalid credentials")
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
}

type UserProvider interface {
	User(ctx context.Context, email string) (*user.User, error)
}

type UserSaver interface {
	Set(ctx context.Context, user *user.User) error
}

func New(log *slog.Logger, userSaver UserSaver, userProvider UserProvider) *Auth {
	return &Auth{log: log,
		userSaver:    userSaver,
		userProvider: userProvider}
}

func (a *Auth) Register(ctx context.Context, email, password string) error {
	const op = "Auth.Register"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Info("failed to generate password hash")
	}
	user := user.New(email, string(passwordHash))
	err = a.userSaver.Set(ctx, user)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Info("user already exists")
			return fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
		}
		log.Error("failed to register user")
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *Auth) Login(ctx context.Context, email, password string) (string, error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	log.Info("attempting to login user")
	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Info("user not found")
			return "", fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		log.Info("failed to login user")
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		log.Info("invalid credentials")
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredenitals)
	}
	log.Info("user logged in successfully")
	token, err := jwt.NewToken(*user, time.Hour*12)
	if err != nil {
		log.Info("failed to generate token")
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}
