package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/tcnksm/go-gitconfig"
	"github.com/zmalik/git-pif/config"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func CreateFork() (owner, repo, userLogin string, err error) {
	origin, err := gitconfig.OriginURL()
	if err != nil {
		return "", "", "", errors.New("Not a git repository (or any of the parent directories): .git")
	}
	owner, repo, err = ParseUpstreamURL(origin)
	if err != nil {
		return "", "", "", err
	}
	ctx := context.Background()
	client, err := config.GetGithubClient(ctx)
	if err != nil {
		return "", "", "", err
	}
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return "", "", "", errors.New(fmt.Sprintf("Problem with the GITHUB_TOKEN: %s", err.Error()))
	}
	userLogin = user.GetLogin()
	fmt.Printf("Checking if user %s have a fork of %s/%s\n", user.GetLogin(), owner, repo)
	_, _, err = client.Repositories.Get(ctx, user.GetLogin(), repo)
	if err == nil {
		fmt.Printf("Fork %s/%s already exists\n", user.GetLogin(), repo)
	} else {
		fmt.Printf("Forking the repo %s/%s to %s/%s\n", owner, repo, user.GetLogin(), repo)
		err = ForkRepo(ctx, client, owner, repo)
		if err == nil {
			gitconfig.OriginURL()
			_, err := pollRepo(ctx, client, user.GetLogin(), repo, 5*time.Minute)
			if err == nil {
				fmt.Printf("Fork %s/%s created successfully.\n", user.GetLogin(), repo)
			} else {
				return "", "", "", err
			}
		} else {
			fmt.Printf("error forking the repo %s/%s to %s/%s\n. Error: %v", owner, repo, user.GetLogin(), repo, err)
		}
	}

	return owner, repo, user.GetLogin(), err
}

func pollRepo(ctx context.Context, client *github.Client, owner, repo string, duration time.Duration) (bool, error) {
	timeout := time.After(duration)
	tick := time.Tick(5 * time.Second)
	for {
		select {
		case <-timeout:
			return false, errors.New(fmt.Sprintf("timed out while checking the created fork %s/%s", owner, repo))
		case <-tick:
			_, _, err := client.Repositories.Get(ctx, owner, repo)
			fmt.Printf("Waiting the creation of the fork %s/%s.\n", owner, repo)
			if err == nil {
				return true, err
			}
		}
	}

}

//Fork the repo from the given owner and repo
func ForkRepo(ctx context.Context, client *github.Client, owner, repo string) error {
	_, resp, err := client.Repositories.CreateFork(ctx, owner, repo, &github.RepositoryCreateForkOptions{})
	if resp != nil && resp.StatusCode == 202 {
		return nil
	}
	if err != nil {
		msg := fmt.Sprintf("Error creating the fork %s/%s. Error: %v\n", owner, repo, err)
		return errors.New(msg)
	}
	return nil

}

func ParseUpstreamURL(remoteURL string) (owner, repo string, err error) {
	defer func() {
		// Strip trailing .git if present.
		if strings.HasSuffix(repo, ".git") {
			repo = repo[:len(repo)-len(".git")]
		}
	}()

	// Try to parse the URL as an SSH URL first.
	owner, repo, ok := tryParseUpstreamAsSSH(remoteURL)
	if ok {
		return owner, repo, nil
	}

	// Try to parse the URL as a regular URL.
	owner, repo, ok = tryParseUpstreamAsURL(remoteURL)
	if ok {
		return owner, repo, nil
	}

	// No success, return an error.
	err = fmt.Errorf("failed to parse git remote URL: %v", remoteURL)
	return "", "", err
}

// tryParseUpstreamAsSSH tries to parse the address as an SSH address,
// e.g. git@github.com:owner/repo.git
func tryParseUpstreamAsSSH(remoteURL string) (owner, repo string, ok bool) {
	re := regexp.MustCompile("^git@[^:]+:([^/]+)/(.+)$")
	match := re.FindStringSubmatch(remoteURL)
	if len(match) != 0 {
		owner, repo := match[1], match[2]
		return owner, repo, true
	}

	return "", "", false
}

// tryParseUpstreamAsURL tries to parse the address as a regular URL,
// e.g. https://github.com/owner/repo
func tryParseUpstreamAsURL(remoteURL string) (owner, repo string, ok bool) {
	u, err := url.Parse(remoteURL)
	if err != nil {
		return "", "", false
	}

	switch u.Scheme {
	case "ssh":
		fallthrough
	case "https":
		re := regexp.MustCompile("^/([^/]+)/(.+)$")
		match := re.FindStringSubmatch(u.Path)
		if len(match) != 0 {
			owner, repo := match[1], match[2]
			return owner, repo, true
		}
	}

	return "", "", false
}
