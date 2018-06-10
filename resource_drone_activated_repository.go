package main

import (
	"log"
	"errors"
	"strings"
        "github.com/hashicorp/terraform/helper/schema"
	"github.com/drone/drone-go/drone"
	"golang.org/x/oauth2"
)

const (
	token = ""
	host  = ""
)

func droneActivatedRepository() *schema.Resource {
        return &schema.Resource{
                Create: resourceServerCreate,
                Read:   resourceServerRead,
                Update: resourceServerUpdate,
                Delete: resourceServerDelete,

                Schema: map[string]*schema.Schema{
                        "name": &schema.Schema{
                                Type:     schema.TypeString,
                                Required: true,
                        },
                        "hooks": &schema.Schema{
                                Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
                                Optional: true,
                        },
			"is_protected": &schema.Schema{
                                Type:     schema.TypeBool,
                                Optional: true,
                        },
			"is_trusted": &schema.Schema{
                                Type:     schema.TypeBool,
                                Optional: true,
                        },
			"visibility": &schema.Schema{
                                Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
                                Optional: true,
                        },
			"timeout": &schema.Schema{
                                Type:     schema.TypeInt,
                                Optional: true,
                        },
                },
        }
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	config := new(oauth2.Config)
	auther := config.Client(
		oauth2.NoContext,
		&oauth2.Token{
			AccessToken: token,
		},
	)
        name := d.Get("name").(string)
	nameParts := strings.Split(name, "/")
	if len(nameParts) != 2 {
		return errors.New("repo name must be 'org/name'")
	}

	client := drone.NewClient(host, auther)

	user, err := client.Self()
	log.Println(user, err)

	_, err = client.RepoPost(nameParts[0], nameParts[1])
	if err != nil {
		return err
	}

        d.SetId(name)

        return nil
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
        return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
        return nil
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	config := new(oauth2.Config)
	auther := config.Client(
		oauth2.NoContext,
		&oauth2.Token{
			AccessToken: token,
		},
	)

        name := d.Get("name").(string)
        nameParts := strings.Split(name, "/")
        if len(nameParts) != 2 {
                return errors.New("repo name must be 'org/name'")
        }

	client := drone.NewClient(host, auther)

	user, err := client.Self()
	log.Println(user, err)

	err = client.RepoDel(nameParts[0], nameParts[1])
	if err != nil {
		return err
	}
        d.SetId("")
        return nil
}
