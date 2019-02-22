package main

import (
	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/helper/schema"
)

func droneActivatedRepository() *schema.Resource {
	return &schema.Resource{
		Create: resourceActivatedRepositoryCreate,
		Read:   resourceActivatedRepositoryRead,
		Update: resourceActivatedRepositoryUpdate,
		Delete: resourceActivatedRepositoryDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"is_protected": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_trusted": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"visibility": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"allow_pull": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"allow_push": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"allow_tag": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"allow_deploy": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Default:  60,
				Optional: true,
			},
		},
	}
}

func resourceActivatedRepositoryCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(drone.Client)

	repoFullName := d.Get("name").(string)

	owner, repoName, err := splitRepoName(repoFullName)
	if err != nil {
		return err
	}

	// Check if repo is already
	repo, err := client.Repo(owner, repoName)
	if repo != nil && !repo.Active {
		_, err = client.RepoEnable(owner, repoName)
		if err != nil {
			return err
		}
	}
	repoPatch := drone.RepoPatch{}

	timeout, ok := d.GetOk("timeout")
	if ok {
		timeoutTmp := int64(timeout.(int))
		repoPatch.Timeout = &timeoutTmp
	}

	repo, err = client.RepoUpdate(owner, repoName, &repoPatch)
	if err != nil {
		return err
	}

	d.Set("timeout", repo.Timeout)

	d.SetId(repoFullName)

	return nil
}

func resourceActivatedRepositoryRead(d *schema.ResourceData, m interface{}) error {
	client := m.(drone.Client)

	repoFullName := d.Get("name").(string)

	_, repoName, err := splitRepoName(repoFullName)
	if err != nil {
		return err
	}

	repoList, err := client.RepoList()
	if err != nil {
		return err
	}

	for _, repo := range repoList {
		if repo.Name == repoName {
			d.Set("timeout", repo.Timeout)
			return nil
		}
	}

	d.SetId("")

	return nil
}

func resourceActivatedRepositoryUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(drone.Client)

	repoFullName := d.Get("name").(string)

	owner, repoName, err := splitRepoName(repoFullName)
	if err != nil {
		return err
	}

	repoPatch := drone.RepoPatch{}

	timeout, ok := d.GetOk("timeout")
	if ok {
		timeoutTmp := int64(timeout.(int))
		repoPatch.Timeout = &timeoutTmp
	}

	repo, err := client.RepoUpdate(owner, repoName, &repoPatch)
	if err != nil {
		return err
	}

	d.Set("timeout", repo.Timeout)

	return nil
}

func resourceActivatedRepositoryDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(drone.Client)

	owner, repoName, err := splitRepoName(d.Get("name").(string))
	if err != nil {
		return err
	}

	err = client.RepoDisable(owner, repoName)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
