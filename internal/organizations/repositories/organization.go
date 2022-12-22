package repositories

import (
	"context"
	model "github.com/huyhvq/backend-dev-test/internal/organizations/models"
	"github.com/uptrace/bun"
)

type OrganizationManager interface {
	Search(ctx context.Context, term string) ([]*SearchResult, error)
}

type organization struct {
	db *bun.DB
}

func NewOrganization(db *bun.DB) *organization {
	return &organization{
		db: db,
	}
}

type SearchResult struct {
	HubID       uint64 `json:"hub_id"`
	HubName     string `json:"hub_name"`
	HubLocation string `json:"hub_location"`
	TeamID      uint64 `json:"team_id"`
	TeamName    string `json:"team_name"`
	TeamType    string `json:"team_type"`
}

func (sq *organization) Search(ctx context.Context, term string) ([]*SearchResult, error) {
	var rows []*SearchResult
	q := sq.db.NewSelect().Model((*model.Hub)(nil)).
		ColumnExpr("hub.id AS hub_id, hub.name AS hub_name, hub.location AS hub_location").
		ColumnExpr("t.id AS team_id, t.name AS team_name, t.type AS team_type").
		Join("FULL JOIN teams AS t ON t.hub_id = hub.id").
		Where("to_tsvector(concat_ws(' ', hub.name, hub.location, t.name, t.type)) @@ to_tsquery(?)", term)
	if err := q.Scan(ctx, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}
