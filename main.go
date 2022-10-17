package main

import (
	"encoding/json"
	"github.com/gocolly/colly"
	"log"
	"os"
	"strconv"
	"strings"
)

const SicUrl = "https://www.sec.gov/corpfin/division-of-corporation-finance-standard-industrial-classification-sic-code-list"

type SIC struct {
	SICCode       int    `json:"sicCode"`
	IndustryTitle string `json:"industryTitle"`
	Office        string `json:"office"`
}

func main() {
	c := colly.NewCollector()

	data := make(map[string][]SIC)

	c.OnHTML("table.sic", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			if !strings.Contains(el.Text, "SIC Code") {
				sic := SIC{}

				el.ForEach("td", func(i int, td *colly.HTMLElement) {
					switch i {
					case 0:
						sicCode, err := strconv.Atoi(td.Text)
						if err != nil {
							log.Fatal(err)
						}
						sic.SICCode = sicCode
					case 1:
						sic.Office = td.Text
					case 2:
						sic.IndustryTitle = td.Text
					}
				})
				if _, ok := data[sic.Office]; !ok {
					data[sic.Office] = make([]SIC, 0)
				}

				data[sic.Office] = append(data[sic.Office], sic)
			}
		})
	})

	err := c.Visit(SicUrl)
	if err != nil {
		log.Fatalln(err)
	}

	j, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}

	err = os.WriteFile("sic.json", j, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
