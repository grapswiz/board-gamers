package board_gamers

import (
	"encoding/json"
	"github.com/mjibson/goon"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"net/http"
)

type Auth struct {
	IsLoggedIn	bool `json:"isLoggedIn"`
	User	*User `json:"user"`
}

func ArrivalOfGamesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	g := goon.NewGoon(r)

	var a []ArrivalOfGames
	if _, err := g.GetAll(datastore.NewQuery("ArrivalOfGames").Order("-CreatedAt"), &a); err != nil {
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

func AuthHandler(w http.ResponseWriter, r *http.Request)  {
	ctx := appengine.NewContext(r)

	var isLoggedIn bool
	var id string

	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		isLoggedIn = false;
		id = ""

		log.Infof(ctx, "sessionStore.get error: %v", err)
	} else {
		value := session.Values[sessionUserKey]
		log.Infof(ctx, "id: %v", value)

		valueStr, ok := value.(string)
		if ok {
			isLoggedIn = true;
			id = valueStr
		} else {
			isLoggedIn = false;
			id = ""
		}
	}

	var a *Auth
	if isLoggedIn {
		g := goon.NewGoon(r)
		u := &User{
			UserId: id,
		}
		g.Get(u)

		a = &Auth{
			IsLoggedIn: isLoggedIn,
			User: u,
		}
	} else {
		a = &Auth{
			IsLoggedIn: isLoggedIn,
			User: nil,
		}
	}

	bodyBytes, err := json.Marshal(a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Errorf(ctx, "json marshal error: %v", err)
		return
	}

	log.Infof(ctx, "auth: %v %v", a.IsLoggedIn, a.User)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bodyBytes)
}
