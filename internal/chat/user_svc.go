package chat

import (
	"errors"
	"fmt"

	"github.com/boxuanduan/gochat/config"
	"github.com/boxuanduan/gochat/internal/repo/mysql"
	"github.com/boxuanduan/gochat/pkg/auth"
	"github.com/boxuanduan/gochat/pkg/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo   *mysql.UserRepo
	jwtCfg config.JWTConfig
}

func NewUserService(repo *mysql.UserRepo, jwtCfg config.JWTConfig) *UserService {
	return &UserService{
		repo:   repo,
		jwtCfg: jwtCfg,
	}
}

func (s *UserService) Register(username, password, nickname string) (string, error) {
	// Check if username already exists
	_, err := s.repo.FindByUsername(username)

	if err == nil {
		return "", fmt.Errorf("username already exists")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	// Create new user

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password),
		bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}

	user := &model.User{
		Username: username,
		Password: string(hashedPassword),
		Nickname: nickname,
	}

	err = s.repo.Create(user)

	if err != nil {
		return "", fmt.Errorf("failed to create user: %v", err)
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, s.jwtCfg.Secret, s.jwtCfg.Expire)

	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return token, nil
}

func (s *UserService) Login(username, password string) (string, error) {
	// Find user by username
	user, err := s.repo.FindByUsername(username)

	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, s.jwtCfg.Secret, s.jwtCfg.Expire)

	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return token, nil
}
