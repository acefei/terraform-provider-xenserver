package xenserver

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccVmDataSourceConfig(name_label string) string {
	return fmt.Sprintf(`
data "xenserver_vm" "test_vm_data" {
	name_label = "%s"
}
`, name_label)
}

func TestAccVmDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + testAccVmDataSourceConfig("virtual machine 0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.xenserver_vm.test_vm_data", "name_label", "virtual machine 0"),
				),
			},
		},
	})
}
