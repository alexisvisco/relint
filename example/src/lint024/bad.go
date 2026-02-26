package handler

// AuthBodyRequest does not match {Name}BodyInput or {Name}BodyOutput
type AuthBodyRequest struct{} // want `LINT-024: body type "AuthBodyRequest" must be named "{Name}BodyInput" or "{Name}BodyOutput"`

// AuthBodyInput is valid
type AuthBodyInput struct{} // ok

// AuthBodyOutput is valid
type AuthBodyOutput struct{} // ok
