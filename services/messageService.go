package services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/lafathalfath/go-chatserver/database"
	"github.com/lafathalfath/go-chatserver/graph/models"
	"github.com/lafathalfath/go-chatserver/helpers"
	"log"

	"github.com/google/uuid"
)

type MessageService interface {
	SendMessage(ctx context.Context, input *models.NewMessage) (*models.Message, error)
	SubscribeMessage(ctx context.Context, roomId string) (<-chan *models.Message, error)
	GetMessages(roomId string) ([]*models.Message, error)
}

type messageService struct {
	*database.ConnectionObj
}

func MessageServices() MessageService {
	conn := &database.DBConnection
	return &messageService{conn}
}

func (s *messageService) SendMessage(ctx context.Context, input *models.NewMessage) (*models.Message, error) {
	userID, ok := helpers.GetUserId(ctx)
	if !ok {
		return nil, errors.New("Unauthorized")
	}
	id := uuid.NewString()
	msg := &models.Message{
		ID:       id,
		SenderID: userID,
		RoomID:   input.RoomID,
		Content:  input.Content,
	}
	if err := s.DB.Create(&msg).Error; err != nil {
		return nil, err
	}
	if err := s.DB.
		Select("id", "content", "sender_id", "room_id", "created_at").
		Preload("Sender").
		First(&msg, "id = ?", id).Error; err != nil {
		return nil, err
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	publishKey := "room:" + input.RoomID
	s.RDB.Publish(ctx, publishKey, payload)
	return msg, nil
}

func (s *messageService) SubscribeMessage(ctx context.Context, roomId string) (<-chan *models.Message, error) {
	_, ok := helpers.GetUserId(ctx)
	if !ok {
		return nil, errors.New("Unauthorized")
	}
	ch := make(chan *models.Message)
	publishKey := "room:" + roomId
	sub := s.RDB.Subscribe(ctx, publishKey)

	go func() {
		defer close(ch)
		defer sub.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-sub.Channel():
				if !ok {
					return
				}
				var m models.Message
				if err := json.Unmarshal([]byte(msg.Payload), &m); err != nil {
					log.Println("invalid payload:", err)
					continue
				}
				
				ch <- &m
			}
		}
	}()

	return ch, nil
}

func (s *messageService) GetMessages(roomId string) ([]*models.Message, error) {
	var messages []*models.Message
	if err := s.DB.
		Select("id", "content", "sender_id", "room_id", "created_at").
		Preload("Sender").
		Find(&messages, "room_id = ?", roomId).Error; err != nil {
		return nil, err
	}
	log.Println(messages)
	return messages, nil
}
