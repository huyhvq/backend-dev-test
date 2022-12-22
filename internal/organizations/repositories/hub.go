package repositories

import (
	"context"
	model "github.com/huyhvq/backend-dev-test/internal/organizations/models"
	"github.com/uptrace/bun"
)

type HubManager interface {
	Create(ctx context.Context, hub *model.Hub) error
	FindByID(ctx context.Context, id uint64) (*model.Hub, error)
}

type hub struct {
	db *bun.DB
}

func NewHub(db *bun.DB) *hub {
	return &hub{
		db: db,
	}
}

func (sq *hub) Create(ctx context.Context, hub *model.Hub) error {
	_, err := sq.db.NewInsert().Model(hub).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (sq *hub) FindByID(ctx context.Context, id uint64) (*model.Hub, error) {
	var m model.Hub
	if err := sq.db.NewSelect().Model(&m).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return &m, nil
}
