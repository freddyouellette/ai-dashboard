package bots

import (
	"fmt"

	"github.com/freddyouellette/ai-dashboard/internal/models"
)

type BotRepository interface {
	GetAll() ([]models.Bot, error)
	GetByID(id uint) (models.Bot, error)
}

type BotService struct {
	botRepository BotRepository
}

var (
	ErrRepository = fmt.Errorf("repository error")
)

func NewBotService(botRepository BotRepository) *BotService {
	return &BotService{
		botRepository: botRepository,
	}
}

func (s *BotService) GetBots() ([]models.Bot, error) {
	bots, err := s.botRepository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRepository, err.Error())
	}
	return bots, nil
}

func (s *BotService) GetBotById(id uint) (*models.Bot, error) {
	bot, err := s.botRepository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRepository, err.Error())
	}
	return &bot, nil
}
