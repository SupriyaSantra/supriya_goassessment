package main

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type TypeVal struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Resp struct {
	Event           string             `json:"event"`
	EventType       string             `json:"event_type"`
	AppID           string             `json:"app_id"`
	UserID          string             `json:"user_id"`
	MessageID       string             `json:"message_id"`
	PageTitle       string             `json:"page_title"`
	PageUrl         string             `json:"page_url"`
	BrowserLanguage string             `json:"browser_language"`
	ScreenSize      string             `json:"screen_size"`
	Attributes      map[string]TypeVal `json:"attributes"`
	Traits          map[string]TypeVal `json:"traits"`
}

func main() {
	app := fiber.New()

	app.Post("/", func(c *fiber.Ctx) error {
		data := map[string]string{}
		req := c.BodyRaw()
		json.Unmarshal(req, &data)
		ch := make(chan string)
		go worker(data, ch)
		out := <-ch
		return c.SendString(out)
	})

	app.Listen(":3000")
}

func worker(data map[string]string, ch chan<- string) {
	resp := new(Resp)
	resp.Attributes = make(map[string]TypeVal)
	resp.Traits = make(map[string]TypeVal)
	for k, v := range data {
		switch {
		case k == "ev":
			resp.Event = v
		case k == "et":
			resp.EventType = v
		case k == "id":
			resp.AppID = v
		case k == "uid":
			resp.UserID = v
		case k == "mid":
			resp.MessageID = v
		case k == "t":
			resp.PageTitle = v
		case k == "p":
			resp.PageUrl = v
		case k == "l":
			resp.BrowserLanguage = v
		case k == "sc":
			resp.ScreenSize = v
		case strings.HasPrefix(k, "atrk"):
			index := k[4:]
			var tva TypeVal
			tva.Type = data["atrt"+index]
			tva.Value = data["atrv"+index]
			resp.Attributes[v] = tva
		case strings.HasPrefix(k, "uatrk"):
			index := k[5:]
			var tvt TypeVal
			tvt.Type = data["uatrt"+index]
			tvt.Value = data["uatrv"+index]
			resp.Traits[v] = tvt
		}
	}
	out, _ := json.Marshal(resp)
	ch <- string(out)
}
