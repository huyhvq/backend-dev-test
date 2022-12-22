package repositories

import "github.com/uptrace/bun"

type Manager interface {
	Hub() HubManager
	Organization() OrganizationManager
	Team() TeamManager
	User() UserManager
}

var _ Manager = &repo{}

type repo struct {
	hub          *hub
	organization *organization
	team         *team
	user         *user
}

func New(db *bun.DB) *repo {
	return &repo{
		hub:          NewHub(db),
		organization: NewOrganization(db),
		team:         NewTeam(db),
		user:         NewUser(db),
	}
}

func (r repo) Hub() HubManager {
	return r.hub
}

func (r repo) Organization() OrganizationManager {
	return r.organization
}

func (r repo) Team() TeamManager {
	return r.team
}

func (r repo) User() UserManager {
	return r.user
}
