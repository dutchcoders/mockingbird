package mockingbird

import (
	"encoding/json"
	"net/http"
)

var slackOauthAccess = Handle("//slack.com:443/api/oauth.access", SlackOAuthAccessHandler)
var slackOauthTest = Handle("//slack.com:443/api/auth.test", SlackOAuthTestHandler)

func SlackOAuthAccessHandler(w ResponseWriter, req *http.Request) {
	w.WriteString("HTTP/1.1" + " 200 OK\r\n\r\n")

	bla := map[string]interface{}{
		"access_token": "valid-token",
		"scope":        "read",
	}
	json.NewEncoder(w).Encode(bla)
}

func SlackOAuthTestHandler(w ResponseWriter, req *http.Request) {
	w.WriteString("HTTP/1.1" + " 200 OK\r\n\r\n")

	bla := map[string]interface{}{
		"ok":      true,
		"url":     "https://myteam.slack.com/",
		"team":    "My Team",
		"user":    "cal",
		"team_id": "T12345",
		"user_id": "U12345",
	}
	json.NewEncoder(w).Encode(bla)
}
