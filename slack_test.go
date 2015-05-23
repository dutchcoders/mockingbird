package mockingbird

import (
	"testing"

	"github.com/dutchcoders/slack"
	. "github.com/smartystreets/goconvey/convey"
	slack_oauth "github.com/tappleby/slack_auth_proxy/slack"
)

func Test(t *testing.T) {
	Convey("slack", t, func() {
		client := slack_oauth.NewOAuthClient("", "", "")

		Convey("redeem code", func() {
			NewServer()

			accessToken, _ := client.RedeemCode("test")

			So(accessToken.Token, ShouldEqual, "valid-token")
		})

		Convey("auth.test", func() {
			NewServer()

			api := slack.New("")
			response, _ := api.AuthTest()

			So(response, ShouldEqual, &slack.AuthTestResponse{Url: "https://myteam.slack.com/", Team: "My Team", User: "cal", TeamId: "T12345", UserId: "U12345"})
		})

	})
}
