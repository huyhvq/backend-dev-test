package repositories

import (
	"context"
	model "github.com/huyhvq/backend-dev-test/internal/organizations/models"
	"github.com/uptrace/bun"
)

type TeamManager interface {
	Create(ctx context.Context, team *model.Team) error
	FindByID(ctx context.Context, id uint64) (*model.Team, error)
	Save(ctx context.Context, t *model.Team) error
}

var _ TeamManager = &team{}

type team struct {
	db *bun.DB
}

func NewTeam(db *bun.DB) *team {
	return &team{
		db: db,
	}
}

func (sq *team) Create(ctx context.Context, team *model.Team) error {
	_, err := sq.db.NewInsert().Model(team).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (sq *team) FindByID(ctx context.Context, id uint64) (*model.Team, error) {
	var t model.Team
	if err := sq.db.NewSelect().Model(&t).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return &t, nil
}

func (sq *team) Save(ctx context.Context, t *model.Team) error {
	if _, err := sq.db.NewUpdate().Model(t).Where("id = ?", t.ID).Exec(ctx); err != nil {
		return err
	}
	return nil
}
