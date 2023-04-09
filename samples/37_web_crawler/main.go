package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Book struct {
	Name    string `json:"name"`
	Price   string `json:"price"`
	InStock string `json:"instock"`
	Rate    int    `json:"rate"`
	Link    string `json:"link"`
}

type RatingEnum int

const (
	One RatingEnum = iota
	Two
	Three
	Four
	Five
)

var ratingEnumMap = map[string]RatingEnum{
	"One":   One,
	"Two":   Two,
	"Three": Three,
	"Four":  Four,
	"Five":  Five,
}

func (e RatingEnum) Int() int {
	return int(e) + 1
}

func getRateFromString(value string) int {
	e, ok := ratingEnumMap[value]
	if !ok {
		return 0
	}

	return e.Int()
}

func getRateFromClass(input string) string {
	parts := strings.Split(input, " ")
	return parts[1]
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v", name, time.Since(start))
	}
}

func main() {
	defer timer("main")()
	c := colly.NewCollector(colly.Async(true))
	books := []Book{}
	startPage := "https://books.toscrape.com/catalogue/category/books/travel_2/index.html"

	// Go to categories from left sidebar
	c.OnHTML("div.side_categories li ul li", func(h *colly.HTMLElement) {
		// Most websites have blocking function so a waiting time is necessary

		link := h.ChildAttr("a", "href")
		time.Sleep(time.Millisecond * 2)
		c.Visit(h.Request.AbsoluteURL(link))
	})

	// Goto the next pages
	c.OnHTML("li.next a", func(h *colly.HTMLElement) {
		time.Sleep(time.Millisecond * 2)
		c.Visit(h.Request.AbsoluteURL(h.Attr("href")))
	})

	// Get books
	c.OnHTML("article.product_pod", func(h *colly.HTMLElement) {
		// Get rate from <p class=star-rating Two>
		rateStr := getRateFromClass(h.ChildAttr("p[class^=star-rating]", "class"))
		book := Book{
			Name:    h.ChildAttr("h3 a", "title"),
			Price:   h.ChildText("p.price_color"),
			InStock: h.ChildText("p.instock"),
			Rate:    getRateFromString(rateStr),
			Link:    h.ChildAttr("a", "href"),
		}

		books = append(books, book)
	})

	// ==== Custom rules ====
	// c.Limit(
	// 	&colly.LimitRule{
	// 		Parallelism: 2,
	// 		// Delay:       time.Second * 2,
	// 	},
	// )

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visit -", r.URL)
	})

	err := c.Visit(startPage)
	if err != nil {
		log.Fatal(err)
	}

	c.Wait()
	time.Sleep(time.Microsecond * 2)

	data, err := json.MarshalIndent(books, "\t", "")
	if err != nil {
		log.Fatal(err)
	}

	// Write crawled books to file
	os.WriteFile("books.json", data, 0644)
}
