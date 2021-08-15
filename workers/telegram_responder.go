package workers

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/oltdaniel/rwthcorona-bot/utils"
)

var TelegramResponderQueue chan (utils.TelegramWebhookUpdate) = make(chan utils.TelegramWebhookUpdate, 10)

var TELEGRAM_USERNAME = os.Getenv("TELEGRAM_USERNAME")

func TelegramResponder() {
	stKreisDetails, err := utils.DATABASE.Prepare("SELECT tag, label, anzahlWoche, rateWoche, anteilWoche FROM corona_data WHERE plz=? AND altersgruppe=? AND tag=(SELECT tag FROM corona_data ORDER BY tag DESC LIMIT 1) LIMIT 1")
	defer stKreisDetails.Close()
	if err != nil {
		log.Fatal(err)
	}
	stAllAltersgruppen, err := utils.DATABASE.Prepare("SELECT altersgruppe FROM corona_data GROUP BY altersgruppe ORDER BY altersgruppe DESC")
	defer stAllAltersgruppen.Close()
	if err != nil {
		log.Fatal(err)
	}
	stAltersgruppe, err := utils.DATABASE.Prepare("SELECT tag, label, anzahlWoche, rateWoche, anteilWoche FROM corona_data WHERE altersgruppe=? AND tag=(SELECT tag FROM corona_data ORDER BY tag DESC LIMIT 1)")
	defer stAltersgruppe.Close()
	if err != nil {
		log.Fatal(err)
	}
	for {
		update := <-TelegramResponderQueue
		msg := update.Message.Text
		if update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup" {
			if strings.HasSuffix(msg, "@RWTHcorona_bot") {
				msg = msg[:len(msg)-len("@RWTHcorona_bot")]
			} else if update.Message.ReplyToMessage.From.Username != TELEGRAM_USERNAME {
				continue
			}
		}
		switch {
		case msg == "/info":
			err := utils.TelegramSendMessage(&utils.TelegramRequestSendMessage{
				ChatId:           update.Message.Chat.Id,
				Text:             utils.EscapeTelegramMessage("🦠 @RWTHcorona\\_bot \nIch kann nur Nachrichten lesen die mit `/` starten, Antworten auf meine Nachrichten \\(auch geschachtelte Antworten\\) und Benachrichtigungen zu z.B. angepinnten Nachrichten.\n\nJede Nachricht die mich somit erreicht, werde ich temporär bearbeiten. Wenn ich fertig bin, verschwindet alles.\n\nSpeziell in Gruppen verlange ich, dass die Kommandos auf meinen Benutzernamen enden. Dies passiert automatisch wenn mehrere Bots in der Gruppe sind. Zusätzlich ermögliche ich die Benutzung vomMarkup Keyboard, worüber die Altersgruppen ausgewählt werden.\n\n📫 *Datenschutz*\nIch speichere nichts über dich. Die gesendeten Chatverläufe sind aber natürlich bei Telegram entsprechend gespeichert.\n\n📂 *Datenquelle*\nDie Daten werden von dem [Land NRW](https://www.lzg.nrw.de/covid19/covid19_mags.html) bezogen und stündlich von uns auf Updates überprüft. Da es nu rein Update pro tag gibt, wird nur das Datum der Daten angegeben."),
				ReplyToMessageId: update.Message.MessageId,
				ParseMode:        "MarkdownV2",
			})
			if err != nil {
				log.Println(err)
			}
		case msg == "/aktuell":
			resp := ""
			getKreisDetails := func(kreis int64) (string, error) {
				rresp := ""
				r, err := stKreisDetails.Query(kreis, "gesamt")
				if err != nil {
					return "", err
				}
				for r.Next() {
					var (
						tag         string
						label       string
						anzahlWoche int64
						rateWoche   float64
						anteilWoche float64
					)
					err = r.Scan(&tag, &label, &anzahlWoche, &rateWoche, &anteilWoche)
					if err != nil {
						return "", err
					}
					lastUpdate := "unbekannt"
					date, err := time.Parse("2006-01-02T15:04:05Z", tag)
					if err == nil {
						lastUpdate = date.Format("2006-01-02")
					}
					rresp += fmt.Sprintf("\n\n📌 *%v* _\\(`%v`\\)_\nAnzahl: `%d`\nIndzidenz: `%0.2f`", label, lastUpdate, anzahlWoche, rateWoche)
				}
				r.Close()
				rresp = utils.EscapeTelegramMessage(rresp)
				return rresp, nil
			}
			for _, kreis := range []int64{5, 5334} {
				data, err := getKreisDetails(kreis)
				if err != nil {
					utils.TelegramSendMessage(&utils.TelegramRequestSendMessage{
						ChatId:           update.Message.Chat.Id,
						Text:             err.Error(),
						ReplyToMessageId: update.Message.MessageId,
					})
					continue
				}
				resp += data
			}
			err = utils.TelegramSendMessage(&utils.TelegramRequestSendMessage{
				ChatId:           update.Message.Chat.Id,
				Text:             resp,
				ReplyToMessageId: update.Message.MessageId,
				ParseMode:        "MarkdownV2",
			})
			if err != nil {
				log.Println(err)
			}
		case strings.HasPrefix(msg, "/altersgruppe"):
			commands := strings.SplitN(msg, " ", 2)
			if len(commands) == 1 {
				getAllAltersgruppen := func() ([]utils.TelegramReplyMarkupReplyKeyboardButton, error) {
					rresp := []utils.TelegramReplyMarkupReplyKeyboardButton{}
					r, err := stAllAltersgruppen.Query()
					if err != nil {
						return rresp, err
					}
					for r.Next() {
						var (
							altersgruppe string
						)
						err = r.Scan(&altersgruppe)
						if err != nil {
							return rresp, err
						}
						rresp = append(rresp, utils.TelegramReplyMarkupReplyKeyboardButton{Text: altersgruppe})
					}
					r.Close()
					return rresp, nil
				}
				resp, err := getAllAltersgruppen()
				if err != nil {
					utils.TelegramSendMessage(&utils.TelegramRequestSendMessage{
						ChatId:           update.Message.Chat.Id,
						Text:             err.Error(),
						ReplyToMessageId: update.Message.MessageId,
					})
					continue
				}
				buttons := [][]utils.TelegramReplyMarkupReplyKeyboardButton{}
				index := -1
				for i, v := range resp {
					if i%4 == 0 {
						index += 1
						buttons = append(buttons, []utils.TelegramReplyMarkupReplyKeyboardButton{})
					}
					buttons[index] = append(buttons[index], v)
				}
				err = utils.TelegramSendMessage(&utils.TelegramRequestSendMessage{
					ChatId:           update.Message.Chat.Id,
					ReplyToMessageId: update.Message.MessageId,
					Text:             "Wähle Altersgruppe:",
					ReplyMarkup: utils.TelegramRequestSendMessageReplyMarkup{
						Keyboard:        buttons,
						OneTimeKeyboard: true,
						Selective:       true,
					},
				})
				if err != nil {
					log.Println(err)
				}
				continue
			}
			if len(commands) != 2 {
				err = utils.TelegramSendMessage(&utils.TelegramRequestSendMessage{
					ChatId:           update.Message.Chat.Id,
					Text:             "this is not allowed",
					ReplyToMessageId: update.Message.MessageId,
				})
				if err != nil {
					log.Println(err)
				}
				continue
			}
			altersgruppe := commands[1]
			altersgruppe = strings.ReplaceAll(altersgruppe, "-", " - ")
			altersgruppe = strings.ReplaceAll(altersgruppe, "+", " u. älter")
			altersgruppe = strings.ReplaceAll(altersgruppe, "  ", " ")
			getAltersgruppe := func(altersgruppe string) (string, error) {
				rresp := ""
				r, err := stAltersgruppe.Query(altersgruppe)
				if err != nil {
					return "", err
				}
				for r.Next() {
					var (
						tag         string
						label       string
						anzahlWoche int64
						rateWoche   float64
						anteilWoche float64
					)
					err = r.Scan(&tag, &label, &anzahlWoche, &rateWoche, &anteilWoche)
					if err != nil {
						return "", err
					}
					lastUpdate := "unbekannt"
					date, err := time.Parse("2006-01-02T15:04:05Z", tag)
					if err == nil {
						lastUpdate = date.Format("2006-01-02")
					}
					rresp += fmt.Sprintf("\n\n📌 *%v* _\\(`%v`\\)_\nAltersgruppe: `%v`\nAnzahl: `%d`\nIndzidenz: `%0.2f`", label, lastUpdate, altersgruppe, anzahlWoche, rateWoche)
				}
				r.Close()
				rresp = utils.EscapeTelegramMessage(rresp)
				return rresp, nil
			}
			resp, err := getAltersgruppe(altersgruppe)
			if err != nil {
				err = utils.TelegramSendMessage(&utils.TelegramRequestSendMessage{
					ChatId:           update.Message.Chat.Id,
					Text:             err.Error(),
					ReplyToMessageId: update.Message.MessageId,
				})
				if err != nil {
					log.Println(err)
				}
				continue
			}
			err = utils.TelegramSendMessage(&utils.TelegramRequestSendMessage{
				ChatId:           update.Message.Chat.Id,
				Text:             resp,
				ReplyToMessageId: update.Message.MessageId,
				ParseMode:        "MarkdownV2",
			})
			if err != nil {
				log.Println(err)
			}
		case update.Message.ReplyToMessage.From.Username == TELEGRAM_USERNAME || update.Message.Chat.Type == "private":
			altersgruppe := msg
			altersgruppe = strings.ReplaceAll(altersgruppe, "-", " - ")
			altersgruppe = strings.ReplaceAll(altersgruppe, "+", " u. älter")
			altersgruppe = strings.ReplaceAll(altersgruppe, "  ", " ")
			getAltersgruppe := func(altersgruppe string) (string, error) {
				rresp := ""
				r, err := stAltersgruppe.Query(altersgruppe)
				if err != nil {
					return "", err
				}
				for r.Next() {
					var (
						tag         string
						label       string
						anzahlWoche int64
						rateWoche   float64
						anteilWoche float64
					)
					err = r.Scan(&tag, &label, &anzahlWoche, &rateWoche, &anteilWoche)
					if err != nil {
						return "", err
					}
					lastUpdate := "unbekannt"
					date, err := time.Parse("2006-01-02T15:04:05Z", tag)
					if err == nil {
						lastUpdate = date.Format("2006-01-02")
					}
					rresp += fmt.Sprintf("\n\n📌 *%v* _\\(`%v`\\)_\nAltersgruppe: `%v`\nAnzahl: `%d`\nIndzidenz: `%0.2f`", label, lastUpdate, altersgruppe, anzahlWoche, rateWoche)
				}
				r.Close()
				rresp = utils.EscapeTelegramMessage(rresp)
				return rresp, nil
			}
			resp, err := getAltersgruppe(altersgruppe)
			if err != nil {
				err = utils.TelegramSendMessage(&utils.TelegramRequestSendMessage{
					ChatId:           update.Message.Chat.Id,
					Text:             err.Error(),
					ReplyToMessageId: update.Message.MessageId,
				})
				if err != nil {
					log.Println(err)
				}
				continue
			}
			err = utils.TelegramSendMessage(&utils.TelegramRequestSendMessage{
				ChatId:           update.Message.Chat.Id,
				Text:             resp,
				ReplyToMessageId: update.Message.MessageId,
				ParseMode:        "MarkdownV2",
			})
			if err != nil {
				log.Println(err)
			}
			if update.Message.ReplyToMessage.From.Username == TELEGRAM_USERNAME {
				err = utils.TelegramDeleteMessage(&utils.TelegramRequestDeleteMessage{
					ChatId:    update.Message.ReplyToMessage.Chat.Id,
					MessageId: update.Message.ReplyToMessage.MessageId,
				})
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}
