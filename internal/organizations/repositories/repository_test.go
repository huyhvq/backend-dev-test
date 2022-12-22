package repositories

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	sqlDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := bun.NewDB(sqlDB, pgdialect.New())
	defer sqlDB.Close()

	h := NewHub(db)
	o := NewOrganization(db)
	tm := NewTeam(db)
	u := NewUser(db)

	r := &repo{
		hub:          h,
		organization: o,
		team:         tm,
		user:         u,
	}

	type args struct {
		db *bun.DB
	}
	tests := []struct {
		name             string
		args             args
		want             Manager
		wantHub          HubManager
		wantOrganization OrganizationManager
		wantTeam         TeamManager
		wantUser         UserManager
	}{
		{
			name: "succeed",
			args: args{
				db: db,
			},
			want:             r,
			wantHub:          h,
			wantOrganization: o,
			wantTeam:         tm,
			wantUser:         u,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
			if got := New(tt.args.db); !reflect.DeepEqual(got.Hub(), tt.wantHub) {
				t.Errorf("New() = %v, wantHub %v", got.Hub(), tt.wantHub)
			}
			if got := New(tt.args.db); !reflect.DeepEqual(got.Organization(), tt.wantOrganization) {
				t.Errorf("New() = %v, wantOrganization %v", got.Organization(), tt.wantOrganization)
			}
			if got := New(tt.args.db); !reflect.DeepEqual(got.Team(), tt.wantTeam) {
				t.Errorf("New() = %v, wantTeam %v", got.Team(), tt.wantTeam)
			}
			if got := New(tt.args.db); !reflect.DeepEqual(got.User(), tt.wantUser) {
				t.Errorf("New() = %v, wantUser %v", got.User(), tt.wantUser)
			}
		})
	}
}
