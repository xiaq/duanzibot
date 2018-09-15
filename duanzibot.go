package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/xiaq/tg"
)

type DuanziBot struct {
	*tg.CommandBot
}

type Reply struct {
	Body string `json:"body"`
}

var errormsg = "发生了错误哦。快 /pia 我的主人吧。"

func get(text string) string {
	/*
		_, err := strconv.Atoi(text)
		if err != nil {
			return "要给我数字哦。"
		}
	*/
	url := "https://api.github.com/repos/tuna/collection/issues/" + url.QueryEscape(text)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("error http.Get(%q): %v\n", url, err)
		return errormsg
	}
	defer resp.Body.Close()

	var reply Reply
	err = json.NewDecoder(resp.Body).Decode(&reply)
	if err != nil {
		log.Printf("error decoding: %v\n", err)
		return "发生了错误哦。快 /pia 我的主人吧。"
	}
	if reply.Body == "" {
		return "GET 不到这个段子哦。"
	}
	return "https://github.com/tuna/collection/issues/" + text
	//return reply.Body
}

func (b *DuanziBot) handleDuanzi(_ *tg.CommandBot, text string, msg *tg.Message) {
	log.Println("/duanzi", text)
	err := b.Get("/sendMessage", tg.Query{
		"chat_id": msg.Chat.ID,
		"text":    get(text),
	}, nil)
	if err != nil {
		log.Println("error /sendMessage:", err)
	}
}

func NewDuanziBot(token string) *DuanziBot {
	b := &DuanziBot{tg.NewCommandBot(token)}
	b.OnCommand("duanzi", b.handleDuanzi)
	return b
}

func main() {
	buf, err := ioutil.ReadFile("token.txt")
	if err != nil {
		log.Fatalf("cannot read token file: %s", err)
	}
	token := strings.TrimSpace(string(buf))

	NewDuanziBot(token).Main()
}
