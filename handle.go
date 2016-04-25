package board_gamers

import (
	"encoding/json"
	"github.com/mjibson/goon"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/dghubble/sessions"
	"github.com/garyburd/go-oauth/oauth"
	"google.golang.org/appengine/taskqueue"
	"io/ioutil"
)

const (
	layout         = "January 02, 2006 at 15:04PM"
	layout2        = "January 2, 2006 at 15:04PM"
	sessionName    = "board-gamers"
	sessionSecret  = "board-gamers-secret"
	sessionUserKey = "twitterID"
)

var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

var oauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authorize",
	TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
}

var secretKey string

type Tweet struct {
	UserName    string
	Text        string
	LinkToTweet string
	CreatedAt   string
	SecretKey   string
}

type Values struct {
	Value1 string `json:"value1"`
	Value2 string `json:"value2"`
}

type ArrivalOfGames struct {
	Id        int64     `datastore:"-" goon:"id" json:"id"`
	Shop      string    `json:"shop"`
	Games     []string  `json:"games"`
	CreatedAt time.Time `json:"createdAt"`
	Url       string    `json:"url" datastore:",noindex"`
}

type Config struct {
	TwitterConsumerKey    string
	TwitterConsumerSecret string
}

type User struct {
	UserId          string   `json:"userId" goon:"id"`
	ScreenName      string   `json:"screenName" datastore:",noindex"`
	Shops           []string `json:"shops"`
	NotificationKey string   `json:"notificationKey"`
}

type Shop struct {
	Name             string   `json:"name"`
	NotificationKeys []string `json:"notificationKeys"`
}

func init() {
	{
		b, err := ioutil.ReadFile("config.json")
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(b, &oauthClient.Credentials); err != nil {
			panic(err)
		}
	}

	{
		b, err := ioutil.ReadFile("secretKey")
		if err != nil {
			panic(err)
		}
		secretKey = string(b)
	}

	http.HandleFunc("/webhook/trickplay", TrickplayHandler)
	http.HandleFunc("/webhook/tendays", TendaysHandler)

	http.HandleFunc("/api/v1/arrivalOfGames", ArrivalOfGamesHandler)
	http.HandleFunc("/api/v1/user", UserHandler)
	http.HandleFunc("/api/v1/auth", AuthHandler)

	http.HandleFunc("/twitter/login", TwitterLoginHandler)
	http.HandleFunc("/twitter/callback", TwitterCallbackHandler)
	http.HandleFunc("/twitter/logout", TwitterLogoutHandler)
}

func TrickplayHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	decoder := json.NewDecoder(r.Body)
	var t Tweet
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(w, "json parse error", 500)
		log.Errorf(ctx, "json parse error: %v", err)
		return
	}

	if t.SecretKey != secretKey {
		log.Infof(ctx, "invalid secretKey received: %v", t.SecretKey)
		return
	}

	if !strings.Contains(t.Text, "#トリックプレイ") {
		log.Infof(ctx, "no nyuuka")
		return
	}

	games := extractTrickplayGames(t.Text)

	if len(games) == 0 {
		log.Infof(ctx, "this is no nyuuka")
		return
	}

	saveArrivalOfGames(ctx, w, "トリックプレイ", games, t.CreatedAt, t.LinkToTweet)

	postToIOS(ctx, w, "トリックプレイ", games, t.CreatedAt, t.LinkToTweet)
}

func extractTrickplayGames(text string) (games []string) {
	re := regexp.MustCompile("、?「(.+?)」|、?([^「」]+拡張「.+?」)|、?[^「」]+「(.+?)」")
	submatch := re.FindAllStringSubmatch(text, -1)
	for _, v := range submatch {
		if v[1] != "" {
			games = append(games, v[1])
		} else if v[2] != "" {
			games = append(games, v[2])
		} else if v[3] != "" {
			games = append(games, v[3])
		}

	}
}

func TendaysHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	decoder := json.NewDecoder(r.Body)
	var t Tweet
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(w, "json parse error", 500)
		log.Errorf(ctx, "json parse error: %v", err)
	}

	if t.SecretKey != secretKey {
		log.Infof(ctx, "invalid secretKey received: %v", t.SecretKey)
		return
	}

	if !strings.Contains(t.Text, "新入荷") && !strings.Contains(t.Text, "再入荷") {
		log.Infof(ctx, "no nyuuka")
		return
	}

	games := extractTendaysGames(t.Text)

	log.Infof(ctx, "%v", games)

	if len(games) == 0 {
		log.Infof(ctx, "this is no nyuuka")
		return
	}

	saveArrivalOfGames(ctx, w, "テンデイズ", games, t.CreatedAt, t.LinkToTweet)

	postToIOS(ctx, w, "テンデイズ", games, t.CreatedAt, t.LinkToTweet)
}

