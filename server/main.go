package main

import (
	"fmt"

	telebot "gopkg.in/telebot.v3"
)
func main() {
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