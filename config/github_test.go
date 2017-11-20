package config

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"context"
	"os"
	"github.com/spf13/viper"
)

func TestGetGithubClient(t *testing.T) {

	Convey("Get github client with no GITHUB_TOKEN", t, func() {
		_, err := GetGithubClient(context.Background())

		So(err, ShouldNotBeNil)
	})

	Convey("Get github client with a GITHUB_TOKEN", t, func() {
		viper.AutomaticEnv()
		os.Unsetenv("GITHUB_TOKEN")
		os.Setenv("GITHUB_TOKEN", "value")
		_, err := GetGithubClient(context.Background())

		So(err, ShouldBeNil)
	})
}