func extractTendaysGames(text string) (games []string) {
	re := regexp.MustCompile("、?「(.+?)」|、?([^「」]+拡張「.+?」)|\n(([^「].+[^、」])、?)を再入荷")
	submatch := re.FindAllStringSubmatch(text, -1)
	for _, v := range submatch {
		if v[1] != "" {
			games = append(games, v[1])
		} else if v[2] != "" {
			games = append(games, v[2])
		} else if v[3] != "" {
			games = append(games, strings.Split(v[3], "、")...)
		}

	}
}

func TwitterLoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	httpClient := urlfetch.Client(ctx)
	tmpCred, err := oauthClient.RequestTemporaryCredentials(httpClient, "http://"+r.Host+"/twitter/callback", nil)
	if err != nil {
		http.Error(w, "tmpCred error", http.StatusInternalServerError)
		log.Errorf(ctx, "tmpCred error: %v", err)
		return
	}

	http.Redirect(w, r, oauthClient.AuthorizationURL(tmpCred, nil), http.StatusFound)
	return
}

func TwitterCallbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	token := r.FormValue("oauth_token")
	tmpCred := &oauth.Credentials{
		Token:  token,
		Secret: oauthClient.Credentials.Secret,
	}
	httpClient := urlfetch.Client(ctx)
	tokenCred, v, err := oauthClient.RequestToken(httpClient, tmpCred, r.FormValue("oauth_verifier"))
	if err != nil {
		http.Error(w, "request token error", http.StatusInternalServerError)
		log.Errorf(ctx, "request token error: %v", err)
		return
	}
	log.Infof(ctx, "token cred: %v", tokenCred)
	log.Infof(ctx, "url.Values: %v", v)

	// セッションに保存
	session := sessionStore.New(sessionName)
	session.Values[sessionUserKey] = v["user_id"][0]
	session.Save(w)

	// ユーザIDを保存する
	u := &User{
		UserId:     v["user_id"][0],
		ScreenName: v["screen_name"][0],
	}
	log.Infof(ctx, "user: %v", u)
	g := goon.NewGoon(r)
	if _, err = g.Put(u); err != nil {
		log.Errorf(ctx, "goon put error: %v", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
	return
}

func TwitterLogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	sessionStore.Destroy(w, sessionName)
	log.Infof(ctx, "session destroyed")

	http.Redirect(w, r, "/", http.StatusFound)
	return
}

func isAuthenticated(req *http.Request) bool {
	if _, err := sessionStore.Get(req, sessionName); err == nil {
		return true
	}
	return false
}

var notificationPost = delay.Func("notificationPost", func(ctx context.Context, bodyStr string) {
	log.Infof(ctx, "delay httpPost")

	client := urlfetch.Client(ctx)

	if err != nil {
		log.Errorf(ctx, "httpPost request error: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = client.Do(req)
	if err != nil {
		log.Errorf(ctx, "httpPost client do error: %v", err)
		return
	}
})

var save = delay.Func("save", func(ctx context.Context, shop string, games []string, createdAt string, url string) {
	g := goon.FromContext(ctx)

	at, err := time.Parse(layout, createdAt)
	if err != nil {
		log.Errorf(ctx, "Time Parse error layout: %v", err)

		at, err = time.Parse(layout2, createdAt)
		if err != nil {
			log.Errorf(ctx, "Time Parse error layout2: %v", err)
			return
		}
	}
	a := &ArrivalOfGames{
		Shop:      shop,
		Games:     games,
		CreatedAt: at,
		Url:       url,
	}

	if _, err := g.Put(a); err != nil {
		log.Errorf(ctx, "Datastore put error: %v", err)
		return
	}
})

func saveArrivalOfGames(ctx context.Context, w http.ResponseWriter, shop string, games []string, createdAt string, url string) {
	t, err := save.Task(shop, games, createdAt, url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := taskqueue.Add(ctx, t, ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func postToIOS(ctx context.Context, w http.ResponseWriter, shop string, games []string, createdAt string, url string) {
	param := Values{
		Value1: shop,
		Value2: strings.Join(games, ","),
	}
	paramBytes, err := json.Marshal(param)
	if err != nil {
		log.Errorf(ctx, "json marshal error: %v", err)
		return
	}
	paramStr := string(paramBytes[:len(paramBytes)])

	t, err := notificationPost.Task(paramStr)
	if err != nil {
		log.Errorf(ctx, "notificationPost.Task error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := taskqueue.Add(ctx, t, ""); err != nil {
		log.Errorf(ctx, "taskqueue.Add error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
