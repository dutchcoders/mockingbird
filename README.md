# Mockingbird

Mockingbird is a framework that mocks traffic to standard services. This can be used for testing api implementations. We want to create a framework that will mock all kind of different available services and extensible with own implementations.

[![Build Status](https://travis-ci.org/dutchcoders/mockingbird.svg?branch=master)](https://travis-ci.org/dutchcoders/mockingbird)

## Implementation

```
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
```


## Contributions

Contributions are welcome.

## Creators

**Remco Verhoef**
- <https://twitter.com/remco_verhoef>
- <https://twitter.com/dutchcoders>

**Gerred Dillon** 

## Copyright and license

Code and documentation copyright 2011-2015 Remco Verhoef.

Code released under [the MIT license](LICENSE).
