package main

import (
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
	name := d.Get("name").(string)
	value := d.Get("value").(string)

	client := drone.NewClient(host, auther)

	secret := drone.Secret{
		Name: "jlkjlkj",
		Value : "qsdqdsd",
	}

	_ = secret
	_ = client
	_ = repository
	_ = value
	_, err := client.SecretCreate(repository, name, nil)
	if err != nil {
		return err
	}

	//d.SetId(name)

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
	return nil
}
