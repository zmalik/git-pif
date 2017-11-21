package utils

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
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


func TestCreateFork(t *testing.T) {

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
