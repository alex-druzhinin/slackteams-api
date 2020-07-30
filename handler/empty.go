package handler

import "net/http"

type Empty struct{}

func (h Empty) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		respond(w, errorJSON("only GET requests are supported"), http.StatusMethodNotAllowed)
		return
	}

	w.Write(emptyPage)
}

var emptyPage = []byte(`
<!DOCTYPE html>
<html>
	<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
		Hello there
	</body>
</html>
`)
