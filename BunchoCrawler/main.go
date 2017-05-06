package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
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
	target := "http://pets-kojima.com/small_list/?topics_group_id=4&group=&shop%5B%5D=tokyo01&freeword=%E3%83%96%E3%83%B3%E3%83%81%E3%83%A7%E3%82%A6&price_bottom=&price_upper=&order_type=2"
	doc, err := goquery.NewDocument(target)
	if err != nil {
		fmt.Println(err)
	}

	var attachments []Attachment
	doc.Find("div.sca_table2").Each(func(_ int, s *goquery.Selection) {
		// title := s.Find("td.sca_name2 a").Text()
		attachment := Map(s)
		attachments = append(attachments, attachment)
	})

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
