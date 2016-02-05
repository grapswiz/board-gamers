package board_gamers

import (
	"net/http"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"encoding/json"
	"strings"
	"regexp"
	"bytes"
)

type Tweet struct {
	UserName	string
	Text		string
	LinkToTweet	string
	CreatedAt	string
}

type Values struct {
	Value1		string        `json:"value1"`
	Value2		string        `json:"value2"`
	Value3		string        `json:"value3"`
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
	//TODO 入荷した判定をする
	if !strings.Contains(t.Text, "入荷しております") {
		log.Infof(ctx, "no nyuuka")
		return
	}

	log.Infof(ctx, "this is 入荷 tweet: " + t.Text)

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

	client := urlfetch.Client(ctx)

	param := Values{
		Value1: "あるボドゲ屋",
		Value2: games[1],
		Value3: games[2],
	}
	paramBytes, err := json.Marshal(param)
	if err != nil {
		log.Errorf(ctx, "%v", err)
		return
	}
	if err != nil {
		log.Errorf(ctx, "%v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = client.Do(req)
	if err != nil {
		log.Errorf(ctx, "%v", err)
		return
	}
}
