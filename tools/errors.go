package tools

import (
	"errors"
	"fmt"
)

// define constant for error management.
const (
	ClientError                       = "Client Error"
	ResourceError                     = "Resource Error"
	DataSourceError                   = "Data Source Error"
	UnexpectedImportIdentifier        = "Unexpected Import Identifier"
	UnexpectedResourceConfigureType   = "Unexpected Resource Configure Type"
	UnexpectedDataSourceConfigureType = "Unexpected DataSource Configure Type"
)

var ErrDataNotFound = errors.New("data source not found")

func ErrDataNotFoundError(kind, field, search string) error {
	return fmt.Errorf("%w: no %s with %s '%s'", ErrDataNotFound, kind, field, search)
}

func UnableToRead(name string, err error) string {
	return fmt.Sprintf("Unable to read %s, got error: %s", name, err)
}

func WrongClient(clientType string, providerData interface{}) string {
	return fmt.Sprintf("Expected %s, got: %T. Please report this issue to the provider developers.", clientType, providerData)
}
