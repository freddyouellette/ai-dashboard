package bots

import (
	"fmt"

	"github.com/freddyouellette/ai-dashboard/internal/models"
)

type BotRepository interface {
	GetAll() ([]models.Bot, error)
	GetByID(id uint) (models.Bot, error)
	Create(bot models.Bot) (models.Bot, error)
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

func (s *BotService) GetAll() ([]models.Bot, error) {
	bots, err := s.botRepository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRepository, err.Error())
	}
	return bots, nil
}

func (s *BotService) Create(bot models.Bot) (models.Bot, error) {
	bot, err := s.botRepository.Create(bot)
	if err != nil {
		return bot, fmt.Errorf("%w: %s", ErrRepository, err.Error())
	}
	return bot, nil
}

func (s *BotService) GetById(id uint) (*models.Bot, error) {
	bot, err := s.botRepository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRepository, err.Error())
	}
	return &bot, nil
}
