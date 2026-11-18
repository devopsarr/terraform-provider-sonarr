package helpers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/devopsarr/sonarr-go/sonarr"
)

// define constant for error management.
const (
	Create                            = "create"
	Read                              = "read"
	Update                            = "update"
	Delete                            = "delete"
	List                              = "list"
	ClientError                       = "Client Error"
	ResourceError                     = "Resource Error"
	DataSourceError                   = "Data Source Error"
	UnexpectedImportIdentifier        = "Unexpected Import Identifier"
	UnexpectedResourceConfigureType   = "Unexpected Resource Configure Type"
	UnexpectedDataSourceConfigureType = "Unexpected DataSource Configure Type"
)

func ParseNotFoundError(kind, field, search string) string {
	return fmt.Sprintf("Unable to find %s, got error: data source not found: no %s with %s '%s'", kind, kind, field, search)
}

func WrongClient(clientType string, providerData interface{}) string {
	return fmt.Sprintf("Expected %s, got: %T. Please report this issue to the provider developers.", clientType, providerData)
}

func ParseClientError(action, name string, err error) string {
	if e, ok := err.(*sonarr.GenericOpenAPIError); ok {
		return fmt.Sprintf("Unable to %s %s, got error: %s\nDetails:\n%s", action, name, err, string(e.Body()))
	}

	return fmt.Sprintf("Unable to %s %s, got error: %s", action, name, err)
}

// IsNotFound reports whether the API answered 404 Not Found, meaning the object no longer exists.
// The client puts the HTTP status line ("404 Not Found") in the error; only the code is matched.
func IsNotFound(err error) bool {
	var apiErr *sonarr.GenericOpenAPIError

	return errors.As(err, &apiErr) && strings.HasPrefix(apiErr.Error(), "404 ")
}
