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

func handleEvent(e []byte) {
	var result TauEvent
	fmt.Printf("%s\n", e)
	json.Unmarshal(e, &result)

	if strings.Contains(result.EventType, "follow") {
		message := fmt.Sprintf(" %s gave us a follow", result.EventData.UserName)
		cmd := exec.Command("/home/rex/bin/event-message.sh", message)
		log.Printf("Running command and waiting for it to finish...")
		err := cmd.Run()
		log.Printf("Command finished with error: %v", err)
	} else if strings.Contains(result.EventType, "point-redemption") {
		title := result.EventData.Reward.Title
		prompt := result.EventData.Reward.Prompt
	} else if strings.Contains(result.EventType, "subscribe") {
		message := fmt.Sprintf(" %s dropped a sub for %s months!  : %s ",
			result.EventData.Data.Message.UserName,
			result.EventData.Data.Message.StreakMonths,
			result.EventData.Data.Message.SubMessage.Message)
		cmd := exec.Command("/home/rex/bin/event-message.sh", message)
		log.Printf("Running command and waiting for it to finish...")
		err := cmd.Run()
		log.Printf("Command finished with error: %v", err)
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
