package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/telebot.v3"
)

type Stickerpack struct {
    Name string `json:"name"`
}

func LaunchBot() {
    bot, err := telebot.NewBot(telebot.Settings {
        Token: BotToken,
    })

    if (err != nil) {
        print(err)
    }

    // bot.createNewStickerSet
    bot.Handle("/start", func(c telebot.Context) error {
        fmt.Printf("Message from %d\n", c.Chat().ID)
		return c.Send(fmt.Sprintf("%d", c.Chat().ID))
	})

	bot.Start()
}

func main() {
    http.HandleFunc("/api/stickerpacks", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*") 
        w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")  
        
        monkas := Stickerpack{
            Name: "MONKAS",
        }

        json.NewEncoder(w).Encode(monkas)
    })

	http.ListenAndServe(":8080", nil)
}