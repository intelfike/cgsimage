package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var numreg = regexp.MustCompile("[0-9]+(-2)?")

func init() {
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(rw, getHTML())
	})
}

func main() {
	log.Print(":8080 server start")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getHTML() string {
	doc, err := goquery.NewDocument("http://imas-cg.net/2017/02/21/50703874.html")
	if err != nil {
		log.Fatal(err)
	}
	// 必要の要素を削除
	doc.Find(`
		#extended, 
		#category, 
		#title,
		#popularArticlesWithImageTagezBcoCDwMMzo0ssf,
		#utilities,
		#comment-form,
		.articleSocial,
		.message-board,
		.google-2ad-m,
		.titleRssBottom
	`).Remove()
	// 全体を左寄せ
	doc.Find("#main").SetAttr("style", "float:left;")

	// 画像を挿入するための要素
	doc.Find(".wrapper").AppendHtml(`<div class="cgs-images" style="display:flex; flex-wrap:wrap; clear:both;"><div>`)

	// 各コメントを編集
	doc.Find(".cmArea").Each(func(_ int, s *goquery.Selection) {
		comment := s.Text()
		comment = strings.NewReplacer("０", "0", "１", "1", "２", "2", "３", "3", "４", "4", "５", "5", "６", "6", "７", "7", "８", "8", "９", "9", "－", "-").Replace(comment)
		numlist := numreg.FindAllString(comment, -1)
		for _, v := range numlist {
			i, err := strconv.Atoi(v)
			// 118-2、146-2などの特例を考慮して計算
			if i > 146 {
				i++
			}
			if i > 118 {
				i++
			}
			if err != nil {
				switch v {
				case "118-2":
					i = 119
				case "146-2":
					i = 148
				default:
					i = 200
				}
				fmt.Println(v)
			}

			// 画像生成
			top, left := calRect(i)
			tmp := imgTemplate
			tmp = strings.Replace(tmp, "toppx", strconv.Itoa(top)+"px", 1)
			tmp = strings.Replace(tmp, "leftpx", strconv.Itoa(left)+"px", 1)
			// 画像追加
			s.Parent().Find(".cgs-images").AppendHtml(tmp)
		}
	})
	html, _ := doc.Html()
	return html
}

// 座標計算
func calRect(i int) (top, left int) {
	i--
	top = 90 * (i / 16)
	left = 30 + (i%16)*120
	return
}

// 画像用のテンプレート
var imgTemplate = `
<div style="
	position: relative;
	overflow:hidden;
	width:60px;
	height:91px;">
<img src="http://livedoor.blogimg.jp/sr_cobra/imgs/4/8/48b6cab9.jpg" style="
	position:absolute;
	top: -toppx;
	left: -leftpx;
">
</div>

`
