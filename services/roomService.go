package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/lafathalfath/go-chatserver/cache"
	"github.com/lafathalfath/go-chatserver/database"
	"github.com/lafathalfath/go-chatserver/graph/models"
	"github.com/lafathalfath/go-chatserver/helpers"

	"github.com/google/uuid"
)

type RoomService interface {
	MyRooms(ctx context.Context) ([]*models.Room, error)
	MyRoom(ctx context.Context, id string) (*models.Room, error)
	CreateDM(ctx context.Context, email string) (*models.Room, error)
	CreateGroup(ctx context.Context, input *models.NewRoom) (*models.Room, error)

	Typing(ctx context.Context, roomId string, isTyping bool) (bool, error)
	SubscribeUserTyping(ctx context.Context, roomId string) (<-chan *models.TypingEvent, error)
}

type roomService struct {
	*database.ConnectionObj
}

func RoomServices() RoomService {
	conn := &database.DBConnection
	return &roomService{conn}
}

func (s *roomService) Typing(ctx context.Context, roomId string, isTyping bool) (bool, error) {
	userId, ok := helpers.GetUserId(ctx)
	if !ok {
		return false, errors.New("Unauthorized")
	}
	userCacheKey := "users:"+userId
	var user models.User
	if err := cache.Scan(userCacheKey, &user); err != nil {
		if err := s.DB.Select("id", "name", "email").First(&user, "id = ?", userId).Error; err != nil {
			return false, err
		}
		cache.Set(userCacheKey, user, time.Hour)
	}
	typingEvent := models.TypingEvent{
		RoomID: roomId,
		User: &user,
		Typing: isTyping,
	}
	publishKey := "room:"+roomId+":usersTyping"
	payload, err := json.Marshal(typingEvent)
	if err != nil {
		return false, err
	}
	s.RDB.Publish(ctx, publishKey, payload)
	return true, nil
}

func (s *roomService) SubscribeUserTyping(ctx context.Context, roomId string) (<-chan *models.TypingEvent, error) {
	_, ok := helpers.GetUserId(ctx)
	if !ok {
		return nil, errors.New("Unauthorized")
	}

	ch := make(chan *models.TypingEvent)
	publishKey := "room:"+roomId+":usersTyping"
	sub := s.RDB.Subscribe(ctx, publishKey)

	go func()  {
		defer close(ch)
		defer sub.Close()

		for {
			select {
			case <-ctx.Done():
				return 
			case typ, ok := <-sub.Channel():
				if !ok {
					return 
				}
				var t models.TypingEvent
				if err := json.Unmarshal([]byte(typ.Payload), &t); err != nil {
					log.Println("invalid payload:", err)
					continue
				}
				ch <- &t
			}
		}
	}()
	return ch, nil
}

func (s *roomService) MyRooms(ctx context.Context) ([]*models.Room, error) {
	userID, ok := helpers.GetUserId(ctx)
	if !ok {
		return nil, errors.New("Unauthorized")
	}
	var rooms []*models.Room
	if err := s.DB.Joins("JOIN room_members rm ON rm.room_id = rooms.id").
		Select("rooms.id", "rooms.type", "rooms.name", "rooms.created_at").
		Where("rm.user_id = ?", userID).
		Preload("Members").
		Find(&rooms).Error; err != nil {
		return nil, err
	}
	for _, room := range rooms {
		if room.Type != "dm" {
			continue
		}
		for _, member := range room.Members {
			if member.ID != userID {
				room.Name = &member.Name
				break
			}
		}
	}
	return rooms, nil
}

func (s *roomService) MyRoom(ctx context.Context, id string) (*models.Room, error) {
	userId, ok := helpers.GetUserId(ctx)
	if !ok {
		return nil, errors.New("Unauthorized")
	}
	var room *models.Room
	if err := s.DB.Select("id", "type", "name", "created_at").Preload("Members").First(&room, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if room.Type == "dm" {
		for _, member := range room.Members {
			if member.ID == userId {
				continue
			}
			room.Name = &member.Name
		}
	}
	return room, nil
}

func (s *roomService) CreateDM(ctx context.Context, email string) (*models.Room, error) {
	myUserID, ok := helpers.GetUserId(ctx)
	if !ok {
		return nil, errors.New("Unauthorized")
	}
	var user models.User
	if err := s.DB.Select("id", "name").First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	if user.ID == myUserID {
		return nil, errors.New("Cannot create DM with this email")
	}
	members := []*models.User{
		{ID: myUserID},
		{ID: user.ID},
	}
	id := uuid.NewString()
	name := ""
	room := models.Room{
		ID:      id,
		Type:    "dm",
		Name:    &name,
		Members: members,
	}
	if err := s.DB.Create(&room).Error; err != nil {
		return nil, err
	}
	if err := s.DB.Joins("JOIN room_members rm ON rm.room_id = rooms.id").
		Select("rooms.id", "rooms.type", "rooms.name", "rooms.created_at").
		Where("rm.user_id = ?", myUserID).
		Preload("Members").
		Find(&room).Error; err != nil {
		return nil, err
	}
	room.Name = &user.Name
	return &room, nil
}

func (s *roomService) CreateGroup(ctx context.Context, input *models.NewRoom) (*models.Room, error) {
	myUserID, ok := helpers.GetUserId(ctx)
	if !ok {
		return nil, errors.New("Unauthorized")
	}
	var users []*models.User
	if err := s.DB.Select("id").
		Where("id != ?", myUserID).
		Where("email IN ?", input.MembersEmail).
		Find(&users).Error; err != nil {
		return nil, err
	}
	users = append(users, &models.User{ID: myUserID})
	id := uuid.NewString()
	room := models.Room{
		ID:      id,
		Type:    "group",
		Name:    &input.Name,
		Members: users,
	}
	if err := s.DB.Create(&room).Error; err != nil {
		return nil, err
	}
	if err := s.DB.Joins("JOIN room_members rm ON rm.room_id = rooms.id").
		Select("rooms.id", "rooms.type", "rooms.name", "rooms.created_at").
		Where("rm.user_id = ?", myUserID).
		Preload("Members").
		Find(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}
