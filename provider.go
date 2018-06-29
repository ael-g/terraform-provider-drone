package main

import (
        "github.com/hashicorp/terraform/helper/schema"
        "github.com/drone/drone-go/drone"
        "golang.org/x/oauth2"
)

func Provider() *schema.Provider {
        return &schema.Provider{
                Schema: map[string]*schema.Schema{
                        "host": {
                                Type:        schema.TypeString,
                                Required:    true,
                                DefaultFunc: schema.EnvDefaultFunc("DRONE_SERVER", "http://localhost/"),
                                Description: "The url of your Drone server",
                        },
			"token": {
                                Type:        schema.TypeString,
                                Required:    true,
                                DefaultFunc: schema.EnvDefaultFunc("DRONE_TOKEN", ""),
                                Description: "",
                        },
		},
                ResourcesMap: map[string]*schema.Resource{
                        "drone_activated_repository": droneActivatedRepository(),
                        "drone_secret": droneSecret(),
                },
                ConfigureFunc: providerConfigure,
        }
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
        hostInterface, ok := d.GetOk("host")
        var host, token string
        if ok {
                host = hostInterface.(string)
        }

        tokenInterface, ok := d.GetOk("token")
        if ok {
                token = tokenInterface.(string)
        }


        config := new(oauth2.Config)
        auther := config.Client(
                oauth2.NoContext,
                &oauth2.Token{
                        AccessToken: token,
                },
        )

        return drone.NewClient(host, auther), nil
}
