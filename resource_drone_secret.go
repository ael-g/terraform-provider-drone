package main

import (
	"log"

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
			"pull_request": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"pull_request_push": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  true,
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
	pullRequest := d.Get("pull_request").(bool)
	pullRequestPush := d.Get("pull_request_push").(bool)

	secret := drone.Secret{
		Name:            name,
		Data:            value,
		PullRequest:     pullRequest,
		PullRequestPush: pullRequestPush,
	}

	_, err = client.SecretCreate(owner, repoName, &secret)
	if err != nil {
		return err
	}

	d.SetId(repoName + "/" + name)

	return nil
}

func resourceSecretRead(d *schema.ResourceData, m interface{}) error {
	client := m.(drone.Client)
	owner, repoName, err := splitRepoName(d.Get("repository").(string))
	if err != nil {
		return err
	}

	secrets, err := client.SecretList(owner, repoName)
	if err != nil {
		return err
	}
	for _, secret := range secrets {
		log.Println(secret)
		if repoName+"/"+secret.Name == d.Id() {
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
		secret.Data = d.Get("value").(string)
	}

	if d.HasChange("pull_request") {
		secret.PullRequest = d.Get("pull_request").(bool)
	}

	if d.HasChange("pull_request_push") {
		secret.PullRequestPush = d.Get("pull_request_push").(bool)
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
