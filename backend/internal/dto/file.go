package dto

import (
	"github.com/oapi-codegen/runtime/types"
	"time"

	"github.com/google/uuid"
)

type File struct {
	Uuid         uuid.UUID `json:"uuid"`
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	CreatedAt    time.Time `json:"created_at"`
	SymmetricKey *string   `json:"symmetric_key,omitempty"`
}

type CreateFileInput struct {
	Name         string     `json:"name"`
	SymmetricKey *string    `json:"symmetric_key,omitempty"`
	File         types.File `json:"file"`
}

type ShareFileInput struct {
	RecipientUUID uuid.UUID `json:"recipient_uuid"`
	SymmetricKey  string    `json:"symmetric_key"`
}

type DeleteFileAccessInput struct {
	OwnerUUID     uuid.UUID `json:"owner_uuid"`
	FileUUID      uuid.UUID `json:"file_uuid"`
	RecipientUUID uuid.UUID `json:"recipient_uuid"`
}
