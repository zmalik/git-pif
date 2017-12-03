package utils

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"context"
	"github.com/zmalik/git-pif/config"
	"gopkg.in/jarcoal/httpmock.v1"
	"os"
	"github.com/spf13/viper"
	"time"
)

func TestParseUpstreamURL(t *testing.T) {

	Convey("Given a valid github repo url", t, func() {

		Convey("https://github.com/octagen/anyrepo.git", func() {
			owner, repo, err := ParseUpstreamURL("https://github.com/octagen/anyrepo.git")
			So(err, ShouldBeNil)
			So(owner, ShouldEqual, "octagen")
			So(repo, ShouldEqual, "anyrepo")
		})
		Convey("ssh://git@github.com/octagen/anyrepo.git", func() {
			owner, repo, err := ParseUpstreamURL("ssh://git@github.com/octagen/anyrepo.git")
			So(err, ShouldBeNil)
			So(owner, ShouldEqual, "octagen")
			So(repo, ShouldEqual, "anyrepo")
		})

		Convey("https://github.com/octagen/anyrepo", func() {
			owner, repo, err := ParseUpstreamURL("https://github.com/octagen/anyrepo")
			So(err, ShouldBeNil)
			So(owner, ShouldEqual, "octagen")
			So(repo, ShouldEqual, "anyrepo")
		})

		Convey("ssh://git@github.com/octagen/anyrepo", func() {
			owner, repo, err := ParseUpstreamURL("ssh://git@github.com/octagen/anyrepo")
			So(err, ShouldBeNil)
			So(owner, ShouldEqual, "octagen")
			So(repo, ShouldEqual, "anyrepo")
		})
		Convey("git@github.com:octagen/anyrepo", func() {
			owner, repo, err := ParseUpstreamURL("git@github.com:octagen/anyrepo")
			So(err, ShouldBeNil)
			So(owner, ShouldEqual, "octagen")
			So(repo, ShouldEqual, "anyrepo")
		})
		Convey("git@github.com:octagen/anyrepo.git", func() {
			owner, repo, err := ParseUpstreamURL("git@github.com:octagen/anyrepo.git")
			So(err, ShouldBeNil)
			So(owner, ShouldEqual, "octagen")
			So(repo, ShouldEqual, "anyrepo")
		})

	})
	Convey("Given an invalid github repo url", t, func() {
		Convey("git@github.com:octagen.git", func() {
			_, _, err := ParseUpstreamURL("git@github.com:octagen.git")
			So(err, ShouldNotBeNil)
		})
		Convey("https://github.com/octagen.git", func() {
			_, _, err := ParseUpstreamURL("https://github.com/octagen.git")
			So(err, ShouldNotBeNil)
		})
		Convey("letsjusttest", func() {
			_, _, err := ParseUpstreamURL("letsjusttest")
			So(err, ShouldNotBeNil)
		})
	})

}

func TestForkRepo(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://api.github.com/repos/owner/repo/forks",
		httpmock.NewStringResponder(202, "{}"))
	Convey("Call github to fork the repo", t, func() {
		ctx := context.Background()
		viper.AutomaticEnv()
		os.Setenv(config.ENV_GITHUB_VAR, "fakeToken")
		client, err := config.GetGithubClient(ctx)
		So(err, ShouldBeNil)
		err = forkRepo(ctx, client, "owner", "repo")
		So(err, ShouldBeNil)
	})
	httpmock.RegisterResponder("POST", "https://api.github.com/repos/owner/repo/forks",
		httpmock.NewStringResponder(400, "{}"))
	Convey("Call github to fork the repo with auth error", t, func() {
		ctx := context.Background()
		viper.AutomaticEnv()
		os.Setenv(config.ENV_GITHUB_VAR, "fakeToken")
		client, err := config.GetGithubClient(ctx)
		So(err, ShouldBeNil)
		err = forkRepo(ctx, client, "owner", "repo")
		So(err, ShouldNotBeNil)
	})
	httpmock.RegisterResponder("POST", "https://api.github.com/repos/owner/repo/forks",
		httpmock.NewStringResponder(200, "{}"))
	Convey("Call github to fork the repo with no error but unexpected code", t, func() {
		ctx := context.Background()
		viper.AutomaticEnv()
		os.Setenv(config.ENV_GITHUB_VAR, "fakeToken")
		client, err := config.GetGithubClient(ctx)
		So(err, ShouldBeNil)
		err = forkRepo(ctx, client, "owner", "repo")
		So(err, ShouldBeNil)
	})
}

