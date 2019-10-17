package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var WroteUsers *BotUsers

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_BOT_API_KEY"))
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	var titslink string
	var randomTitsNum int
	var successSend bool

	updates.Clear()

	WroteUsers = &BotUsers{
		userList: make(map[int]*BotUser),
	}

	http.HandleFunc("/", hello)
	go http.ListenAndServe(":8080", nil)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Command() == "start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "@"+update.Message.From.String()+" Ну какой start то? Я умею только /tits номер и /randomtits")
			bot.Send(msg)
			continue
		}

		if update.Message.Command() == "tits" {
			titsNum, _ := strconv.Atoi(update.Message.CommandArguments())
			titslink = "http://media.oboobs.ru/boobs_preview/" + strconv.Itoa(titsNum) + ".jpg"

			resp, err := http.Head(titslink)
			if err != nil {
				log.Print(err)
			} else {
				if resp.StatusCode == 200 {
					successSend = true
					SendTits(update, bot, titsNum, titslink)
				} else {
					log.Print("Send tits to " + update.Message.From.String() + " failed: " + titslink + " - return " + resp.Status + " :(")

					if resp.StatusCode == 404 {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "@"+update.Message.From.String()+" Нет такой картинки")
						bot.Send(msg)
					} else {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "@"+update.Message.From.String()+" Я сломался, извини :(")
						bot.Send(msg)
					}
				}
				resp.Body.Close()
			}
		}

		if update.Message.Command() == "randomtits" {
			successSend = false
			for i := 0; i < 3; i++ {
				randomTitsNum = Random(10000, 11549)
				titslink = "http://media.oboobs.ru/boobs_preview/" + strconv.Itoa(randomTitsNum) + ".jpg"

				resp, err := http.Head(titslink)
				if err != nil {
					log.Print(err)
				} else {
					if resp.StatusCode == 200 {
						successSend = true
						SendTits(update, bot, randomTitsNum, titslink)
						break
					} else {
						log.Print("Send tits to " + update.Message.From.String() + " failed: " + titslink + " - return " + resp.Status + ", retry :(")
					}
					resp.Body.Close()
				}
			}
			if successSend == false {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "@"+update.Message.From.String()+" Я сломался, извини :(")
				bot.Send(msg)
			}
		} else {
			log.Print("Unknown command: " + update.Message.Text + " from user @" + update.Message.From.String())
		}
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

func SendTits(update tgbotapi.Update, bot *tgbotapi.BotAPI, titsnum int, titslink string) {
	if update.Message.From.UserName == "soljarka" || (update.Message.From.UserName != "soljarka" && WroteUsers.CheckInterval(update.Message.From)) || update.Message.Chat.IsPrivate() {
		WroteUsers.AddUser(update.Message.From)

		file, err := DownloadFile(titslink)
		if err != nil {
			log.Print(err)
		}

		msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, file.Name())
		msg.Caption = "Сиськи №" + strconv.Itoa(titsnum) + " для @" + update.Message.From.String()
		//msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
		if update.Message.Chat.IsPrivate() {
			log.Print("Send tits to PRIVATE " + update.Message.From.String() + " - success")
		} else {
			log.Print("Send tits to CHANNEL " + update.Message.Chat.Title + " - success")
		}

		os.Remove(file.Name())
	} else {
		if !WroteUsers.CheckIgnore(update.Message.From) {
			WroteUsers.Ignore(update.Message.From)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "@"+update.Message.From.String()+" Вы уже получали сиськи, хватит спамить! Не более чем раз в 15 минут.")
			bot.Send(msg)
		}
		log.Print("Ignore request from user " + update.Message.From.String())
	}
}
