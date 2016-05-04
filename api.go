package board_gamers

import (
	"encoding/json"
	"github.com/mjibson/goon"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"net/http"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
	"encoding/base64"
	"strings"
	"github.com/googlechrome/push-encryption-go/webpush"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/taskqueue"
)

type Auth struct {
	IsLoggedIn bool  `json:"isLoggedIn"`
	User       *User `json:"user"`
}

type Subscription struct {
	StatusType string `json:"statusType"`
	Endpoint string `json:"endpoint"`
	Keys Keys `json:"keys"`
	Shops []string `json:"shops"`
}

type Keys struct {
	P256dh string `json:"p256dh"`
	Auth string `json:"auth"`
}

type Notification struct {
	Title string `json:"title"`
	Body string `json:"body"`
	Tag string `json:"tag"`
	Icon string `json:"icon"`
}

type Shop struct {
	Name             string   `json:"name" goon:"id"`
	NotificationKeys []NotificationKey `json:"notificationKeys"`
}

type NotificationKey struct {
	Endpoint string `json:"endpoint"`
	Keys Keys `json:"keys"`
}

type Push struct {
	Shop string `json:"shop"`
	Notification Notification `json:"notification"`
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
func PostUser(w http.ResponseWriter, r *http.Request) {
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

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var isLoggedIn bool
	var id string

	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		isLoggedIn = false
		id = ""

		log.Infof(ctx, "sessionStore.get error: %v", err)
	} else {
		value := session.Values[sessionUserKey]
		log.Infof(ctx, "id: %v", value)

		valueStr, ok := value.(string)
		if ok {
			isLoggedIn = true
			id = valueStr
		} else {
			isLoggedIn = false
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
			User:       u,
		}
	} else {
		a = &Auth{
			IsLoggedIn: isLoggedIn,
			User:       nil,
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

func SubscriptionHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "GET" {
		return
	}

	ctx := appengine.NewContext(r)

	var sub Subscription

	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		log.Errorf(ctx, "json decode error: %v", err)
		return
	}

	if (sub.StatusType == "unsubscribe") {
		unsubscribe(ctx, sub.Shops, sub.Endpoint, sub.Keys)
		return
	}

	subscribe(ctx, sub.Shops, sub.Endpoint, sub.Keys)
}

func subscribe(ctx context.Context, shops []string, endpoint string, keys Keys)  {
	g := goon.FromContext(ctx)

	// load
	ss := []Shop{}
	for _, shop := range shops {
		ss = append(ss, Shop{
			Name: shop,
		})
	}
	log.Infof(ctx, "shops: %v", ss)
	if err := g.GetMulti(ss); err != nil {
		if me, ok := err.(appengine.MultiError); ok {
			for i, merr := range me {
				if merr == datastore.ErrNoSuchEntity {
					log.Infof(ctx, "key is not found. create it: %v", ss[i])
					g.Put(ss[i])
				}
			}
		} else {
			log.Errorf(ctx, "g.GetMulti(ss) error: %v", err)
			return
		}
	}

	for is, s := range ss {
		// check duplicated
		for _, v := range s.NotificationKeys {
			if (v.Endpoint == endpoint && v.Keys.P256dh == keys.P256dh && v.Keys.Auth == keys.Auth) {
				return
			}
		}

		// add
		ss[is].NotificationKeys = append(ss[is].NotificationKeys, NotificationKey{
			Endpoint: endpoint,
			Keys: keys,
		})
	}
	log.Infof(ctx, "ss: %v", ss)

	// save
	if _, err := g.PutMulti(ss); err != nil {
		log.Errorf(ctx, "g.PutMulti(ss) error: %v", err)
		return
	}
}

func unsubscribe(ctx context.Context, shops []string, endpoint string, keys Keys)  {
	g := goon.FromContext(ctx)

	// load
	ss := []Shop{}
	for _, shop := range shops {
		ss = append(ss, Shop{
			Name: shop,
		})
	}
	log.Infof(ctx, "shops: %v", ss)
	if err := g.GetMulti(ss); err != nil {
		if me, ok := err.(appengine.MultiError); ok {
			for i, merr := range me {
				if merr == datastore.ErrNoSuchEntity {
					log.Infof(ctx, "key is not found. create it: %v", ss[i])
					g.Put(ss[i])
				}
			}
		} else {
			log.Errorf(ctx, "g.GetMulti(ss) error: %v", err)
			return
		}
	}

	for is, s := range ss {
		var newNK []NotificationKey

		for _, v := range s.NotificationKeys {
			if (v.Endpoint != endpoint && v.Keys.P256dh != keys.P256dh && v.Keys.Auth != keys.Auth) {
				newNK = append(newNK, v)
			}
		}

		ss[is].NotificationKeys = newNK
	}
	log.Infof(ctx, "ss: %v", ss)

	// save
	if _, err := g.PutMulti(ss); err != nil {
		log.Errorf(ctx, "g.PutMulti(ss) error: %v", err)
		return
	}
}

//func PushHandler(w http.ResponseWriter, r *http.Request) {
//	if r.Method == "GET" {
//		return
//	}
//
//	ctx := appengine.NewContext(r)
//	g := goon.NewGoon(r)
//
//	var p Push
//
//	decoder := json.NewDecoder(r.Body)
//	if err := decoder.Decode(&p); err != nil {
//		log.Errorf(ctx, "json decode error: %v", err)
//		return
//	}
//
//	log.Infof(ctx, "push: %v", p)
//
//	s := &Shop{
//		Name: p.Shop,
//	}
//	if err := g.Get(s); err != nil {
//		log.Errorf(ctx, "g.Get(s) error: %v", err)
//		return
//	}
//	for _, key := range s.NotificationKeys {
//		bodyBytes, err := json.Marshal(p.Notification)
//		if err != nil {
//			log.Errorf(ctx, "json marshal error: %v", err)
//			return
//		}
//		push(ctx, key, bodyBytes)
//	}
//}
//
//func push(ctx context.Context, nk NotificationKey, bodyBytes []byte)  {
//	httpClient := urlfetch.Client(ctx)
//
//	b64 := base64.URLEncoding.WithPadding(base64.NoPadding)
//	key, err := b64.DecodeString(strings.TrimRight(nk.Keys.P256dh, "="))
//	if err != nil {
//		log.Errorf(ctx, "key decode error: %v", err)
//	}
//	auth, err := b64.DecodeString(strings.TrimRight(nk.Keys.Auth, "="))
//	if err != nil {
//		log.Errorf(ctx, "auth decode error: %v", err)
//		return
//	}
//
//	s := &webpush.Subscription{
//		Endpoint: nk.Endpoint,
//		Key: []byte(key),
//		Auth: []byte(auth),
//	}
//
//	_, err = webpush.Send(httpClient, s, string(bodyBytes), "AIzaSyAA5mxCYKrwjCZwB1E4Jqpp3UYLAEchJ6o")
//	if err != nil {
//		log.Errorf(ctx, "webpush.Send error: %v", err)
//		return
//	}
//}

var pushNotification = delay.Func("push", func(ctx context.Context, shop string, games []string) {
	log.Infof(ctx, "delay push")

	httpClient := urlfetch.Client(ctx)
	g := goon.FromContext(ctx)

	s := &Shop{
		Name: shop,
	}
	if err := g.Get(s); err != nil {
		log.Errorf(ctx, "g.Get(s) error: %v", err)
		return
	}
	n := &Notification{
		Title: "入荷速報",
		Body: shop + "さんに " + strings.Join(games, " ,") + " が入荷しました！",
		Tag: "push",
	}
	for _, nk := range s.NotificationKeys {
		bodyBytes, err := json.Marshal(n)
		if err != nil {
			log.Errorf(ctx, "json marshal error: %v", err)
			return
		}

		b64 := base64.URLEncoding.WithPadding(base64.NoPadding)
		key, err := b64.DecodeString(strings.TrimRight(nk.Keys.P256dh, "="))
		if err != nil {
			log.Errorf(ctx, "key decode error: %v", err)
		}
		auth, err := b64.DecodeString(strings.TrimRight(nk.Keys.Auth, "="))
		if err != nil {
			log.Errorf(ctx, "auth decode error: %v", err)
			return
		}

		s := &webpush.Subscription{
			Endpoint: nk.Endpoint,
			Key: []byte(key),
			Auth: []byte(auth),
		}

		_, err = webpush.Send(httpClient, s, string(bodyBytes), "AIzaSyAA5mxCYKrwjCZwB1E4Jqpp3UYLAEchJ6o")
		if err != nil {
			log.Errorf(ctx, "webpush.Send error: %v", err)
			return
		}
	}
})

func pushNotificationTask(ctx context.Context, w http.ResponseWriter, shop string, games []string)  {
	t, err := pushNotification.Task(shop, games)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := taskqueue.Add(ctx, t, ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}