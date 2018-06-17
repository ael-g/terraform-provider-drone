package main

import (
	"strings"
	"fmt"
	"errors"
	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/helper/schema"
	"golang.org/x/oauth2"
)

func droneSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecretCreate,
		Read:   resourceSecretRead,
		Update: resourceSecretUpdate,
		Delete: resourceSecretDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"repository": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
                        "events": &schema.Schema{
                                Type:     schema.TypeList,
                                Elem:     &schema.Schema{Type: schema.TypeString},
                                Optional: true,
                        },
		},
	}
}

func resourceSecretCreate(d *schema.ResourceData, m interface{}) error {
	config := new(oauth2.Config)
	auther := config.Client(
		oauth2.NoContext,
		&oauth2.Token{
			AccessToken: token,
		},
	)
	repository := d.Get("repository").(string)
        repositoryNameParts := strings.Split(repository, "/")
        if len(repositoryNameParts) != 2 {
                return errors.New("repo name must be 'org/name'")
        }

	name := d.Get("name").(string)
	value := d.Get("value").(string)
	events := d.Get("events").([]string)
	fmt.Println(events)

	client := drone.NewClient(host, auther)

	secret := drone.Secret{
		Name: name,
		Value: value,
	//	Events: events,
	}

	_, err := client.SecretCreate(repositoryNameParts[0], repositoryNameParts[1], &secret)
	if err != nil {
		return err
	}

	d.SetId(name)

	return nil
}

func resourceSecretRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSecretUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSecretDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	config := new(oauth2.Config)
	auther := config.Client(
		oauth2.NoContext,
		&oauth2.Token{
			AccessToken: token,
		},
	)
	repository := d.Get("repository").(string)
        repositoryNameParts := strings.Split(repository, "/")
        if len(repositoryNameParts) != 2 {
                return errors.New("repo name must be 'org/name'")
        }

	name := d.Get("name").(string)

	client := drone.NewClient(host, auther)

	_, err := client.Secret(repositoryNameParts[0], repositoryNameParts[1], name)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
