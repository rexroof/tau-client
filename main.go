package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8000", "tau websocket")

/*
todo:  have a config file (yaml?)
      - each line of config could be a twitch event with a go template formatted command to run
*/

// strips string of characters for printing
func sanitize(_in string) string {
	pat := regexp.MustCompile(`[^A-Za-z0-9-_+]+`)
	space := regexp.MustCompile(`\s+`)

	// replace all non-matching characters with spaces.
	_out := pat.ReplaceAllString(_in, " ")
	// any duplicate spaces become one
	_out = space.ReplaceAllString(_out, " ")

	return _out
}

func execute(_cmd []string) bool {
	app, args := _cmd[0], _cmd[1:]
	cmd := exec.Command(app, args...)
	log.Printf("execute: %s\n", strings.Join(args, " "))
	err := cmd.Run()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
		return false
	} else {
		return true
	}
}

func handleEvent(e []byte) {
	var result TauEvent
	fmt.Printf("%s\n", e)
	json.Unmarshal(e, &result)

	if strings.Contains(result.EventType, "follow") {
		message := fmt.Sprintf("%s followed!", result.EventData.UserName)
		execute([]string{"/home/rex/bin/follow-message.sh", message})

		// my different point rewards
	} else if strings.Contains(result.EventType, "channel-channel_points_custom_reward_redemption-add") {
		title := result.EventData.Reward.Title
		input := sanitize(result.EventData.UserInput)
		user := sanitize(result.EventData.UserName)

		log.Printf("channel points: %s %s\n", title, input)

		switch title {
    case "Eton Treats":
			execute([]string{"/home/rex/bin/eton-treats.sh", user})
		case "Eton Pets":
			execute([]string{"/home/rex/bin/eton-pets.sh", user})
		case "Test Reward, Do Not Use":
			execute([]string{"/home/rex/bin/test-reward.sh", user})
		case "Change Terminal Font":
			execute([]string{"/home/rex/bin/change-font.sh", input})
		default:
			message := fmt.Sprintf("%s %s: %s", user, title, input)
			execute([]string{"/home/rex/bin/event-message.sh", message})
		}

	} else if strings.Contains(result.EventType, "raid") {
		user := result.EventData.FromBroadcasterUserName
		raiders := result.EventData.Viewers
		message := fmt.Sprintf("%s raided with %d viewers", user, raiders)
		execute([]string{"/home/rex/bin/event-message.sh", message})
	} else if strings.Contains(result.EventType, "subscribe") {
		u := result.EventData.Data.Message.UserName
		mon := result.EventData.Data.Message.StreakMonths
		msg := sanitize(result.EventData.Data.Message.SubMessage.Message)
		message := fmt.Sprintf("%s subbed", u)
		if mon > 1 {
			message = fmt.Sprintf("%s x%d", message, mon)
		}
		if len(msg) > 0 {
			message = fmt.Sprintf("%s %s", message, msg)
		}
		execute([]string{"/home/rex/bin/sub-message.sh", message})
	} else {
		log.Println(result)
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws/twitch-events/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	tokenJSON := fmt.Sprintf("{\"token\": \"%s\"}", os.Getenv("TWITCH_WEBHOOK_SECRET"))

	err = c.WriteMessage(websocket.TextMessage, []byte(tokenJSON))
	if err != nil {
		log.Fatal("json token err:", err)
	}

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			handleEvent(message)
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			return
		}
	}
}
