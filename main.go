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
	"strings"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8000", "tau websocket")

/*
todo:  have a config file (yaml?)
      - each line of config could be a twitch event with a go template formatted command to run
*/

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
		message := fmt.Sprintf(" %s gave us a follow", result.EventData.UserName)
		execute([]string{"/home/rex/bin/event-message.sh", message})
	} else if strings.Contains(result.EventType, "point-redemption") {
		title := result.EventData.Reward.Title
		prompt := result.EventData.Reward.Prompt
		// add user here
		message := fmt.Sprintf("points: redeemed %s : %s ", title, prompt)
		execute([]string{"/home/rex/bin/event-message.sh", message})
	} else if strings.Contains(result.EventType, "raid") {
		user := result.EventData.FromBroadcasterUserName
		raiders := result.EventData.Viewers
		message := fmt.Sprintf("%s raided with %d viewers", user, raiders)
		execute([]string{"/home/rex/bin/event-message.sh", message})
	} else if strings.Contains(result.EventType, "subscribe") {
		message := fmt.Sprintf(" %s dropped a sub for %d months!  : %s ",
			result.EventData.Data.Message.UserName,
			result.EventData.Data.Message.StreakMonths,
			result.EventData.Data.Message.SubMessage.Message)
		execute([]string{"/home/rex/bin/event-message.sh", message})
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
