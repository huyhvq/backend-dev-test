package repositories

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"reflect"
	"testing"
)

func Test_organization_Search(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer sqlDB.Close()
	db := bun.NewDB(sqlDB, pgdialect.New())

	mock.ExpectQuery(`SELECT hub.id AS hub_id, hub.name AS hub_name, hub.location AS hub_location, t.id AS team_id, t.name AS team_name, t.type AS team_type FROM "hubs" AS "hub" FULL JOIN teams AS t ON t.hub_id = hub.id WHERE \(to_tsvector\(concat_ws\(' ', hub.name, hub.location, t.name, t.type\)\) @@ to_tsquery\('backend'\)\)`).
		WillReturnRows(sqlmock.NewRows([]string{"hub_id", "hub_name", "hub_location", "team_id", "team_name", "team_type"}).
			AddRow(1, "hub test", "SG", 2, "team name", "backend"))

	mock.ExpectQuery(`SELECT hub.id AS hub_id, hub.name AS hub_name, hub.location AS hub_location, t.id AS team_id, t.name AS team_name, t.type AS team_type FROM "hubs" AS "hub" FULL JOIN teams AS t ON t.hub_id = hub.id WHERE \(to_tsvector\(concat_ws\(' ', hub.name, hub.location, t.name, t.type\)\) @@ to_tsquery\('frontend'\)\)`).
		WillReturnRows(sqlmock.NewRows([]string{"hub_id", "hub_name", "hub_location", "team_id", "team_name", "team_type"}))

	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx  context.Context
		term string
	}
	w := []*SearchResult{
		{
			HubID:       1,
			HubName:     "hub test",
			HubLocation: "SG",
			TeamID:      2,
			TeamName:    "team name",
			TeamType:    "backend",
		},
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*SearchResult
		wantErr bool
	}{
		{
			name: "found",
			fields: fields{
				db: db,
			},
			args: args{
				ctx:  context.Background(),
				term: "backend",
			},
			want:    w,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sq := &organization{
				db: tt.fields.db,
			}
			got, err := sq.Search(tt.args.ctx, tt.args.term)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Search() got = %v, want %v", got, tt.want)
			}
		})
	}
}
