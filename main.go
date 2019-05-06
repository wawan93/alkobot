package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"flag"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/wawan93/bot-framework"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var rnd = 2
var r *rand.Rand
var randomRange, randomStart int

var lastPinTime time.Time

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

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	api.Debug = true

	log.Printf("logged in as %v", api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := api.GetUpdatesChan(u)

	bot := tgbot.NewBotFramework(api)

	if err := bot.RegisterCommand("/pin", PinMessage, 0); err != nil {
		log.Fatalf("can't register command: %+v", err)
	}
	if err := bot.RegisterPlainTextHandler(RandomPhrase(chat), 0); err != nil {
		log.Fatalf("can't register handler: %+v", err)
	}
	bot.Send(tgbotapi.NewMessage(chat, "Я жив. Я легитимный."))

	bot.HandleUpdates(updates)
}

func PinMessage(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
	if update.Message == nil {
		return errors.New("message is empty")
	}
	if update.Message.ReplyToMessage == nil {
		return errors.New("message is not reply")
	}

	if time.Now().Before(lastPinTime.Add(time.Minute * 10)) {
		_, err := bot.Send(tgbotapi.NewMessage(bot.GetChatID(update), "Нельзя так часто пинить"))
		return err
	}

	msg := &tgbotapi.PinChatMessageConfig{
		ChatID:              bot.GetChatID(update),
		MessageID:           update.Message.ReplyToMessage.MessageID,
		DisableNotification: false,
	}

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("error pinning message: %+v", err)
	} else {
		lastPinTime = time.Now()
	}
	return err
}

func RandomPhrase(targetChat int64) tgbot.CommonHandler {
	return func(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
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
		} else if chatID == targetChat {
			rnd--
			log.Println(rnd)
			if rnd == 0 {
				rnd = r.Intn(randomRange) + randomStart
				msg := tgbotapi.NewMessage(chatID, GetRandomPhrase())
				_, err := bot.Send(msg)
				return err
			}
		} else if update.Message != nil && update.Message.Text != "" {
			if strings.Contains(update.Message.Text, "изыди") {
				bot.LeaveChat(tgbotapi.ChatConfig{ChatID: chatID})
			}
		}
		return nil
	}
}

func GetRandomPhrase() string {
	phrases := []string{
		"ууу, сексизм!",
		"Когда Кац так делает, вы этого не замечаете.",
		"У вас здесь токсичная атмосфера...",
		"Таким в Яблоке не место.",
		"Мур мур мур",
		"Вступайте в гендерную фракцию.",
		"Группа блядей",
		"Здесь блядство",
		"Давайте запретим Минаеву шутить",
		"Почему все разговоры в алкочате сводятся к ебле?",
		"Ныть не надо",
		"Давайте притеснять женщин",
		"А меня смущают Кацевские игры в демократию",
		"Нежное черебурчанье",
		"Утютю",
		"Внимание всем, помните, логики нет. Её просто нет. Рефлексировать бесполезно.",
		"Только настоящие демократы могут к хуям переругаться на пустом месте из-за невесть чего",
		"Кажется, все уже выпили :)",
		"Овчаренко наш вождь",
		"Блядь объявляю собрание алкочата",
		"Извинись!",
		"У нас тут блять Тартуга.",
		"Слышала про эту историю. Там не все так однозначно.",
		"А Кац уже ушёл?",
		"Митрохин же оспорит в суде",
		"Давайте завершим эту дискуссию.",
		"Давайте не будем о здравом смысле",
		"Аааа сделайте меня развидеть это!",
		"слыш ты чо гендеры попутал",
		"Как-то гомофобненько тут сегодня!",
		"Вы нихуя не понимаете в дизайне.",
		"как тебе такое влад неймарк",
		"это похоже на пиздёж",
		"во губу раскатали",
		"ваше ощущение ошибочно",
		"явная хрень",
		"поздравляю",
		"алё",
		"Клуб недовольных граждан объявляется закрытым",
		"Виновные выявлены, к ним уже выехали с паяльником",
		"только отвернёшься на минуту, сразу демократию разведут :)",
		"как это вот взять и никого не наебать",
		"Утипути",
		"коммунизму бой",
		"Левакам бой",
		"халяве — бой",
		"Футуризму бой",
		"либертарианству бой",
		"социализму бой",
	}

	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	return phrases[r.Intn(len(phrases))]
}
