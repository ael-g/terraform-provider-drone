package main

import (
        "github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
        return &schema.Provider{
                ResourcesMap: map[string]*schema.Resource{
                        "drone_activated_repository": droneActivatedRepository(),
                        "drone_secret": droneSecret(),
                },
        }
}
