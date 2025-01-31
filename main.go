package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type item struct {
	Name      string `json:"name"`
	Price     string `json:"price"`
	Link      string `json:"link"`
	OwnerName string `json:"ownername"`
	Location  string `json:"location"`
}

var items []item

func main() {
	// Initialize Colly collector
	c := colly.NewCollector(
		colly.AllowedDomains("uae.dubizzle.com"),
	)

	// Extract each item within the list
	c.OnHTML("div[data-testid=lpv-list] [data-testid^=listing-]", func(h *colly.HTMLElement) {
		// Extract name and price
		name := h.ChildText("h2[data-testid=subheading-text]")
		price := h.ChildText("div[data-testid=listing-price]")
		link := h.Request.AbsoluteURL(h.Attr("href"))

		// Trim spaces from name and price
		name = strings.TrimSpace(name)
		price = strings.TrimSpace(price)
		link = strings.TrimSpace(link)

		// Validate name and price before adding to items
		if name != "" && price != "" && link != "" {
			// Create an item struct and append to items slice
			item := item{Name: name, Price: price, Link: link}
			items = append(items, item)
		}

		c.Visit(link)
	})

	c.OnHTML("p[data-testid=agent-name]", func(h *colly.HTMLElement) {

		ownername := h.Text

		link := h.Request.URL.String()
		for i := range items {
			if items[i].Link == link {
				items[i].OwnerName = strings.TrimSpace(ownername)
				break
			}
		}

	})

	c.OnHTML("div[data-testid=listing-location-map]", func(h *colly.HTMLElement) {

		location := h.Text
		link := h.Request.URL.String()
		for i := range items {
			if items[i].Link == link {
				items[i].Location = strings.TrimSpace(location)
				break
			}
		}

	})

	c.OnHTML("[data-testid=page-next]", func(h *colly.HTMLElement) {
		next_page := h.Request.AbsoluteURL(h.Attr("href"))
		c.Visit(next_page)

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL.String())
	})

	// Start scraping
	var url string
	fmt.Println("Enter The Link Of The Listings You Want To Scrape")
	fmt.Scanln(&url)
	c.Visit(url)

	content, err := json.Marshal(items)
	if err != nil {
		fmt.Println(err.Error())
	}
	os.WriteFile("products.json", content, 0644)

	// Print items as JSON
	itemsJSON, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling items:", err)
		return
	}
	fmt.Println(string(itemsJSON))
}
