package main

import (
        "errors"
	"strings"
)

func splitRepoName(repo string) (string, string, error) {
	repositoryNameParts := strings.Split(repo, "/")
	if len(repositoryNameParts) != 2 {
		return "", "", errors.New("repo name must be 'org/name'")
	}

	return repositoryNameParts[0], repositoryNameParts[1], nil
}
