package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"bitbucket.org/iwlab-standuply/slackteams-api/auth"
	log "github.com/sirupsen/logrus"
)

var ErrNotFound = errors.New("GetAuthorization: no documents in result")

type GetAuthorization struct {
	Repo AuthorizationsRepository
}

type resultSingle struct {
	OK   bool                   `json:"ok"`
	Auth *SlackBotAuthorization `json:"auth"`
}

func (h GetAuthorization) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	v := r.URL.Query()
	teamId := v["teamId"]
	if len(teamId) == 0 {
		respond(w, errorJSON("server error - no param"), http.StatusInternalServerError)
		return
	}

	auth, err := h.Repo.GetAuthorization(ctx, teamId[0])
	log.WithContext(ctx).Debugf("auth: %+v\n", auth)

	if err != nil && err != ErrNotFound {
		respond(w, errorJSON("server error - DB request failed"), http.StatusInternalServerError)
		return
	}

	res := resultSingle{
		OK:   err == nil,
		Auth: auth,
	}

	resp, err := json.Marshal(res)
	if err != nil {
		respond(w, errorJSON("server error - JSON failed"), http.StatusInternalServerError)
		return
	}
	log.WithContext(ctx).Debugf("resp size: %d\n", len(resp))

	respond(w, resp, http.StatusOK)
}
