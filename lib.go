package main

import (
        "errors"
	"strings"
        "github.com/drone/drone-go/drone"
        "golang.org/x/oauth2"
)

func getDroneClient() drone.Client {
        config := new(oauth2.Config)
        auther := config.Client(
                oauth2.NoContext,
                &oauth2.Token{
                        AccessToken: token,
                },
        )

        return drone.NewClient(host, auther)
}

func splitRepoName(repo string) (string, string, error) {
	repositoryNameParts := strings.Split(repo, "/")
	if len(repositoryNameParts) != 2 {
		return "", "", errors.New("repo name must be 'org/name'")
	}

	return repositoryNameParts[0], repositoryNameParts[1], nil
}
