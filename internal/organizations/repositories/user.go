package repositories

import (
	"context"
	model "github.com/huyhvq/backend-dev-test/internal/organizations/models"
	"github.com/uptrace/bun"
)

type UserManager interface {
	Create(ctx context.Context, m *model.User) error
	Save(ctx context.Context, m *model.User) error
	FindByID(ctx context.Context, id uint64) (*model.User, error)
}

type user struct {
	db *bun.DB
}

func NewUser(db *bun.DB) *user {
	return &user{
		db: db,
	}
}

func (sq *user) Create(ctx context.Context, u *model.User) error {
	_, err := sq.db.NewInsert().Model(u).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (sq *user) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	var m model.User
	if err := sq.db.NewSelect().Model(&m).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return &m, nil
}

func (sq *user) Save(ctx context.Context, m *model.User) error {
	if _, err := sq.db.NewUpdate().Model(m).Where("id = ?", m.ID).Exec(ctx); err != nil {
		return err
	}
	return nil
}
