package grafana_v2

import (
	"context"
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// testAccProviderFactories is a static map containing only the main provider instance
var testAccProviderFactories map[string]func() (*schema.Provider, error)

// testAccProvider is the "main" provider instance
//
// This Provider can be used in testing code for API calls without requiring
// the use of saving and referencing specific ProviderFactories instances.
//
// testAccPreCheck(t) must be called before using this provider instance.
var testAccProvider *schema.Provider

// testAccProviderConfigure ensures that testAccProvider is only configured once.
//
// The testAccPreCheck(t) function is invoked for every test and this prevents
// extraneous reconfiguration to the same values each time. However, this does
// not prevent reconfiguration that may happen should the address of
// testAccProvider be errantly reused in ProviderFactories.
var testAccProviderConfigure sync.Once

func init() {
	testAccProvider = Provider("testacc")()

	// Always allocate a new provider instance each invocation, otherwise gRPC
	// ProviderConfigure() can overwrite configuration during concurrent testing.
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"grafana": func() (*schema.Provider, error) {
			return Provider("testacc")(), nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

// testAccPreCheckEnv contains all environment variables that must be present
// for acceptance tests to run. These are checked in testAccPreCheck.
var testAccPreCheckEnv = []string{
	"GRAFANA_URL",
	"GRAFANA_AUTH",
	"GRAFANA_ORG_ID",
}

// testAccPreCheck verifies required provider testing configuration. It should
// be present in every acceptance test.
//
// These verifications and configuration are preferred at this level to prevent
// provider developers from experiencing less clear errors for every test.
func testAccPreCheck(t *testing.T) {
	for _, e := range testAccPreCheckEnv {
		if v := os.Getenv(e); v == "" {
			t.Fatal(e + " must be set for acceptance tests")
		}
	}
	testAccProviderConfigure.Do(func() {
		// Since we are outside the scope of the Terraform configuration we must
		// call Configure() to properly initialize the provider configuration.
		err := testAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil))
		if err != nil {
			t.Fatal(err)
		}
	})
}

// testAccPreCheckCloud should be called by acceptance tests in files where the
// "cloud" build tag is present.
func testAccPreCheckCloud(t *testing.T) {
	testAccPreCheckEnv = append(testAccPreCheckEnv, "GRAFANA_SM_ACCESS_TOKEN")
	testAccPreCheck(t)
}

// testAccExample returns an example config from the examples directory.
// Examples are used for both documentation and acceptance tests.
func testAccExample(t *testing.T, path string) string {
	example, err := ioutil.ReadFile("../examples/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return string(example)
}
