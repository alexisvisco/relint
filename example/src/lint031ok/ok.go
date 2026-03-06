package userhandler

import hapi "github.com/alexisvisco/relint/example/src/github.com/danielgtaylor/huma/v2"

type UserHandler struct{}

func (h *UserHandler) GetDefinition() {}

type GetDefinitionInput struct {
	InvitationToken string `path:"invitationToken"`
}

func registerRoutes() {
	var group any
	h := &UserHandler{}

	hapi.Get(group, "/api/objects/{objectId}/definition", h.GetDefinition)
	hapi.Post(group, "/api/objects/{objectId}/tenant/{tenantId}", h.GetDefinition)
}
