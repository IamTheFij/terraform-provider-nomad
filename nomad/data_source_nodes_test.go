// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nomad

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataSourceNodes_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-nomad-test")
	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		PreCheck:  func() { testAccPreCheck(t); testCheckMinVersion(t, "1.6.0-beta.1") },
		Steps: []resource.TestStep{
			{
				Config: testDataSourceNodes_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nomad_nodes.all", "nodes.#", "1"),
					resource.TestCheckResourceAttr("data.nomad_nodes.prefix", "nodes.#", "1"),
					resource.TestCheckResourceAttr("data.nomad_nodes.filter", "nodes.#", "0"),
					resource.TestCheckResourceAttr("data.nomad_nodes.filter_with_prefix", "nodes.#", "0"),
				),
			},
		},
	})
}

func testDataSourceNodes_basic(prefix string) string {
	return fmt.Sprintf(`
data "nomad_nodes" "all" {
  depends_on = [
    nomad_node.basic,
    nomad_node.simple,
    nomad_node.different_prefix,
  ]
}

data "nomad_nodes" "prefix" {
  depends_on = [
    nomad_node.basic,
    nomad_node.simple,
    nomad_node.different_prefix,
  ]

  prefix = "%[1]s"
}

data "nomad_nodes" "filter" {
  depends_on = [
    nomad_node.basic,
    nomad_node.simple,
    nomad_node.different_prefix,
  ]

  filter = "Meta.test == \"%[1]s\""
}

data "nomad_nodes" "filter_with_prefix" {
  depends_on = [
    nomad_node.basic,
    nomad_node.simple,
    nomad_node.different_prefix,
  ]

  prefix = "%[1]s"
  filter = "Meta.test == \"%[1]s\""
}
`, prefix)
}
