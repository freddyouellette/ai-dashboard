package messages_repository

import (
	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/freddyouellette/ai-dashboard/internal/repositories/entity_repository"
	"gorm.io/gorm"
)

type MessagesRepository struct {
	*entity_repository.EntityRepository[models.Message]
	db *gorm.DB
}

func NewMessagesRepository(
	entityRepository *entity_repository.EntityRepository[models.Message],
	db *gorm.DB,
) *MessagesRepository {
	return &MessagesRepository{
		EntityRepository: entityRepository,
		db:               db,
	}
}

func (r *MessagesRepository) GetAllPaginated(options *models.GetMessagesOptions) (*models.MessagesDTO, error) {
	dto := &models.MessagesDTO{
		Page:     options.Page,
		PerPage:  options.PerPage,
		Messages: []*models.Message{},
	}

	offset := (options.Page - 1) * options.PerPage
	query := r.db.Model(&models.Message{})
	if options.ChatID != 0 {
		query = query.Where("chat_id = ?", options.ChatID)
	}
	query.
		Order("created_at DESC").
		Offset(offset).
		Limit(options.PerPage)

	if err := query.Find(&dto.Messages).Error; err != nil {
		return nil, err
	}

	var totalMessages int64
	if err := r.db.
		Model(&models.Message{}).
		Where("chat_id = ?", options.ChatID).
		Count(&totalMessages).Error; err != nil {
		return nil, err
	}
	dto.Total = int(totalMessages)

	if len(dto.Messages) == 0 {
		// if there are no entities, return empty slice instead of nil...
		dto.Messages = []*models.Message{}
	}

	return dto, nil
}
