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
	// trim
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
	// left shift
	doc.Find("#main").SetAttr("style", "float:left;")

	cmArea := doc.Find(".cmArea")
	cmArea.Each(func(i int, s *goquery.Selection) {
		comment := s.Text()
		numreper := strings.NewReplacer("０", "0", "１", "1", "２", "2", "３", "3", "４", "4", "５", "5", "６", "6", "７", "7", "８", "8", "９", "9", "ー", "-")
		comment = numreper.Replace(comment)
		numlist := numreg.FindAllString(comment, -1)
		// replist := make([]string, len(numlist)<<1)
		s.Parent().AppendHtml(`<div class="cgs-images" style="display:flex; flex-wrap:wrap; clear:both;"><div>`)
		for _, v := range numlist {
			i, err := strconv.Atoi(v)
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
					i = 147
				default:
					i = 0
				}
			}

			top, left := calRect(i)
			tmp := imgTemplate
			tmp = strings.Replace(tmp, "toppx", strconv.Itoa(top)+"px", 1)
			tmp = strings.Replace(tmp, "leftpx", strconv.Itoa(left)+"px", 1)

			s.Parent().Find(".cgs-images").AppendHtml(tmp)
			// replist[n<<1] = v
			// replist[n<<1+1] = v + tmp
		}

		// reper := strings.NewReplacer(replist...)
		// s.SetHtml(reper.Replace(comment))
	})
	html, _ := doc.Html()
	return html
}

func calRect(i int) (top, left int) {
	i--
	top = 90 * (i / 16)
	left = 30 + (i%16)*120
	return
}

// var imgTemplate = `
// <img src="http://livedoor.blogimg.jp/sr_cobra/imgs/4/8/48b6cab9.jpg" style="
// 	position:absolute;
// 	clip: rect(toppx rightpx bottompx leftpx);
// ">
//  `

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
