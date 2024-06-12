package xenserver

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccNetworkResourceConfig(name_label string) string {
	return fmt.Sprintf(`
resource "xenserver_network" "test_network" {
	name_label = "%s"
	name_description = "Network 0 for DHCP"
	other_config = {}
}
`, name_label)
}

func TestAccNetworkResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + testAccNetworkResourceConfig("test network 1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("xenserver_network.test_network", "name_label", "test network 1"),
					resource.TestCheckResourceAttr("xenserver_network.test_network", "name_description", "Network 0 for DHCP"),
					resource.TestCheckResourceAttr("xenserver_network.test_network", "other_config.%", "0"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("xenserver_network.test_network", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "xenserver_network.test_network",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{},
			},
			// Update and Read testing
			{
				Config: providerConfig + testAccNetworkResourceConfig("test network 2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("xenserver_network.test_network", "name_label", "test network 2"),
					resource.TestCheckResourceAttr("xenserver_network.test_network", "other_config.%", "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
