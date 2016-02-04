package board_gamers

import (
	"net/http"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"encoding/json"
)

type Tweet struct {
	UserName	string
	Text		string
	LinkToTweet	string
	CreatedAt	string
}

func init() {
	http.HandleFunc("/hello", handler)
	http.HandleFunc("/webhook/trickplay", trickplayHandler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}

func trickplayHandler(w http.ResponseWriter, r *http.Request)  {
	ctx := appengine.NewContext(r)

	decoder := json.NewDecoder(r.Body)
	var t Tweet
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(w, "json parse error", 500)
		log.Errorf(ctx, "json parse error: %v", err)
	}
	log.Infof(ctx, "%v %v %v %v", t.UserName, t.Text, t.LinkToTweet, t.CreatedAt)
}
