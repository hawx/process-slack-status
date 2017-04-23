package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-ps"
)

const statusMaxLength = 100

var lastStatus string

type slackClient struct {
	apiToken, apiURL, versionUID string
}

func (slack *slackClient) call(method string, args url.Values) error {
	args.Add("token", slack.apiToken)

	timestamp := time.Now().Unix()
	uri := slack.apiURL + method + "?_x_id=" + slack.versionUID + "-" + strconv.FormatInt(timestamp, 10)

	_, err := http.PostForm(uri, args)
	return err
}

func setStatus(slack *slackClient, emoji, text string) error {
	if emoji+text == lastStatus {
		return nil
	}
	lastStatus = emoji + text

	if len(text) > statusMaxLength {
		text = text[:statusMaxLength-2] + "…"
	}

	log.Printf("Setting status [%s] %s\n", emoji, text)

	payload, _ := json.Marshal(map[string]string{
		"status_text":  text,
		"status_emoji": emoji,
	})

	return slack.call("users.profile.set", url.Values{
		"profile": {string(payload)},
	})
}

func update(slack *slackClient, conf table) error {
	list, err := ps.Processes()
	if err != nil {
		return err
	}

	for _, process := range list {
		for k, v := range conf {
			if strings.ToLower(process.Executable()) == k {
				return setStatus(slack, v.Emoji, v.Text)
			}
		}
	}

	if defaults, ok := conf["*"]; ok {
		return setStatus(slack, defaults.Emoji, defaults.Text)
	}

	return nil
}

func start(slack *slackClient, conf table, tick time.Duration) error {
	if err := update(slack, conf); err != nil {
		return err
	}

	for _ = range time.Tick(tick) {
		if err := update(slack, conf); err != nil {
			return err
		}
	}

	return nil
}

type status struct {
	Emoji, Text string
}

type table map[string]status

func main() {
	var (
		apiToken   = flag.String("api-token", "", "Your Slack API token")
		apiURL     = flag.String("api-url", "", "Full URL to API path for the Slack team")
		versionUID = flag.String("version-uid", "", "The Slack version uid")
		tick       = flag.String("tick", "10s", "Duration to refresh status")
		config     = flag.String("config", "./config.toml", "")
	)
	flag.Parse()

	tickDuration, err := time.ParseDuration(*tick)
	if err != nil {
		log.Fatal(err)
	}

	var conf table
	if _, err = toml.DecodeFile(*config, &conf); err != nil {
		log.Fatal(err)
	}

	slack := &slackClient{
		apiToken:   *apiToken,
		apiURL:     *apiURL,
		versionUID: *versionUID,
	}

	if err := start(slack, conf, tickDuration); err != nil {
		log.Fatal(err)
	}
}
