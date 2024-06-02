package models

type User struct {
	BaseEntity
	Email string `json:"email"`
}

type UserScopedEntityInterface interface {
	BaseEntityInterface
	SetUserId(userId uint)
	GetUserId() uint
}

type UserScopedEntityInterfacePointer[T any] interface {
	*T
	UserScopedEntityInterface
}

type UserScopedEntity struct {
	BaseEntity
	UserId uint
}

func (e *UserScopedEntity) SetUserId(userId uint) {
	e.UserId = userId
}

func (e *UserScopedEntity) GetUserId() uint {
	return e.UserId
}
