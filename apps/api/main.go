package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Traunin/stickerpack-editor/server/config"
	"gopkg.in/telebot.v3"
)

type StickerPackBot struct {
	*telebot.Bot
}

type Stickerpack struct {
	Name string `json:"name"`
}

var StickerpackMocks = []Stickerpack{
	{"monkaS"},
	{"peepoS"},
	{"Jokerge"},
	{"xdd"},
	{"hiii"},
	{"tuh"},
}

func (bot *StickerPackBot) createMockStickerSet() {
	bot.Handle("/start", func(c telebot.Context) error {
		fmt.Printf("Message from %d\n", c.Chat().ID)
		return c.Send(fmt.Sprintf("%d", c.Chat().ID))
	})
}

func LaunchBot() {
	bot, err := telebot.NewBot(telebot.Settings{
		Token: config.BotToken,
	})

	if err != nil {
		print(err)
	}

	// bot.createNewStickerSet
	bot.Handle("/start", func(c telebot.Context) error {
		fmt.Printf("Message from %d\n", c.Chat().ID)
		return c.Send(fmt.Sprintf("%d", c.Chat().ID))
	})

	bot.Start()
}

type TelegramRecipient struct {
	chatID string
}

func (r TelegramRecipient) Recipient() string {
	return r.chatID
}

func getUserById(id int64) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token: config.BotToken,
	})

	if err != nil {
		return
	}

	recipient := &TelegramRecipient{strconv.Itoa(int(id))}
	err = bot.CreateStickerSet(recipient, &telebot.StickerSet{
		Type:     telebot.StickerRegular,
		Format:   telebot.StickerStatic,
		Name:     "jokerge_by_seventv_stickerpack_bot",
		Title:    "Jokerge",
		Animated: true,
		Video:    true,
		Input: []telebot.InputSticker{
			{
				File: telebot.FromDisk("jokerge.webp"),
				Emojis: []string{
					"🤡",
				},
				Keywords: []string{
					"Jokerge",
				},
			},
		},
	})
	if err != nil {
		println("stickerpack error")
		println(err.Error())
	}

	fmt.Println(id)
	chat, err := bot.ChatByID(id)

	if err != nil {
		println("No user with such id")
	}

	fmt.Println(chat.StickerSet)
}

func handleStickerpacks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	json.NewEncoder(w).Encode(StickerpackMocks)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	getUserById(int64(id))
	// json.NewEncoder(w).Encode(StickerpackMocks)
}

func main() {
	// LaunchBot()
	http.HandleFunc("/api/stickerpacks", handleStickerpacks)
	http.HandleFunc("/api/user", handleUser)

	http.ListenAndServe(":8080", nil)
}
