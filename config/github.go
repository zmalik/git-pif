package config

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

const ENV_GITHUB_VAR = "GITHUB_TOKEN"

func getGithubToken() string {
	token := viper.Get(ENV_GITHUB_VAR)
	if token != nil {
		return viper.Get(ENV_GITHUB_VAR).(string)
	}
	return ""
}

// Creates a github client using env variable GITHUB_TOKEN
func GetGithubClient(ctx context.Context) (*github.Client, error) {
	token := getGithubToken()
	if len(token) == 0 {
		return nil, errors.New(fmt.Sprintf("No %s found.", ENV_GITHUB_VAR))
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc), nil
}
