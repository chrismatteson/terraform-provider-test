package github

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	p := &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"test_scenario":           dataSourceTestScenario(),
		},
	}
}
