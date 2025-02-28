package catalog

import (
	"testing"

	"github.com/databricks/terraform-provider-databricks/qa"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestShareData(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.1/unity-catalog/shares/a?include_shared_data=true",
				Response: ShareInfo{
					Name: "a",
					Objects: []SharedDataObject{
						{
							Name:           "a",
							DataObjectType: "TABLE",
							Comment:        "c",
							SharedAs:       "",
							AddedAt:        0,
							AddedBy:        "",
						},
					},
					CreatedBy: "bob",
					CreatedAt: 1921321,
				},
			},
		},
		Resource:    DataSourceShare(),
		Read:        true,
		NonWritable: true,
		ID:          "_",
		HCL: `
		name = "a"
		`,
	}.Apply(t)
	assert.NoError(t, err, err)
	assert.Equal(t, "bob", d.Get("created_by"))
	assert.Equal(t, 1921321, d.Get("created_at"))
	assert.Equal(t,
		map[string]interface{}{
			"added_at":         0,
			"added_by":         "",
			"comment":          "c",
			"data_object_type": "TABLE",
			"name":             "a",
			"shared_as":        "",
		},
		d.Get("object").(*schema.Set).List()[0])
}

func TestShareData_Error(t *testing.T) {
	qa.ResourceFixture{
		Fixtures:    qa.HTTPFailures,
		Resource:    DataSourceShare(),
		Read:        true,
		NonWritable: true,
		ID:          "_",
	}.ExpectError(t, "I'm a teapot")
}
