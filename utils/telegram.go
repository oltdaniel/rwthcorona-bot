package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

var TOKEN string = os.Getenv("TELEGRAM_TOKEN")

func telegramApiUrl(action string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%v/%v", TOKEN, action)
}

func TelegramSetWebhook(url string, certificateFilepath string) error {
	// prepare
	body := &bytes.Buffer{}
	// construct form fields
	form := multipart.NewWriter(body)
	// certificate field
	certificateFormFile, err := form.CreateFormFile("certificate", "certificate.crt")
	if err != nil {
		return err
	}
	fp, err := os.Open(certificateFilepath)
	if err != nil {
		return err
	}
	_, err = io.Copy(certificateFormFile, fp)
	if err != nil {
		return err
	}
	// url field
	urlFormField, err := form.CreateFormField("url")
	if err != nil {
		return err
	}
	urlFormField.Write([]byte(url))
	form.Close()
	// make request
	req, err := http.NewRequest("POST", telegramApiUrl("setWebhook"), bytes.NewReader(body.Bytes()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", form.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		return nil
	}
	return errors.New("telegram.setWebhook failed")
}

func TelegramSendMessage(message *TelegramRequestSendMessage) error {
	// serialize message
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	// make request
	req, err := http.NewRequest("POST", telegramApiUrl("sendMessage"), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		return nil
	}
	b, _ := ioutil.ReadAll(resp.Body)
	return errors.New(fmt.Sprintf("telegram.sendMessage failed: %v", string(b)))
}

func EscapeTelegramMessage(msg string) string {
	for _, special := range []string{"{", "}", "-", ".", "#", "=", "!", "|", ">"} {
		msg = strings.ReplaceAll(msg, special, "\\"+special)
	}
	return msg
}
