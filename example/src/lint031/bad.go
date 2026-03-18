package userhandler

type Operation struct{}

func (o Operation) WithPattern(_ string) Operation { return o }

func registerRoutes() {
	Operation{}.WithPattern("GET /api/objects/{object_id}/definition") // want `LINT-031: httpapi path param "object_id" must be lowerCamelCase`
	Operation{}.WithPattern("GET /api/objects/{ObjectId}/definition")  // want `LINT-031: httpapi path param "ObjectId" must be lowerCamelCase`
	Operation{}.WithPattern("GET /api/objects/{objectId}/{tenant_id}") // want `LINT-031: httpapi path param "tenant_id" must be lowerCamelCase`
}

type GetDefinitionInput struct {
	InvitationToken string `path:"invitation_token"` // want `LINT-031: path tag "invitation_token" must be lowerCamelCase`
}
