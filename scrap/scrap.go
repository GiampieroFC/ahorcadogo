package scrap

import (
	"github.com/gocolly/colly"
	"net/url"
)

type Paldef struct {
	Palabra    string
	Definicion string
	Link       *url.URL
}

func PrintPalabra() Paldef {

	var p string
	var def string
	var link *url.URL

	url := "https://es.wikipedia.org/wiki/Especial:Aleatoria"

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11"),
	)
	c.AllowURLRevisit = true

	c.OnHTML("span.mw-page-title-main", func(h *colly.HTMLElement) {
		// fmt.Println("🗣 palabra: ", h.Text)
		p = h.Text
	})
	c.OnHTML("#mw-content-text > div:nth-child(1) > p:first-of-type", func(h *colly.HTMLElement) {
		// fmt.Println("📝 definición: ", h.Text)
		def = h.Text
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Add("Accept", "*/*")
		// fmt.Println("entrando...")
	})
	c.OnError(func(_ *colly.Response, err error) {
		// fmt.Println("Something went wrong:", err)
	})
	c.OnResponse(func(r *colly.Response) {
		// fmt.Println("Visited", r.Request.URL)
		link = r.Request.URL
	})

	c.Visit(url)

	palabra := Paldef{
		Palabra:    p,
		Definicion: def,
		Link:       link,
	}

	return palabra

}
