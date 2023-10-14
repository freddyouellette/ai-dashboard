package bots

import (
	"fmt"

	"github.com/freddyouellette/ai-dashboard/internal/models"
)

type BotService struct{}

func NewBotService() *BotService {
	return &BotService{}
}

func (s *BotService) GetBots() ([]models.Bot, error) {
	return []models.Bot{{
		Uuid:        "9f976d75-8aec-4d92-a9e5-c7e810f7051f",
		Name:        "Bot 1",
		Description: "Bot 1 description",
		AiModel:     "Bot 1 model",
		Personality: "Bot 1 personality",
		UserHistory: "Bot 1 user history",
	}}, nil
}

func (s *BotService) GetBotById(id string) (*models.Bot, error) {
	return nil, fmt.Errorf("%w: %s", models.ErrResourceNotFound, fmt.Sprintf("bot with id %s not found", id))
}
