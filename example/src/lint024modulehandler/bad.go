package authhandler

type AuthBodyRequest struct{} // want `LINT-024: body type "AuthBodyRequest" must be named "{Name}BodyInput" or "{Name}BodyOutput"`

type AuthBodyInput struct{}

type AuthBodyOutput struct{}
