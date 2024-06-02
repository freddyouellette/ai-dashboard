package messages_repository

import (
	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/freddyouellette/ai-dashboard/internal/repositories/user_scoped_repository"
	"gorm.io/gorm"
)

type MessagesRepository struct {
	*user_scoped_repository.UserScopedRepository[*models.Message]
	db *gorm.DB
}

func NewMessagesRepository(
	entityRepository *user_scoped_repository.UserScopedRepository[*models.Message],
	db *gorm.DB,
) *MessagesRepository {
	return &MessagesRepository{
		UserScopedRepository: entityRepository,
		db:                   db,
	}
}

func (r *MessagesRepository) GetAllPaginated(userId uint, options *models.GetMessagesOptions) (*models.MessagesDTO, error) {
	dto := &models.MessagesDTO{
		Page:     options.Page,
		PerPage:  options.PerPage,
		Messages: []*models.Message{},
	}

	offset := (options.Page - 1) * options.PerPage
	query := r.db.Model(&models.Message{})
	if options.ChatID != 0 {
		query = query.Where("chat_id = ? AND user_id = ?", options.ChatID, userId)
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
