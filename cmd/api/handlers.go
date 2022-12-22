package api

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	model "github.com/huyhvq/backend-dev-test/internal/organizations/models"
	"github.com/huyhvq/backend-dev-test/pkg/request"
	"github.com/huyhvq/backend-dev-test/pkg/response"
	"github.com/huyhvq/backend-dev-test/pkg/validator"
	"golang.org/x/exp/maps"
	"net/http"
	"strconv"
)

type organizationSearchHub struct {
	ID       uint64                    `json:"id"`
	Name     string                    `json:"name"`
	Location string                    `json:"location"`
	Teams    []*organizationSearchTeam `json:"teams"`
}

type organizationSearchTeam struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type createHubRequest struct {
	Name      string              `json:"name"`
	Location  string              `json:"location"`
	Validator validator.Validator `json:"-"`
}

type createTeamRequest struct {
	Name      string              `json:"name"`
	Type      string              `json:"type"`
	HubID     uint64              `json:"hub_id"`
	Validator validator.Validator `json:"-"`
}

type joinToHubRequest struct {
	HubID     uint64              `json:"hub_id"`
	Validator validator.Validator `json:"-"`
}

type createUserRequest struct {
	Name      string              `json:"name"`
	Title     string              `json:"title"`
	TeamID    uint64              `json:"team_id"`
	Validator validator.Validator `json:"-"`
}

type joinToTeamRequest struct {
	TeamID    uint64              `json:"team_id"`
	Validator validator.Validator `json:"-"`
}

func (app *application) status(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Status": "OK",
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) organizationSearch(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")
	if term == "" {
		app.errorMessage(w, r, http.StatusUnprocessableEntity, "search term is required", nil)
		return
	}
	teams, err := app.repositories.Organization().Search(r.Context(), term)
	if err != nil {
		app.serverError(w, r, err)
	}

	hms := make(map[uint64]*organizationSearchHub, 0)
	for _, team := range teams {
		if _, ok := hms[team.HubID]; !ok {
			hms[team.HubID] = &organizationSearchHub{
				ID:       team.HubID,
				Name:     team.HubName,
				Location: team.HubLocation,
				Teams:    []*organizationSearchTeam{},
			}
		}
		if team.TeamID != 0 {
			hms[team.HubID].Teams = append(hms[team.HubID].Teams, &organizationSearchTeam{
				ID:   team.TeamID,
				Name: team.TeamName,
				Type: team.TeamType,
			})
		}
	}
	if err := response.JSON(w, http.StatusOK, map[string]interface{}{"hubs": maps.Values(hms)}); err != nil {
		app.serverError(w, r, err)
		return
	}
	return
}

func (app *application) createHub(w http.ResponseWriter, r *http.Request) {
	var params createHubRequest
	if err := request.DecodeJSON(w, r, &params); err != nil {
		app.badRequest(w, r, err)
		return
	}

	params.Validator.CheckField(validator.NotBlank(params.Name), "name", "name is required")
	params.Validator.CheckField(validator.NotBlank(params.Location), "location", "location is required")

	if params.Validator.HasErrors() {
		app.failedValidation(w, r, params.Validator)
		return
	}
	hub := model.Hub{
		Name:     params.Name,
		Location: params.Location,
	}
	if err := app.repositories.Hub().Create(r.Context(), &hub); err != nil {
		app.serverError(w, r, err)
		return
	}
	if err := response.JSON(w, http.StatusOK, hub); err != nil {
		app.serverError(w, r, err)
		return
	}
	return
}

func (app *application) createTeam(w http.ResponseWriter, r *http.Request) {
	var params createTeamRequest
	if err := request.DecodeJSON(w, r, &params); err != nil {
		app.badRequest(w, r, err)
		return
	}
	params.Validator.CheckField(validator.NotBlank(params.Name), "name", "name is required")
	params.Validator.CheckField(validator.NotBlank(params.Type), "type", "type is required")
	if params.Validator.HasErrors() {
		app.failedValidation(w, r, params.Validator)
		return
	}
	if params.HubID == 0 {
		params.HubID = model.UnAssignedHub
	}
	hub, err := app.repositories.Hub().FindByID(r.Context(), params.HubID)
	if err != nil {
		if err == sql.ErrNoRows {
			app.errorMessage(w, r, http.StatusUnprocessableEntity, "invalid hub_id", nil)
			return
		}
		app.serverError(w, r, err)
		return
	}
	team := model.Team{
		Name:  params.Name,
		Type:  params.Type,
		HubID: hub.ID,
	}
	if err := app.repositories.Team().Create(r.Context(), &team); err != nil {
		app.serverError(w, r, err)
		return
	}
	if err := response.JSON(w, http.StatusOK, team); err != nil {
		app.serverError(w, r, err)
		return
	}
	return
}

