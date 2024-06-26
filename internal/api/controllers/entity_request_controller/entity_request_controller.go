package entity_request_controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/go-chi/chi/v5"
)

type ResponseHandler interface {
	HandleResponseObject(w http.ResponseWriter, response interface{}, err error)
}

type EntityService[e any] interface {
	GetAll() ([]*e, error)
	GetById(id uint) (*e, error)
	Create(entity *e) (*e, error)
	Update(entity *e) (*e, error)
	Delete(id uint) error
}

type EntityRequestController[e any] struct {
	responseHandler ResponseHandler
	entityService   EntityService[e]
}

func NewEntityRequestController[e any](responseHandler ResponseHandler, entityService EntityService[e]) *EntityRequestController[e] {
	return &EntityRequestController[e]{
		responseHandler: responseHandler,
		entityService:   entityService,
	}
}

func (h *EntityRequestController[e]) HandleGetAllEntitiesRequest(w http.ResponseWriter, r *http.Request) {
	responseObject, err := h.entityService.GetAll()
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}

func (h *EntityRequestController[e]) HandleCreateEntityRequest(w http.ResponseWriter, r *http.Request) {
	var entity e
	err := json.NewDecoder(r.Body).Decode(&entity)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, models.ErrInvalidResourceSyntax)
		return
	}
	responseObject, err := h.entityService.Create(&entity)
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}

func (h *EntityRequestController[e]) HandleGetEntityByIdRequest(w http.ResponseWriter, r *http.Request) {
	entityId, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, fmt.Errorf("%w: %s", models.ErrInvalidId, err.Error()))
		return
	}
	responseObject, err := h.entityService.GetById(uint(entityId))
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}

func (h *EntityRequestController[e]) HandleUpdateEntityRequest(w http.ResponseWriter, r *http.Request) {
	var entity e
	err := json.NewDecoder(r.Body).Decode(&entity)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, models.ErrInvalidResourceSyntax)
		return
	}
	responseObject, err := h.entityService.Update(&entity)
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}

func (h *EntityRequestController[e]) HandleDeleteEntityByIdRequest(w http.ResponseWriter, r *http.Request) {
	entityId, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, fmt.Errorf("%w: %s", models.ErrInvalidId, err.Error()))
		return
	}
	err = h.entityService.Delete(uint(entityId))
	h.responseHandler.HandleResponseObject(w, nil, err)
}
