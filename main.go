package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

type item struct {
	Name   string `json:"name"`
	Price  string `json:"price"`
	ImgUrl string `json:"imgurl"`
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("uae.dubizzle.com"),
	)
	c.OnHTML("div[data-testid=listing-price]", func(h *colly.HTMLElement) {
		fmt.Println(h.Text)
	})

	//

	c.Visit("https://uae.dubizzle.com/search/?keyword=s24")
}
