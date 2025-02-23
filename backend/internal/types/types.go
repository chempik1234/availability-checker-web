package types

import "fmt"

const (
	AuthorizationHeader = "Authorization"
	TokenPrefix         = "Token"
)

var (
	TokenPrefixErrorMessage = fmt.Sprintf("Token prefix must be equal to %s", AuthorizationHeader)
)
