package repositories

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	model "github.com/huyhvq/backend-dev-test/internal/organizations/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"reflect"
	"testing"
	"time"
)

func TestNewHub(t *testing.T) {
	sqlDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer sqlDB.Close()
	db := bun.NewDB(sqlDB, pgdialect.New())
	type args struct {
		db *bun.DB
	}
	tests := []struct {
		name string
		args args
		want *hub
	}{
		{
			name: "succeed",
			args: args{
				db: db,
			},
			want: &hub{db: db},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHub(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hub_Create(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer sqlDB.Close()
	db := bun.NewDB(sqlDB, pgdialect.New())

	var errInsert = errors.New("test insert error")
	mock.ExpectQuery(`^INSERT INTO "hubs" \("id", "name", "location", "created_at", "updated_at"\) VALUES \(DEFAULT, 'hub_name', 'hub_location',(.*),(.*)\) RETURNING "id"$`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectQuery(`^INSERT INTO "hubs" \("id", "name", "location", "created_at", "updated_at"\) VALUES \(DEFAULT, 'hub_name_1', 'hub_location_1',(.*),(.*)\) RETURNING "id"$`).WillReturnError(errInsert)
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		hub *model.Hub
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "succeed",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: context.Background(),
				hub: &model.Hub{
					Name:     "hub_name",
					Location: "hub_location",
				},
			},
			wantErr: false,
		},
		{
			name: "failed",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: context.Background(),
				hub: &model.Hub{
					Name:     "hub_name_1",
					Location: "hub_location_1",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sq := &hub{
				db: tt.fields.db,
			}
			if err := sq.Create(tt.args.ctx, tt.args.hub); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_hub_FindByID(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer sqlDB.Close()
	db := bun.NewDB(sqlDB, pgdialect.New())

	errUnexpected := errors.New("test: unexpected")
	mock.ExpectQuery(`SELECT "hub"."id", "hub"."name", "hub"."location", "hub"."created_at", "hub"."updated_at" FROM "hubs" AS "hub" WHERE \(id = 1\)`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "location", "created_at", "updated_at"}).
			AddRow(1, "hub name", "SG", "2022-12-17 14:07:57", "0001-01-01 00:00:00"))
	mock.ExpectQuery(`SELECT "hub"."id", "hub"."name", "hub"."location", "hub"."created_at", "hub"."updated_at" FROM "hubs" AS "hub" WHERE \(id = 2\)`).WillReturnError(errUnexpected)

	ct, err := time.Parse("2006-01-02 15:04:05", "2022-12-17 14:07:57")
	if err != nil {
		t.Error(err)
	}
	ut, err := time.Parse("2006-01-02 15:04:05", "0001-01-01 00:00:00")
	if err != nil {
		t.Error(err)
	}
	mockHub := model.Hub{
		ID:        1,
		Name:      "hub name",
		Location:  "SG",
		CreatedAt: ct,
		UpdatedAt: ut,
	}
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		id  uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Hub
		wantErr bool
	}{
		{
			name: "succeed",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    &mockHub,
			wantErr: false,
		},
		{
			name: "unexpected error",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: context.Background(),
				id:  2,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sq := &hub{
				db: tt.fields.db,
			}
			got, err := sq.FindByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
