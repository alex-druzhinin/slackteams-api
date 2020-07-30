package handler

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/iwlab-standuply/slackteams-api/auth"
)

type BotInfo struct {
	BotUserId      string
	BotAccessToken string
}

type SlackBotAuthorization struct {
	AccessToken string
	Scope       string
	UserId      string
	TeamName    string
	TeamId      string
	CreatedAt   string
	Enabled     bool
	Bot         BotInfo

	// These fields are not used in bot
	// cloned *string
}

type AllAuthorizations struct {
	Repo AuthorizationsRepository
}

func (h AllAuthorizations) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		respond(w, errorJSON("only GET requests are supported"), http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	user, ok := ctx.Value(auth.CtxKeyAuthUser).(string)
	if !ok || len(user) == 0 {
		respond(w, errorJSON("server error"), http.StatusInternalServerError)
		return
	}

	auths, err := h.Repo.GetAllAuthorizations(ctx)

	if err != nil {
		respond(w, errorJSON("server error - DB request failed"), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(auths)
	if err != nil {
		respond(w, errorJSON("server error - JSON failed"), http.StatusInternalServerError)
		return
	}

	respond(w, resp, http.StatusOK)
}
