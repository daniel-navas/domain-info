package middleware

import (
	"log"

	"github.com/anaskhan96/soup"
)

type TitleAndLogo struct {
	Title string
	Logo  string
}
type TitleAndLogoXtractor struct {
	Get func(string) TitleAndLogo
}

func getTitleAndLogo(url string) [2]string {
	resp, err := soup.Get("https://" + url)
	if err != nil {
		log.Printf("HTTP request failed. %s\n", err)
	}
	doc := soup.HTMLParse(resp)
	title := doc.Find("title")
	logo := doc.Find("link", "rel", "shortcut")
	var titleAndLogo [2]string
	if title.Error == nil {
		titleAndLogo[0] = title.Text()
	} else {
		log.Println(title.Error)
	}
	if logo.Error == nil {
		titleAndLogo[1] = logo.Attrs()["href"]
	} else {
		log.Println(title.Error)
	}
	return titleAndLogo
}

func CreateTitleAndLogoXtractor() *TitleAndLogoXtractor {
	return &TitleAndLogoXtractor{
		Get: func(url string) TitleAndLogo {
			titleAndLogo := getTitleAndLogo(url)
			return TitleAndLogo{
				Title: titleAndLogo[0],
				Logo:  titleAndLogo[1],
			}
		},
	}
}
