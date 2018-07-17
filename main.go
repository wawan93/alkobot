package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/wawan93/bot-framework"
	"os"
)

func main() {
	token := os.Getenv("TOKEN")
	fmt.Printf("%v\n", token)
	api, _ := tgbotapi.NewBotAPI(token)

	u := tgbotapi.NewUpdate(0)
	updates, _ := api.GetUpdatesChan(u)

	bot := tgbot.NewBotFramework(api)

	bot.RegisterCommand("/start", func(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
		chatID := bot.GetChatID(update)
		msg := tgbotapi.NewMessage(chatID, "Hello, World!")
		_, err := bot.Send(msg)
		return err
	}, 0)

	bot.HandleUpdates(updates)
}
