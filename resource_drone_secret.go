package main

import (
	"strconv"

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
				ForceNew: true,
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
	client := m.(drone.Client)

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

	secretRemote, err := client.SecretCreate(owner, repoName, &secret)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(secretRemote.ID, 10))

	return nil
}

func resourceSecretRead(d *schema.ResourceData, m interface{}) error {
	client := m.(drone.Client)

	owner, repoName, err := splitRepoName(d.Get("repository").(string))
	if err != nil {
		return err
	}

	Id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	secrets, err := client.SecretList(owner, repoName)
	if err != nil {
		return err
	}

        for _, secret := range secrets {
                if secret.ID == Id {
			d.Set("events", secret.Events)
			return nil
                }
        }

	d.SetId("")
	return nil
}

func resourceSecretUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(drone.Client)

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
	client := m.(drone.Client)

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
