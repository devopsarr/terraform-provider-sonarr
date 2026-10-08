package provider

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var errResourceNotInState = errors.New("resource not found in state")

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"sonarr": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	t.Helper()

	if v := os.Getenv("SONARR_URL"); v == "" {
		t.Skip("SONARR_URL must be set for acceptance tests")
	}

	if v := os.Getenv("SONARR_API_KEY"); v == "" {
		t.Skip("SONARR_API_KEY must be set for acceptance tests")
	}
}

func testAccAPIClient() *sonarr.APIClient {
	config := sonarr.NewConfiguration()
	config.AddDefaultHeader("X-Api-Key", os.Getenv("SONARR_API_KEY"))
	config.Servers[0].URL = os.Getenv("SONARR_URL")

	return sonarr.NewAPIClient(config)
}

const testUnauthorizedProvider = `
provider "sonarr" {
	url = "http://localhost:8989"
	api_key = "ErrorAPIKey"
	extra_headers = [
		{
			name = "testName"
			value = "testValue"
		}
	]
  }
`

// testAccCheckResourceDisappears deletes a resource outside Terraform, so that the next plan has to create it again.
func testAccCheckResourceDisappears(name string, deleteByID func(*sonarr.APIClient, int32) (*http.Response, error)) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("%w: %s", errResourceNotInState, name)
		}

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 32)
		if err != nil {
			return err
		}

		_, err = deleteByID(testAccAPIClient(), int32(id))

		return err
	}
}
