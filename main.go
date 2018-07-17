package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/wawan93/bot-framework"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var rnd = 2
var r *rand.Rand

func init() {
	godotenv.Load()
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {
	token := os.Getenv("TOKEN")
	chat, _ := strconv.ParseInt(os.Getenv("CHAT"), 10, 64)

	log.Printf("token=%v", token)
	log.Printf("chat=%v", chat)

	api, _ := tgbotapi.NewBotAPI(token)
	//api.Debug = true

	log.Printf("logged in as %v", api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := api.GetUpdatesChan(u)

	bot := tgbot.NewBotFramework(api)

	bot.RegisterUniversalHandler(RandomPhrase, chat)

	bot.HandleUpdates(updates)
}

func RandomPhrase(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
	chatID := bot.GetChatID(update)
	log.Println(rnd)
	log.Println(chatID)
	rnd--
	if rnd == 0 {
		rnd = r.Intn(90) + 10
		msg := tgbotapi.NewMessage(chatID, GetRandomPhrase())
		_, err := bot.Send(msg)
		return err
	}
	return nil
}

func GetRandomPhrase() string {
	phrases := []string{
		"ууу, сексизм!",
		"Вы говорите на языке ненависти!",
		"Когда Кац так делает, вы этого не замечаете.",
		"Гомофонная лексика!",
		"Не хочу быть с вами в одной партии.",
	}

	return phrases[r.Intn(len(phrases))]
}
