package userhandler

type Operation struct{}

func (o Operation) WithPattern(_ string) Operation { return o }

func registerRoutes() {
	Operation{}.WithPattern("GET /api/objects/{objectId}/definition")
	Operation{}.WithPattern("POST /api/objects/{objectId}/tenant/{tenantId}")
}

type GetDefinitionInput struct {
	InvitationToken string `path:"invitationToken"`
}
