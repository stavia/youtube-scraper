package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/stavia/youtube-scraper/scraper"
)

func main() {
	arg1 := flag.String("query", "", "Search query. The suggested format is: Band - Song. For example: \"Arctic Monkeys - Do I Wanna Know?\"")
	flag.Parse()
	if *arg1 == "" {
		flag.PrintDefaults()
		return
	}
	query := strings.TrimSpace(*arg1)
	scraper := scraper.Scraper{}
	results, err := scraper.Search(query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Print results")
	for _, result := range results {
		fmt.Println(result)
	}
	bestResult := scraper.GetBestResult(query, results)
	fmt.Printf("\nBest result: %s %s\n", bestResult.Title, bestResult.Link)
}
