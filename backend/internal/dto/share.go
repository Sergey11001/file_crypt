package dto

import (
	"time"

	"github.com/google/uuid"
)

type UrlSharesView struct {
	UserUUID  uuid.UUID `json:"user_uuid"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateShareInput struct {
	UserUUID     uuid.UUID  `json:"user_uuid"`
	FileUUID     uuid.UUID  `json:"file_uuid"`
	SymmetricKey string     `json:"symmetric_key"`
	ViewLimit    *int64     `json:"view_limit,omitempty"`
	ExpiredAt    *time.Time `json:"expired_at,omitempty"`
}
