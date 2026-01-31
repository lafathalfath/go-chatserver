package services

import (
	"fmt"
	"github.com/lafathalfath/go-chatserver/cache"
	"github.com/lafathalfath/go-chatserver/database"
	"github.com/lafathalfath/go-chatserver/graph/models"
	"github.com/lafathalfath/go-chatserver/helpers"
	"github.com/lafathalfath/go-chatserver/validators"
	"time"

	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(input *models.NewUser) (*models.User, error)
	GetUsers() ([]*models.User, error)
	GetUser(id string) (*models.User, error)
}

type userService struct {
	*database.ConnectionObj
}

func UserServices() UserService {
	conn := &database.DBConnection
	return &userService{conn}
}

func (s *userService) CreateUser(input *models.NewUser) (*models.User, error) {
	if err := validators.ValidateNewUser(s.DB, *input); err != nil {
		return nil, err
	}
	hash, err := helpers.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	id := uuid.NewString()
	user := models.User{
		ID:       id,
		Name:     input.Name,
		Email:    input.Email,
		Password: hash,
	}
	if err := s.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	cacheKey := fmt.Sprintf("users:%s", id)
	cache.Set(cacheKey, user, 10*time.Minute)
	cacheKey = "users"
	var users []models.User
	if err := cache.Scan(cacheKey, &users); err == nil {
		s.DB.Select("id", "name", "email", "password").Find(&users)
		cache.Set(cacheKey, users, 10*time.Minute)
	}
	return &user, nil
}

func (s *userService) GetUsers() ([]*models.User, error) {
	var users []*models.User
	cacheKey := "users"
	if err := cache.Scan(cacheKey, &users); err == nil {
		return users, nil
	} else if err := s.DB.Select("id", "name", "email", "password").Find(&users).Error; err != nil {
		return nil, err
	}
	cache.Set(cacheKey, users, 10*time.Minute)
	return users, nil
}

func (s *userService) GetUser(id string) (*models.User, error) {
	var user *models.User
	cacheKey := fmt.Sprintf("users:%s", id)
	if err := cache.Scan(cacheKey, &user); err == nil {
		return user, nil
	} else if err := s.DB.Select("id", "name", "email", "password").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	cache.Set(cacheKey, user, 10*time.Minute)
	return user, nil
}
