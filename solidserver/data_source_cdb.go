package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/url"
)

func dataSourcecdb() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcecdbRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the custom DB.",
				Required:    true,
			},
			"label1": {
				Type:        schema.TypeString,
				Description: "The name of the label 1",
				Computed:    true,
			},
			"label2": {
				Type:        schema.TypeString,
				Description: "The name of the label 2",
				Computed:    true,
			},
			"label3": {
				Type:        schema.TypeString,
				Description: "The name of the label 3",
				Computed:    true,
			},
			"label4": {
				Type:        schema.TypeString,
				Description: "The name of the label 4",
				Computed:    true,
			},
			"label5": {
				Type:        schema.TypeString,
				Description: "The name of the label 5",
				Computed:    true,
			},
			"label6": {
				Type:        schema.TypeString,
				Description: "The name of the label 6",
				Computed:    true,
			},
			"label7": {
				Type:        schema.TypeString,
				Description: "The name of the label 7",
				Computed:    true,
			},
			"label8": {
				Type:        schema.TypeString,
				Description: "The name of the label 8",
				Computed:    true,
			},
			"label9": {
				Type:        schema.TypeString,
				Description: "The name of the label 9",
				Computed:    true,
			},
			"label10": {
				Type:        schema.TypeString,
				Description: "The name of the label 10",
				Computed:    true,
			},
		},
	}
}

func dataSourcecdbRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)
	d.SetId("")

	// Building parameters
	parameters := url.Values{}
	parameters.Add("WHERE", "name='"+d.Get("name").(string)+"'")

	// Sending the read request
	resp, body, err := s.Request("get", "rest/custom_db_name_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.SetId(buf[0]["custom_db_name_id"].(string))

			d.Set("name", buf[0]["name"].(string))
			d.Set("label1", buf[0]["label1"].(string))
			d.Set("label2", buf[0]["label2"].(string))
			d.Set("label3", buf[0]["label3"].(string))
			d.Set("label4", buf[0]["label4"].(string))
			d.Set("label5", buf[0]["label5"].(string))
			d.Set("label6", buf[0]["label6"].(string))
			d.Set("label7", buf[0]["label7"].(string))
			d.Set("label8", buf[0]["label8"].(string))
			d.Set("label9", buf[0]["label9"].(string))
			d.Set("label10", buf[0]["label10"].(string))

			return nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to read information from custom DB: %s (%s)\n", d.Get("name").(string), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to read information from custom DB: %s\n", d.Get("name").(string))
		}

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find custom DB: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}
