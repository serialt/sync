package service

import (
	"context"

	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	Token string
}

func NewClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client
}