func TestCreateFork(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://api.github.com/repos/owner/repo/forks",
		httpmock.NewStringResponder(202, "{}"))
	httpmock.RegisterResponder("GET", "https://api.github.com/user",
		httpmock.NewStringResponder(200, "{\"login\":\"octagen\"}"))
	httpmock.RegisterResponder("POST", "https://api.github.com/repos/zmalik/git-pif/forks",
		httpmock.NewStringResponder(202, "{}"))
	viper.AutomaticEnv()
	os.Setenv(config.ENV_GITHUB_VAR, "fakeToken")
	Convey("Create fork of an existing fork", t, func() {
		httpmock.RegisterResponder("GET", "https://api.github.com/repos/octagen/git-pif",
			httpmock.NewStringResponder(200, "{}"))
		owner, repo, login, err := CreateFork()
		So(err, ShouldBeNil)
		So(owner, ShouldEqual, "zmalik")
		So(repo, ShouldEqual, "git-pif")
		So(login, ShouldEqual, "octagen")

	})
	Convey("Create fork for first time", t, func() {
		httpmock.RegisterResponder("GET", "https://api.github.com/repos/octagen/git-pif",
			httpmock.NewStringResponder(404, "{}"))
		go setRepoGetReponse(200)
		owner, repo, login, err := CreateFork()

		So(err, ShouldBeNil)
		So(owner, ShouldEqual, "zmalik")
		So(repo, ShouldEqual, "git-pif")
		So(login, ShouldEqual, "octagen")

	})

	Convey("Create fork with no github token", t, func() {
		os.Setenv(config.ENV_GITHUB_VAR, "")
		_, _, _, err := CreateFork()

		So(err, ShouldNotBeNil)

	})

	Convey("Create fork with invalid github token", t, func() {
		os.Setenv(config.ENV_GITHUB_VAR, "fake")
		httpmock.RegisterResponder("GET", "https://api.github.com/user",
			httpmock.NewStringResponder(403, "{\"login\":\"octagen\"}"))

		_, _, _, err := CreateFork()

		So(err, ShouldNotBeNil)
		httpmock.RegisterResponder("GET", "https://api.github.com/user",
			httpmock.NewStringResponder(200, "{\"login\":\"octagen\"}"))


	})

	Convey("Create fork with error in fork call", t, func() {
		setRepoGetReponse(404)
		httpmock.RegisterResponder("POST", "https://api.github.com/repos/zmalik/git-pif/forks",
			httpmock.NewStringResponder(400, "{}"))
		_, _, _, err := CreateFork()

		So(err, ShouldNotBeNil)
	})

	Convey("Create fork and timeout will polling", t, func() {
		setRepoGetReponse(404)
		httpmock.RegisterResponder("POST", "https://api.github.com/repos/zmalik/git-pif/forks",
			httpmock.NewStringResponder(202, "{}"))
		os.Setenv(timeoutPolling, "20ms")
		viper.AutomaticEnv()
		_, _, _, err := CreateFork()

		So(err, ShouldNotBeNil)
	})

}

func setRepoGetReponse(response int) {
	time.Sleep(1 * time.Second)
	httpmock.RegisterResponder("GET", "https://api.github.com/repos/octagen/git-pif",
		httpmock.NewStringResponder(response, "{}"))
}
