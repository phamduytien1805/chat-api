package interfaces

import (
	"github.com/google/uuid"
)

type CreateDMForm struct {
	Recipient uuid.UUID `json:"recipient" validate:"required"`
}
