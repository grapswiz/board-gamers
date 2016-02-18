package board_gamers

import (
	"encoding/json"
	"github.com/mjibson/goon"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"net/http"
)

func ArrivalOfGamesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	g := goon.NewGoon(r)

	var a []ArrivalOfGames
	if _, err := g.GetAll(datastore.NewQuery("ArrivalOfGames").Order("CreatedAt"), &a); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Errorf(ctx, "GetAll error: %v", err)
		return
	}

	bodyBytes, err := json.Marshal(a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Errorf(ctx, "json marshal error: %v", err)
		return
	}

	log.Infof(ctx, "a: %v", a)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bodyBytes)
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

	} else if r.Method == "POST" {
		PostUser(w, r)
	}
}

// jsonを受け取ってdatastoreに保存
func PostUser(w http.ResponseWriter, r *http.Request)  {
	ctx := appengine.NewContext(r)
	g := goon.NewGoon(r)

	var u User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		log.Errorf(ctx, "json decode error: %v", err)
		return
	}

	log.Infof(ctx, "user: %v", u)

	if _, err := g.Put(&u); err != nil {
		log.Errorf(ctx, "user put error: %v", err)
		return
	}
}
