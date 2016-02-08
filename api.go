package board_gamers

import (
	"net/http"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"github.com/mjibson/goon"
)

func ArrivalOfGamesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	g := goon.NewGoon(r)

	var a []ArrivalOfGames
	if _, err := g.GetAll(datastore.NewQuery("ArrivalOfGames").Order("CreatedAt"), &a); err != nil {
		log.Errorf(ctx, "GetAll error: %v", err)
		return
	}
	log.Infof(ctx, "%v", a)
}
