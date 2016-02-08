package board_gamers

import (
	"encoding/json"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"
	"regexp"
	"strings"
	"time"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
	"bytes"
	"github.com/mjibson/goon"
)

const (
	layout = "January 02, 2006 at 15:04PM"
)

type Tweet struct {
	UserName    string
	Text        string
	LinkToTweet string
	CreatedAt   string
}

type Values struct {
	Value1 string `json:"value1"`
	Value2 string `json:"value2"`
	Value3 string `json:"value3"`
}

type ArrivalOfGames struct {
	Id	int64	`datastore:"-" goon:"id"`
	Shop      string    `json:"shop"`
	Games     []string  `json:"games"`
	CreatedAt time.Time `json:"createdAt"`
	Url       string    `json:"url" datastore:",noindex"`
}

func init() {
	http.HandleFunc("/hello", handler)
	http.HandleFunc("/webhook/trickplay", trickplayHandler)


	http.HandleFunc("/api/v1/arrivalOfGames", ArrivalOfGamesHandler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}

func trickplayHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	g := goon.NewGoon(r)

	decoder := json.NewDecoder(r.Body)
	var t Tweet
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(w, "json parse error", 500)
		log.Errorf(ctx, "json parse error: %v", err)
	}
	//TODO 入荷した判定をする
	if !strings.Contains(t.Text, "入荷しております") {
		log.Infof(ctx, "no nyuuka")
		return
	}

	log.Infof(ctx, "this is 入荷 tweet: "+t.Text)

	//TODO 入荷した商品名を抽出する
	//TODO 全ての、の後ろにスペースを挿入する
	re := regexp.MustCompile("、?「(.+?)」|、?([^「」]+「.+?」)")
	submatch := re.FindAllStringSubmatch(t.Text, -1)
	var games []string
	for _, v := range submatch {
		if v[1] != "" {
			games = append(games, v[1])
		} else if v[2] != "" {
			games = append(games, v[2])
		}

	}
	log.Infof(ctx, "%v", games)

	createdAt, err := time.Parse(layout, t.CreatedAt)
	if err != nil {
		log.Errorf(ctx, "Time Parse error: %v", err)
		return
	}
	a := &ArrivalOfGames{
		Shop:      "トリックプレイ",
		Games:     games,
		CreatedAt: createdAt,
		Url:       t.LinkToTweet,
	}
	if _, err := g.Put(a); err != nil {
		log.Errorf(ctx, "Datastore put error: %v", err)
		return
	}

	postToIOS(ctx, a)
}

func postToIOS(ctx context.Context, a *ArrivalOfGames) {
	client := urlfetch.Client(ctx)

	param := Values{
		Value1:	a.Shop,
		Value2: a.Games[0],
		Value3: a.Games[1],
	}
	paramBytes, err := json.Marshal(param)
	if err != nil {
		log.Errorf(ctx, "json marshal error: %v", err)
		return
	}
	if err != nil {
		log.Errorf(ctx, "http request error: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = client.Do(req)
	if err != nil {
		log.Errorf(ctx, "client do error: %v", err)
		return
	}
}
