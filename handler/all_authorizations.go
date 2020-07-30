package handler

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/iwlab-standuply/slackteams-api/auth"
	log "github.com/sirupsen/logrus"
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

type result struct {
	OK    bool                     `json:"ok"`
	Auths []*SlackBotAuthorization `json:"auths"`
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
	log.WithContext(ctx).Debugf("auths size: %d\n", len(auths))
	log.WithContext(ctx).Debugf("auth: %+v\n", auths[0])

	if err != nil {
		respond(w, errorJSON("server error - DB request failed"), http.StatusInternalServerError)
		return
	}

	res := result{
		OK:    true,
		Auths: removeDuplicates(auths),
	}

	resp, err := json.Marshal(res)
	if err != nil {
		respond(w, errorJSON("server error - JSON failed"), http.StatusInternalServerError)
		return
	}
	log.WithContext(ctx).Debugf("resp size: %d\n", len(resp))

	respond(w, resp, http.StatusOK)
}

func removeDuplicates(elements []*SlackBotAuthorization) []*SlackBotAuthorization {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := make([]*SlackBotAuthorization, len(elements))
	i := 0

	for v := range elements {
		if encountered[elements[v].TeamId] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v].TeamId] = true
			// Append to result slice.
			result[i] = elements[v]
			i++
		}
	}
	// Return the new slice.
	return result
}
