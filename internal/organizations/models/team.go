package model

import (
	"github.com/uptrace/bun"
	"time"
)

var UnAssignedTeam uint64 = 1

type Team struct {
	bun.BaseModel
	ID        uint64    `json:"id" bun:"id,pk,autoincrement"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	HubID     uint64    `json:"hub_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
