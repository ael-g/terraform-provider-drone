package main

import (
	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/helper/schema"
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
				ForceNew: true,
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
	client := getDroneClient()

	owner, repoName, err := splitRepoName(d.Get("repository").(string))
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	value := d.Get("value").(string)
	eventsRaw := d.Get("events")
	events := []string{}
	for _, event := range eventsRaw.([]interface{}) {
		events = append(events, event.(string))
	}

	secret := drone.Secret{
		Name: name,
		Value: value,
		Events: events,
	}

	_, err = client.SecretCreate(owner, repoName, &secret)
	if err != nil {
		return err
	}

	d.SetId(name)

	return nil
}

func resourceSecretRead(d *schema.ResourceData, m interface{}) error {
	client := getDroneClient()

	owner, repoName, err := splitRepoName(d.Get("repository").(string))
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	secret, err := client.Secret(owner, repoName, name)
	if err != nil {
		d.SetId("")
		return nil
	}

	d.Set("events", secret.Events)

	return nil
}

func resourceSecretUpdate(d *schema.ResourceData, m interface{}) error {
	client := getDroneClient()

	owner, repoName, err := splitRepoName(d.Get("repository").(string))
	if err != nil {
		return err
	}

	secret := drone.Secret{
		Name: d.Get("name").(string),
	}

	if d.HasChange("value") {
		secret.Value = d.Get("value").(string)
	}

	if d.HasChange("events") {
		eventsRaw := d.Get("events")
		events := []string{}
		for _, event := range eventsRaw.([]interface{}) {
			events = append(events, event.(string))
		}
		secret.Events = events
	}

	_, err = client.SecretUpdate(owner, repoName, &secret)
	if err != nil {
		return err
	}

	return nil
}

func resourceSecretDelete(d *schema.ResourceData, m interface{}) error {
	client := getDroneClient()

	name := d.Get("name").(string)
	owner, repoName, err := splitRepoName(d.Get("repository").(string))
	if err != nil {
		return err
	}

	err = client.SecretDelete(owner, repoName, name)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
