package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tgbot "github.com/wawan93/bot-framework"
)

var rnd = 2
var r *rand.Rand
var randomRange, randomStart int

var lastPinTime time.Time

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {
	token := os.Getenv("TOKEN")
	chat, err := strconv.ParseInt(os.Getenv("CHAT"), 10, 64)
	if err != nil {
		log.Println(os.Getenv("CHAT"))
		log.Panic(err)
	}

	webhookAddress := os.Getenv("WEBHOOK_ADDRESS")
	if webhookAddress == "" {
		log.Panic("WEBHOOK_ADDRESS is empty")
	}

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

	bot := tgbot.NewBotFramework(api)
	updates := getUpdatesChannel(api, webhookAddress)

	if err := bot.RegisterCommand("/pin", PinMessage, chat); err != nil {
		log.Fatalf("can't register command: %+v", err)
	}
	if err := bot.RegisterPlainTextHandler(RandomPhrase(chat), chat); err != nil {
		log.Fatalf("can't register handler: %+v", err)
	}
	bot.Send(tgbotapi.NewMessage(chat, "Я жив. Я легитимный."))

	bot.HandleUpdates(updates)
}

func getUpdatesChannel(api *tgbotapi.BotAPI, webhookAddress string) tgbotapi.UpdatesChannel {
	var updates tgbotapi.UpdatesChannel
	if os.Getenv("APP_ENV") == "development" {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates, _ = api.GetUpdatesChan(u)
	} else if os.Getenv("APP_ENV") == "production" {
		_, err := api.SetWebhook(tgbotapi.NewWebhook(
			"https://" + webhookAddress + "/alkobot",
		))
		if err != nil {
			log.Fatal(err)
		}
		updates = api.ListenForWebhook("/alkobot")
		go http.ListenAndServe("0.0.0.0:80", nil)
	}
	return updates
}

func PinMessage(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
	if update.Message == nil {
		return errors.New("message is empty")
	}
	if update.Message.ReplyToMessage == nil {
		return errors.New("message is not reply")
	}

	if time.Now().Before(lastPinTime.Add(time.Minute * 30)) {
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
		_, _ = bot.Send(&tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           bot.GetChatID(update),
				ReplyToMessageID: update.Message.MessageID,
			},
			Text: "Не выбрали меня админом, вот хуй вам теперь, а не пины! 😡",
		})
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
		"Утютю",
		"Внимание всем, помните, логики нет. Её просто нет. Рефлексировать бесполезно.",
		"Только настоящие демократы могут к хуям переругаться на пустом месте из-за невесть чего",
		"Кажется, все уже выпили :)",
		"Овчаренко наш вождь",
		"Блядь объявляю собрание алкочата",
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
		"Клуб недовольных граждан объявляется закрытым",
		"Виновные выявлены, к ним уже выехали с паяльником",
		"только отвернёшься на минуту, сразу демократию разведут :)",
		"как это вот взять и никого не наебать",
		"коммунизму бой",
		"Левакам бой",
		"халяве — бой",
		"Футуризму бой",
		"либертарианству бой",
		"социализму бой",
		"Не звони мне",
		"Модерация - очень мягкая, исключаются только посты с нецензурщиной. Возможно высказывание любого мнения.",
		"Лаврентьева бы прошла",
		"ступайте в ЖПА",
		"я вижу направленную в мой адрес микроагрессию и отказываюсь продолжать диалог",
		"боевая формация совершила плевок в лицо шаткой конструкции наших взглядов. хипстеры, не унывая, запустили месседж-бокс, но демократические дедушки сменили токинг поинт на фигу.",
		"Месседж бокс, собравший все токинг поинты в шаткую конструкцию",
		"на проекты в которых участвуют сомнительные личности типа Неймарка я донатить не буду",
		"* говорит, что собрал",
		"Либерализмом мальчики в детстве занимаются",
		"/jump_v_okno",
		"При всем уважении всех ко всем, все всё понимают",
		"Напоминаю про правило 34",
	}

	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	return phrases[r.Intn(len(phrases))]
}
