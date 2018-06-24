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
			"hooks": &schema.Schema{
				Type:     schema.TypeList,
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
				Type:     schema.TypeList,
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

func resourceActivatedRepositoryCreate(d *schema.ResourceData, m interface{}) error {
        client := getDroneClient()

	repoFullName := d.Get("name").(string)

        owner, repoName, err := splitRepoName(repoFullName)
        if err != nil {
                return err
        }

	_, err = client.RepoPost(owner, repoName)
	if err != nil {
		return err
	}

	repoPatch := drone.RepoPatch{}

	isTrusted, ok := d.GetOk("is_trusted")
	if ok {
		isTrustedTmp := isTrusted.(bool)
		repoPatch.IsTrusted = &isTrustedTmp
	}

	timeout, ok := d.GetOk("timeout")
	if ok {
		timeoutTmp := timeout.(int64)
		repoPatch.Timeout = &timeoutTmp
	}

	repo, err := client.RepoPatch(owner, repoName, &repoPatch)
	if err != nil {
		return err
	}

	d.Set("is_trusted", repo.IsTrusted)
	d.Set("timeout", repo.Timeout)

	d.SetId(repoFullName)

	return nil
}

func resourceActivatedRepositoryRead(d *schema.ResourceData, m interface{}) error {
	client := getDroneClient()

	repoFullName := d.Get("name").(string)

        owner, repoName, err := splitRepoName(repoFullName)
        if err != nil {
                return err
        }

	repoList, err := client.RepoList()
	if err != nil {
		return err
	}

	notFound := true
	for _, repo := range repoList {
		if repo.Name == repoName && repo.Owner == owner {
			notFound = false
			//d.Set("is_trusted", repo.IsTrusted)
			//d.Set("timeout", repo.Timeout)
		}
	}

	if notFound {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceActivatedRepositoryUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceActivatedRepositoryDelete(d *schema.ResourceData, m interface{}) error {
	client := getDroneClient()

        owner, repoName, err := splitRepoName(d.Get("name").(string))
        if err != nil {
                return err
        }

	err = client.RepoDel(owner, repoName)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

