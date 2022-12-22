package model

import (
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel
	ID        uint64    `json:"id" bun:"id,pk,autoincrement"`
	Name      string    `json:"name"`
	Title     string    `json:"type"`
	TeamID    uint64    `json:"team_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
