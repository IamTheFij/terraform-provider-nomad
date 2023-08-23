// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nomad

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNodesRead,

		Schema: map[string]*schema.Schema{
			"prefix": {
				Description: "Specifies a string to filter nodes based on a name prefix.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"filter": {
				Description: "Specifies the expression used to filter the results.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"nodes": {
				Description: "List of nodes returned",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Description: "Address for this node.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"id": {
							Description: "ID for this node.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"datacenter": {
							Description: "Datacenter for this node.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "Unique name for this node.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"node_class": {
							Description: "Node class for this node.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"node_pool": {
							Description: "Node pool for this node.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func dataSourceNodesRead(d *schema.ResourceData, meta any) error {
	client := meta.(ProviderConfig).client

	prefix := d.Get("prefix").(string)
	filter := d.Get("filter").(string)
	id := strconv.Itoa(schema.HashString(prefix + filter))

	log.Printf("[DEBUG] Reading node list")
	resp, _, err := client.Nodes().List(&api.QueryOptions{
		Prefix: prefix,
		Filter: filter,
	})
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading nodes: %w", err)
	}

	nodes := make([]map[string]any, len(resp))
	for i, p := range resp {
		nodes[i] = map[string]any{
			"address":    p.Address,
			"id":         p.ID,
			"datacenter": p.Datacenter,
			"name":       p.Name,
			"node_class": p.NodeClass,
			"node_pool":  p.NodePool,
		}
	}
	log.Printf("[DEBUG] Read node list")

	d.SetId(id)
	return d.Set("nodes", nodes)
}
