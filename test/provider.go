package test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"test_scenario": dataSourceTestScenario(),
		},
		ResourcesMap: map[string]*schema.Resource{},
	}
}