func (app *application) joinIntoHub(w http.ResponseWriter, r *http.Request) {
	var params joinToHubRequest
	if err := request.DecodeJSON(w, r, &params); err != nil {
		app.badRequest(w, r, err)
		return
	}
	var teamID uint64 = 0
	if id := chi.URLParam(r, "id"); id != "" {
		i, err := strconv.Atoi(id)
		if err != nil {
			app.errorMessage(w, r, http.StatusUnprocessableEntity, "invalid hub_id", nil)
			return
		}
		teamID = uint64(i)
	}

	team, err := app.repositories.Team().FindByID(r.Context(), teamID)
	if err != nil {
		if err == sql.ErrNoRows {
			app.errorMessage(w, r, http.StatusUnprocessableEntity, "invalid team_id", nil)
			return
		}
		app.serverError(w, r, err)
		return
	}
	hub, err := app.repositories.Hub().FindByID(r.Context(), params.HubID)
	if err != nil {
		if err == sql.ErrNoRows {
			app.errorMessage(w, r, http.StatusUnprocessableEntity, "invalid hub_id", nil)
			return
		}
		app.serverError(w, r, err)
		return
	}
	team.HubID = hub.ID
	if err := app.repositories.Team().Save(r.Context(), team); err != nil {
		app.serverError(w, r, err)
		return
	}
	if err := response.JSON(w, http.StatusOK, team); err != nil {
		app.serverError(w, r, err)
		return
	}
	return
}

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	var params createUserRequest
	if err := request.DecodeJSON(w, r, &params); err != nil {
		app.badRequest(w, r, err)
		return
	}
	params.Validator.CheckField(validator.NotBlank(params.Name), "name", "name is required")
	params.Validator.CheckField(validator.NotBlank(params.Title), "type", "type is required")
	if params.Validator.HasErrors() {
		app.failedValidation(w, r, params.Validator)
		return
	}
	if params.TeamID == 0 {
		params.TeamID = model.UnAssignedTeam
	}
	team, err := app.repositories.Hub().FindByID(r.Context(), params.TeamID)
	if err != nil {
		if err == sql.ErrNoRows {
			app.errorMessage(w, r, http.StatusUnprocessableEntity, "invalid team_id", nil)
			return
		}
		app.serverError(w, r, err)
		return
	}
	user := model.User{
		Name:   params.Name,
		Title:  params.Title,
		TeamID: team.ID,
	}
	if err := app.repositories.User().Create(r.Context(), &user); err != nil {
		app.serverError(w, r, err)
		return
	}
	if err := response.JSON(w, http.StatusOK, user); err != nil {
		app.serverError(w, r, err)
		return
	}
	return
}

func (app *application) joinIntoTeam(w http.ResponseWriter, r *http.Request) {
	var params joinToTeamRequest
	if err := request.DecodeJSON(w, r, &params); err != nil {
		app.badRequest(w, r, err)
		return
	}
	var userID uint64 = 0
	if id := chi.URLParam(r, "id"); id != "" {
		i, err := strconv.Atoi(id)
		if err != nil {
			app.errorMessage(w, r, http.StatusUnprocessableEntity, "invalid team_id", nil)
			return
		}
		userID = uint64(i)
	}

	team, err := app.repositories.Team().FindByID(r.Context(), params.TeamID)
	if err != nil {
		if err == sql.ErrNoRows {
			app.errorMessage(w, r, http.StatusUnprocessableEntity, "invalid team_id", nil)
			return
		}
		app.serverError(w, r, err)
		return
	}
	user, err := app.repositories.User().FindByID(r.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			app.errorMessage(w, r, http.StatusUnprocessableEntity, "invalid user_id", nil)
			return
		}
		app.serverError(w, r, err)
		return
	}
	user.TeamID = params.TeamID
	if err := app.repositories.User().Save(r.Context(), user); err != nil {
		app.serverError(w, r, err)
		return
	}
	if err := response.JSON(w, http.StatusOK, team); err != nil {
		app.serverError(w, r, err)
		return
	}
	return
}
