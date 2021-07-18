package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/xiaq/tg"
)

func main() {
	buf, err := ioutil.ReadFile("token.txt")
	if err != nil {
		log.Fatalf("cannot read token file: %s", err)
	}
	token := strings.TrimSpace(string(buf))

	newBot(token).Main()
}

type bot struct {
	*tg.CommandBot
}

func newBot(token string) *bot {
	b := &bot{tg.NewCommandBot(token)}
	b.OnCommand("duanzi", b.handleDuanzi)
	return b
}

func (b *bot) handleDuanzi(_ *tg.CommandBot, text string, msg *tg.Message) {
	log.Println("/duanzi", text)
	err := b.Get("/sendMessage", tg.Query{
		"chat_id": msg.Chat.ID,
		"text":    get(text),
	}, nil)
	if err != nil {
		log.Println("error /sendMessage:", err)
	}
}

func get(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		// Find a random duanzi
		n := getMaxIssueNumber()
		i := rand.Intn(n) + 1
		return "https://github.com/tuna/collection/issues/" + strconv.Itoa(i)
	}
	if i, err := strconv.Atoi(text); err == nil {
		return "https://github.com/tuna/collection/issues/" + strconv.Itoa(i)
	}
	return "用法：/duanzi [id] （省略 id 则随机选择段子）"
}

type issuesReply []struct{ Number int }

var maxIssueNumber struct {
	sync.Mutex
	value  int
	update time.Time
}

const cacheTTL = 30 * time.Minute

func getMaxIssueNumber() int {
	maxIssueNumber.Lock()
	defer maxIssueNumber.Unlock()

	if time.Now().Sub(maxIssueNumber.update) > cacheTTL {
		url := "https://api.github.com/repos/tuna/collection/issues?per_page=1"
		resp, err := http.Get(url)
		if err != nil {
			log.Println("list issues:", err)
			return 1
		}
		var reply issuesReply
		err = json.NewDecoder(resp.Body).Decode(&reply)
		if err != nil {
			log.Println("decode issue list:", err)
			return 1
		}
		if len(reply) != 1 {
			log.Println("length of issue list is", len(reply))
			return 1
		}

		maxIssueNumber.value = reply[0].Number
		maxIssueNumber.update = time.Now()
	}
	return maxIssueNumber.value
}
