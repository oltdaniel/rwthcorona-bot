package workers

import (
	"fmt"
	"strings"

	"github.com/oltdaniel/rwthcorona-bot/utils"
)

var TelegramResponderQueue chan (*utils.TelegramWebhookUpdate) = make(chan *utils.TelegramWebhookUpdate)

func TelegramResponder() {
	for {
		update := <-TelegramResponderQueue
		msg := update.Message.Text
		switch {
		case msg == "/aktuell":
			resp := ""
			getKreisDetails := func(kreis int64) (string, error) {
				rresp := ""
				st, err := utils.DATABASE.Prepare("SELECT label, anzahlWoche, rateWoche, anteilWoche FROM corona_data WHERE plz=? AND altersgruppe=? AND tag=date('now', '-1 day') LIMIT 1")
				if err != nil {
					return "", err
				}
				r, err := st.Query(kreis, "gesamt")
				if err != nil {
					return "", err
				}
				for r.Next() {
					var (
						label       string
						anzahlWoche int64
						rateWoche   float64
						anteilWoche float64
					)
					err = r.Scan(&label, &anzahlWoche, &rateWoche, &anteilWoche)
					if err != nil {
						return "", err
					}
					rresp += fmt.Sprintf("\n\n*%v*:\nAnzahl: `%d`\nIndzidenz: `%0.2f`", label, anzahlWoche, rateWoche)
				}
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
					return
				}
				resp += data
			}
			utils.TelegramSendMessage(&utils.TelegramRequestSendMessage{
				ChatId:           update.Message.Chat.Id,
				Text:             resp,
				ReplyToMessageId: update.Message.MessageId,
				ParseMode:        "MarkdownV2",
			})
		case strings.HasPrefix(msg, "/altersgruppe"):
			commands := strings.SplitN(msg, " ", 2)
			if len(commands) == 1 {
				getAllAltersgruppen := func() (string, error) {
					rresp := ""
					st, err := utils.DATABASE.Prepare("SELECT altersgruppe FROM corona_data WHERE tag=date('now', '-1 day') GROUP BY altersgruppe ORDER BY altersgruppe DESC")
					if err != nil {
						return "", err
					}
					r, err := st.Query()
					if err != nil {
						return "", err
					}
					for r.Next() {
						var (
							altersgruppe string
						)
						err = r.Scan(&altersgruppe)
						if err != nil {
							return "", err
						}
						rresp += fmt.Sprintf("_%v_, ", altersgruppe)
					}
					rresp = rresp[:len(rresp)-2]
					rresp = utils.EscapeTelegramMessage(rresp)
					return rresp, nil
				}
				resp, err := getAllAltersgruppen()
				if err != nil {
					utils.TelegramSendMessage(&utils.TelegramRequestSendMessage{
						ChatId:           update.Message.Chat.Id,
						Text:             err.Error(),
						ReplyToMessageId: update.Message.MessageId,
					})
					return
				}
				resp = "*Altersgruppen*:\n" + resp
				utils.TelegramSendMessage(&utils.TelegramRequestSendMessage{
					ChatId:           update.Message.Chat.Id,
					Text:             resp,
					ReplyToMessageId: update.Message.MessageId,
					ParseMode:        "MarkdownV2",
				})
				return
			}
			if len(commands) != 2 {
				utils.TelegramSendMessage(&utils.TelegramRequestSendMessage{
					ChatId:           update.Message.Chat.Id,
					Text:             "this is not allowed",
					ReplyToMessageId: update.Message.MessageId,
				})
				return
			}
			altersgruppe := commands[1]
			altersgruppe = strings.ReplaceAll(altersgruppe, " ", "")
			altersgruppe = strings.ReplaceAll(altersgruppe, "-", " - ")
			altersgruppe = strings.ReplaceAll(altersgruppe, "+", " u. älter")
			getAltersgruppe := func(altersgruppe string) (string, error) {
				rresp := ""
				st, err := utils.DATABASE.Prepare("SELECT label, anzahlWoche, rateWoche, anteilWoche FROM corona_data WHERE altersgruppe=? AND tag=date('now', '-1 day')")
				if err != nil {
					return "", err
				}
				r, err := st.Query(altersgruppe)
				if err != nil {
					return "", err
				}
				for r.Next() {
					var (
						label       string
						anzahlWoche int64
						rateWoche   float64
						anteilWoche float64
					)
					err = r.Scan(&label, &anzahlWoche, &rateWoche, &anteilWoche)
					if err != nil {
						return "", err
					}
					rresp += fmt.Sprintf("\n\n*%v* _(%v)_:\nAnzahl: `%d`\nIndzidenz: `%0.2f`", label, altersgruppe, anzahlWoche, rateWoche)
				}
				rresp = utils.EscapeTelegramMessage(rresp)
				return rresp, nil
			}
			resp, err := getAltersgruppe(altersgruppe)
			if err != nil {
				utils.TelegramSendMessage(&utils.TelegramRequestSendMessage{
					ChatId:           update.Message.Chat.Id,
					Text:             err.Error(),
					ReplyToMessageId: update.Message.MessageId,
				})
				return
			}
			utils.TelegramSendMessage(&utils.TelegramRequestSendMessage{
				ChatId:           update.Message.Chat.Id,
				Text:             resp,
				ReplyToMessageId: update.Message.MessageId,
				ParseMode:        "MarkdownV2",
			})
			return
		default:
			utils.TelegramSendMessage(&utils.TelegramRequestSendMessage{
				ChatId:           update.Message.Chat.Id,
				Text:             "not implemented yet",
				ReplyToMessageId: update.Message.MessageId,
			})
		}
	}
}
