package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
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
var randomRange, randomStart int

func init() {
	godotenv.Load()
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {
	token := os.Getenv("TOKEN")
	chat, _ := strconv.ParseInt(os.Getenv("CHAT"), 10, 64)

	log.Printf("token=%v", token)
	log.Printf("chat=%v", chat)

	flag.IntVar(&randomRange, "r", 90, "range")
	flag.IntVar(&randomStart, "s", 10, "start")
	flag.Parse()

	log.Println(randomRange)
	log.Println(randomStart)

	api, _ := tgbotapi.NewBotAPI(token)
	//api.Debug = true

	log.Printf("logged in as %v", api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := api.GetUpdatesChan(u)

	bot := tgbot.NewBotFramework(api)

	bot.RegisterUniversalHandler(RandomPhrase, 0)

	bot.HandleUpdates(updates)
}

func RandomPhrase(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
	chatID := bot.GetChatID(update)
	if chatID > 0 {
		hash := md5.New()
		hash.Write([]byte(strconv.Itoa(int(chatID))))
		personToken := hex.EncodeToString(hash.Sum(nil))
		log.Println(chatID)
		log.Println(personToken)
		msg := tgbotapi.NewMessage(chatID, personToken)
		_, err := bot.Send(msg)
		return err
	} else {
		rnd--
		log.Println(rnd)
		if rnd == 0 {
			rnd = r.Intn(randomRange) + randomStart
			msg := tgbotapi.NewMessage(chatID, GetRandomPhrase())
			_, err := bot.Send(msg)
			return err
		}
	}
	return nil
}

func GetRandomPhrase() string {
	phrases := []string{
		"ууу, сексизм!",
		"Вы говорите на языке ненависти!",
		"Когда Кац так делает, вы этого не замечаете.",
		"Гомофобная лексика!",
		"Не хочу быть с вами в одной партии.",
		"У вас здесь токсичная атмосфера...",
		"Таким в Яблоке не место.",
		"Мур мур мур",
		"Помню мы на собраниях первички в 2007 и 2008 бухали, а потом шли на чистые и продолжали бухать там",
		"Вступайте в гендерную фракцию.",
		"Вы слишком правые...",
		"Митрохин! Победа!",
		"Митрохин!",
	}

	return phrases[r.Intn(len(phrases))]
}
