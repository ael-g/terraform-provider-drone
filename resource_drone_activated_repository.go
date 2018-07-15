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
				Default: true,
				Optional: true,
			},
			"visibility": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"allow_pull": &schema.Schema{
				Type:     schema.TypeBool,
				Default: true,
				Optional: true,
			},
			"allow_push": &schema.Schema{
				Type:     schema.TypeBool,
				Default: true,
				Optional: true,
			},
			"allow_tag": &schema.Schema{
				Type:     schema.TypeBool,
				Default: false,
				Optional: true,
			},
			"allow_deployment": &schema.Schema{
				Type:     schema.TypeBool,
				Default: false,
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
	if repo != nil && !repo.IsActive {
		_, err = client.RepoPost(owner, repoName)
		if err != nil {
			return err
		}
	}
	repoPatch := drone.RepoPatch{}

	isTrusted, ok := d.GetOk("is_trusted")
	if ok {
		isTrustedTmp := isTrusted.(bool)
		repoPatch.IsTrusted = &isTrustedTmp
	}

	allowPull, ok := d.GetOk("allow_pull")
	if ok {
		allowPullTmp := allowPull.(bool)
		repoPatch.AllowPull = &allowPullTmp
	}

	allowPush, ok := d.GetOk("allow_push")
	if ok {
		allowPushTmp := allowPush.(bool)
		repoPatch.AllowPush = &allowPushTmp
	}

	allowTag, ok := d.GetOk("allow_tag")
	if ok {
		allowTagTmp := allowTag.(bool)
		repoPatch.AllowTag = &allowTagTmp
	}

	allowDeploy, ok := d.GetOk("allow_deploy")
	if ok {
		allowDeployTmp := allowDeploy.(bool)
		repoPatch.AllowDeploy = &allowDeployTmp
	}

	timeout, ok := d.GetOk("timeout")
	if ok {
		timeoutTmp := int64(timeout.(int))
		repoPatch.Timeout = &timeoutTmp
	}

	repo, err = client.RepoPatch(owner, repoName, &repoPatch)
	if err != nil {
		return err
	}

	d.Set("is_trusted", repo.IsTrusted)
	d.Set("timeout", repo.Timeout)
	d.Set("allow_pull", repo.AllowPull)
	d.Set("allow_push", repo.AllowPush)
	d.Set("allow_tag", repo.AllowTag)
	d.Set("allow_deploy", repo.AllowDeploy)

	d.SetId(repoFullName)

	return nil
}

func resourceActivatedRepositoryRead(d *schema.ResourceData, m interface{}) error {
	client := m.(drone.Client)

	repoFullName := d.Get("name").(string)

        owner, repoName, err := splitRepoName(repoFullName)
        if err != nil {
                return err
        }

	repoList, err := client.RepoList()
	if err != nil {
		return err
	}

	for _, repo := range repoList {
		if repo.Name == repoName && repo.Owner == owner {
			d.Set("is_trusted", repo.IsTrusted)
			d.Set("timeout", repo.Timeout)
			d.Set("allow_pull", repo.AllowPull)
			d.Set("allow_push", repo.AllowPush)
			d.Set("allow_tag", repo.AllowTag)
			d.Set("allow_deploy", repo.AllowDeploy)
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

	isTrusted, ok := d.GetOk("is_trusted")
	if ok {
		isTrustedTmp := isTrusted.(bool)
		repoPatch.IsTrusted = &isTrustedTmp
	}

	allowPull, ok := d.GetOk("allow_pull")
	if ok {
		allowPullTmp := allowPull.(bool)
		repoPatch.AllowPull = &allowPullTmp
	}

	allowPush, ok := d.GetOk("allow_push")
	if ok {
		allowPushTmp := allowPush.(bool)
		repoPatch.AllowPush = &allowPushTmp
	}

	allowTag, ok := d.GetOk("allow_tag")
	if ok {
		allowTagTmp := allowTag.(bool)
		repoPatch.AllowTag = &allowTagTmp
	}

	allowDeploy, ok := d.GetOk("allow_deploy")
	if ok {
		allowDeployTmp := allowDeploy.(bool)
		repoPatch.AllowDeploy = &allowDeployTmp
	}

	timeout, ok := d.GetOk("timeout")
	if ok {
		timeoutTmp := int64(timeout.(int))
		repoPatch.Timeout = &timeoutTmp
	}

	repo, err := client.RepoPatch(owner, repoName, &repoPatch)
	if err != nil {
		return err
	}

	d.Set("is_trusted", repo.IsTrusted)
	d.Set("timeout", repo.Timeout)
	d.Set("allow_pull", repo.AllowPull)
	d.Set("allow_push", repo.AllowPush)
	d.Set("allow_tag", repo.AllowTag)
	d.Set("allow_deploy", repo.AllowDeploy)

	return nil
}

func resourceActivatedRepositoryDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(drone.Client)

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

