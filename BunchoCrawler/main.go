package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	IncomingURL = "https://hooks.slack.com/services/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"
)

type Slack struct {
	Text        string       `json:"text"`
	Username    string       `json:"username"`
	IconEmoji   string       `json:"icon_emoji"`
	IconURL     string       `json:"icon_url"`
	Channel     string       `json:"channel"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Fallback   string   `json:"fallback"`
	Color      string   `json:"color"`
	Pretext    string   `json:"pretext"`
	AuthorName string   `json:"author_name"`
	AuthorLink string   `json:"author_link"`
	AuthorIcon string   `json:"author_icon"`
	Title      string   `json:"title"`
	TitleLink  string   `json:"title_link"`
	Text       string   `json:"text"`
	Fields     [3]Field `json:"fields"`
	ImageURL   string   `json:"image_url"`
	ThumbURL   string   `json:"thumb_url"`
	Footer     string   `json:"footer"`
	FooterIcon string   `json:"footer_icon"`
	Ts         int      `json:"ts"`
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type Attachments []Attachment

func Filter(t time.Time) (ret Attachments) {
	return nil
}

type Config struct {
	Site    string   `json:"site"`
	History []string `json:"history"`
}

func Map(s *goquery.Selection) (ret Attachment) {
	title := s.Find("td.sca_name2 a").Text()
	attr, _ := s.Find("td.sca_name2 a").Attr("href")
	price := s.Find("td.price_l").Text()
	number := s.Find("div.info_area2 > table > tbody > tr:nth-child(2) > td:nth-child(2)").Text()
	shop := s.Find("span.shop_nm2_nm").Text()
	comment := s.Find("td.td_p").Text()
	imageURL, _ := s.Find("p.photo img").Attr("src")

	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(255)
	g := rand.Intn(255)
	b := rand.Intn(255)
	color := fmt.Sprintf("#%X%X%X", r, g, b)

	var fields [3]Field
	fields[0] = Field{
		Title: "価格",
		Value: price,
		Short: false,
	}
	fields[1] = Field{
		Title: "匹数",
		Value: number,
		Short: false,
	}
	fields[2] = Field{
		Title: "店舗",
		Value: shop,
		Short: false,
	}
	return Attachment{
		Color:      color,
		Title:      title,
		TitleLink:  attr,
		Fields:     fields,
		Text:       comment,
		ImageURL:   imageURL,
		Footer:     "ペットショップのコジマ",
		FooterIcon: "https://www.google.com/s2/favicons?domain=pets-kojima.com",
	}
}

func main() {
	target := "http://pets-kojima.com/small_list/?topics_group_id=4&group=&shop%5B%5D=56529&shop%5B%5D=15&shop%5B%5D=54&shop%5B%5D=148&shop%5B%5D=149&shop%5B%5D=150&shop%5B%5D=151&shop%5B%5D=152&shop%5B%5D=153&shop%5B%5D=154&shop%5B%5D=155&shop%5B%5D=156&shop%5B%5D=145&shop%5B%5D=157&shop%5B%5D=158&shop%5B%5D=91960&shop%5B%5D=159&shop%5B%5D=160&shop%5B%5D=161&shop%5B%5D=187095&shop%5B%5D=170&price_bottom=&price_upper=&freeword=%E3%83%96%E3%83%B3%E3%83%81%E3%83%A7%E3%82%A6&order_type=2&x=99&y=38"

	var config Config
	file, err := ioutil.ReadFile("config.json")
	if err == nil {
		json.Unmarshal(file, &config)
	}

	doc, err := goquery.NewDocument(target)
	if err != nil {
		fmt.Println(err)
		return
	}

	var attachments []Attachment
	var histories []string
	doc.Find("div.sca_table2").Each(func(_ int, s *goquery.Selection) {
		attachment := Map(s)
		contains := false
		for _, history := range config.History {
			if attachment.TitleLink == history {
				contains = true
				break
			}
		}
		if !contains {
			attachments = append(attachments, attachment)
		}
		histories = append(histories, attachment.TitleLink)
	})

	writeConfig, _ := json.Marshal(Config{
		Site:    "kojima",
		History: histories,
	})
	ioutil.WriteFile("config.json", writeConfig, os.ModePerm)

	if len(attachments) == 0 {
		return
	}

	params, _ := json.Marshal(Slack{
		Text:        fmt.Sprintf("本日のブンチョウたち（%s）", time.Now().Format("2006-01-02")),
		Username:    "Buncho Bot",
		IconEmoji:   "",
		IconURL:     "https://blog-001.west.edge.storage-yahoo.jp/res/blog-a0-01/galuda6/folder/258481/62/31202762/img_0?1263310483",
		Channel:     "#bot_test",
		Attachments: attachments,
	})

	resp, _ := http.PostForm(
		IncomingURL,
		url.Values{"payload": {string(params)}},
	)

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	fmt.Println(string(body))
}
