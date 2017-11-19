package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"syscall"
)

const FORK_UPSTREAM = "fork"

func Push(owner, repo, user string) (err error) {
	branch, err := execGitConfig(true, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return err
	}
	upstream, err := execGitConfig(true, "config", "--get", fmt.Sprintf("branch.%s.remote", branch))
	if err != nil {
		_, err := execGitConfig(true, "config", "--get", fmt.Sprintf("remote.%s.url", FORK_UPSTREAM))
		if err != nil {
			upstreamUrl := fmt.Sprintf("https://github.com/%s/%s", user, repo)
			_, err := execGitConfig(true, "remote", "add", FORK_UPSTREAM, upstreamUrl)
			if err != nil {
				return errors.New(fmt.Sprintf("Error adding the upstream %s in the repo.", upstreamUrl))
			}
		}
	} else {
		remoteUrl, err := execGitConfig(true, "config", "--get", fmt.Sprintf("remote.%s.url", upstream))
		parsedOwner, parsedRepo, err := ParseUpstreamURL(remoteUrl)
		if err != nil {
			return err
		}
		if owner == parsedOwner && repo == parsedRepo {
			return errors.New(fmt.Sprintf("The branch %s is already following an upstream in %s/%s."+
				" Create a new branch.", branch, owner, repo))
		}

		if user != parsedOwner && repo == parsedRepo {
			return errors.New(fmt.Sprintf("The branch %s is already following a different upstream in %s/%s."+
				" Create a new branch.", branch, parsedOwner, repo))
		}
	}

	output, err := execGitConfig(false, "push", FORK_UPSTREAM, "-u", branch)
	fmt.Printf(output)
	return err
}

func execGitConfig(trim bool, args ...string) (string, error) {
	var stdout bytes.Buffer
	cmd := exec.Command("git", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = ioutil.Discard

	err := cmd.Run()
	if exitError, ok := err.(*exec.ExitError); ok {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			if waitStatus.ExitStatus() == 1 {
				return "", errors.New(fmt.Sprintf("Command not found: %s", args[len(args)-1]))
			}
		}
		return "", err
	}
	if trim {
		return strings.TrimRight(stdout.String(), "\n"), nil
	}
	return stdout.String(), nil
}
