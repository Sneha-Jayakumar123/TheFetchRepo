package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type RequestStruct struct {
	Token string
}

func GetRepository(w http.ResponseWriter, r *http.Request) {
	// Your GitHub Personal Access Token (PAT)
	var request RequestStruct

	json.NewDecoder(r.Body).Decode(&request)

	token := request.Token

	// Create an OAuth2 token source with the PAT
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	// Create an authenticated HTTP client
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)

	// Create a new GitHub client
	client := github.NewClient(oauthClient)

	// Replace "username" and "repositoryName" with the actual GitHub username and repository name
	username := "Sneha-Jayakumar123"
	repositoryName := "FirstRepo"

	// Get the repository
	repo, _, err := client.Repositories.Get(context.Background(), username, repositoryName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Repository Name: %s\n", *repo.Name)

}
