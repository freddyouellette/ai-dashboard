package user_scoped_request_controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/freddyouellette/ai-dashboard/plugins/plugin_models"
	"github.com/go-chi/chi/v5"
)

type ResponseHandler interface {
	HandleResponseObject(w http.ResponseWriter, response interface{}, err error)
}

type UserScopedService[e any] interface {
	Delete(id uint) error
	GetAllWithUserId(userId uint) ([]e, error)
	CreateWithUserId(entity e, userId uint) (e, error)
	UpdateWithUserId(entity e, userId uint) (e, error)
	GetByIdAndUserId(id uint, userId uint) (e, error)
	DeleteFromUserId(id uint, userId uint) error
}

type RequestUtils interface {
	GetContextInt(r *http.Request, key any, def int) int
}

type UserScopedRequestController[e models.UserScopedEntityInterface] struct {
	responseHandler ResponseHandler
	entityService   UserScopedService[e]
	requestUtils    RequestUtils
}

func NewUserScopedRequestController[e models.UserScopedEntityInterface](
	responseHandler ResponseHandler,
	entityService UserScopedService[e],
	requestUtils RequestUtils,
) *UserScopedRequestController[e] {
	return &UserScopedRequestController[e]{
		responseHandler: responseHandler,
		entityService:   entityService,
		requestUtils:    requestUtils,
	}
}

func (h *UserScopedRequestController[e]) HandleGetAllEntitiesRequest(w http.ResponseWriter, r *http.Request) {
	userId := h.requestUtils.GetContextInt(r, plugin_models.UserIdContextKey{}, 0)
	responseObject, err := h.entityService.GetAllWithUserId(uint(userId))
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}

func (h *UserScopedRequestController[e]) HandleCreateEntityRequest(w http.ResponseWriter, r *http.Request) {
	var entity e
	err := json.NewDecoder(r.Body).Decode(&entity)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, models.ErrInvalidResourceSyntax)
		return
	}
	userId := h.requestUtils.GetContextInt(r, plugin_models.UserIdContextKey{}, 0)
	responseObject, err := h.entityService.CreateWithUserId(entity, uint(userId))
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}

func (h *UserScopedRequestController[e]) HandleGetEntityByIdRequest(w http.ResponseWriter, r *http.Request) {
	entityId, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, fmt.Errorf("%w: %s", models.ErrInvalidId, err.Error()))
		return
	}
	userId := h.requestUtils.GetContextInt(r, plugin_models.UserIdContextKey{}, 0)
	responseObject, err := h.entityService.GetByIdAndUserId(uint(entityId), uint(userId))
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}

func (h *UserScopedRequestController[e]) HandleUpdateEntityRequest(w http.ResponseWriter, r *http.Request) {
	var entity e
	err := json.NewDecoder(r.Body).Decode(&entity)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, models.ErrInvalidResourceSyntax)
		return
	}
	userId := h.requestUtils.GetContextInt(r, plugin_models.UserIdContextKey{}, 0)
	responseObject, err := h.entityService.UpdateWithUserId(entity, uint(userId))
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}

func (h *UserScopedRequestController[e]) HandleDeleteEntityByIdRequest(w http.ResponseWriter, r *http.Request) {
	entityId, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, fmt.Errorf("%w: %s", models.ErrInvalidId, err.Error()))
		return
	}
	err = h.entityService.Delete(uint(entityId))
	h.responseHandler.HandleResponseObject(w, nil, err)
}
