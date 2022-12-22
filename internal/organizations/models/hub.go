package model

import (
	"github.com/uptrace/bun"
	"time"
)

var UnAssignedHub uint64 = 1

type Hub struct {
	bun.BaseModel
	ID        uint64    `json:"id" bun:"id,pk,autoincrement"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
