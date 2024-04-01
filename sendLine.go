package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func SendLine(yesterday int, continueDays int, today int) {

	token := os.Getenv("TOKEN")
	if token == " " {
		log.Fatal("Errload env:TOKEN")
	}

	lineURL := "https://notify-api.line.me/api/notify"

	u, err := url.ParseRequestURI(lineURL)
	if err != nil {
		log.Fatal(err)
	}

	var text string

	var nowHour = ReturnHour()

	nowHour = 7

	if nowHour == 7 || nowHour == 8 {
		// 昨日までのコントリビューション数と連続コントリビューション日数を送信
		text = "昨日のコントリビューション数:" + strconv.Itoa(yesterday) + "\n" + "昨日までの連続コントリビューション日数:" + strconv.Itoa(continueDays)
	} else if (nowHour == 11 || nowHour == 12 || nowHour == 13) && today == 0 {
		text = "まだコミットしていませんね? 頑張りましょう!"
	} else if (nowHour >= 17) && today == 0 {
		text = "まだ時間はあります! 今すぐコーディングしましょう!"
	} else if today > 0 {
		text = "よく頑張りました! 明日も頑張りましょう!"
	}

	if text != "" {
		// メッセージを送信
		message := url.Values{"message": {text}}
		r, _ := http.NewRequest("POST", u.String(), strings.NewReader(message.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Set("Authorization", "Bearer "+token)
		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(body))
	}
}
