package userhandler

import huma "github.com/alexisvisco/relint/example/src/github.com/danielgtaylor/huma/v2"

type UserHandler struct{}

func (h *UserHandler) GetDefinition() {}

type GetDefinitionInput struct {
	InvitationToken string `path:"invitation_token"` // want `LINT-031: path tag "invitation_token" must be lowerCamelCase`
}

func registerRoutes() {
	var group any
	h := &UserHandler{}

	huma.Get(group, "/api/objects/{object_id}/definition", h.GetDefinition) // want `LINT-031: huma path param "object_id" must be lowerCamelCase`
	huma.Get(group, "/api/objects/{ObjectId}/definition", h.GetDefinition)  // want `LINT-031: huma path param "ObjectId" must be lowerCamelCase`
	huma.Get(group, "/api/objects/{objectId}/{tenant_id}", h.GetDefinition) // want `LINT-031: huma path param "tenant_id" must be lowerCamelCase`
}
