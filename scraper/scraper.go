package scraper

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mozillazg/go-slugify"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

const YoutubeUrl = "https://www.youtube.com"

type SearchResult struct {
	Title string
	Link  string
	Views int
}

type Scraper struct {
	BaseUrl string
}

func (scraper Scraper) Search(query string) (results []SearchResult, err error) {
	if scraper.BaseUrl == "" {
		scraper.BaseUrl = YoutubeUrl
	}
	res, err := http.Get(scraper.BaseUrl + "/results?search_query=" + url.QueryEscape(query))
	if err != nil {
		return results, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return results, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return results, err
	}
	return scraper.GetLinks(doc), nil
}

func (Scraper) GetLinks(doc *goquery.Document) (results []SearchResult) {
	var result SearchResult
	doc.Find(".yt-lockup-content").Each(func(i int, selection *goquery.Selection) {
		title, href := getTitleAndHref(selection)
		views := getViews(selection)
		if title == "" || href == "" {
			return
		}
		result.Title = title
		result.Link = YoutubeUrl + href
		result.Views = views
		results = append(results, result)
	})
	return results
}

func (Scraper) GetBestResult(query string, results []SearchResult) (bestResult SearchResult) {
	var distance int
	query = slugify.Slugify(query)
	for _, result := range results {
		slugTitle := slugify.Slugify(result.Title)
		slugTitle = strings.Replace(slugTitle, "-official-audio", "", 1)
		slugTitle = strings.Replace(slugTitle, "-official-video", "", 1)
		distance = levenshtein.DistanceForStrings([]rune(query), []rune(slugTitle), levenshtein.DefaultOptions)
		if distance <= 5 {
			bestResult = result
			break
		}
	}
	return bestResult
}

func getTitleAndHref(selection *goquery.Selection) (title string, href string) {
	anchor := selection.Find(".yt-lockup-title").Find("a").First()
	return anchor.AttrOr("title", ""), anchor.AttrOr("href", "")
}

func getViews(selection *goquery.Selection) (views int) {
	lists := selection.Find(".yt-lockup-meta-info").Find("li")
	if lists.Length() != 2 {
		return views
	}
	text := lists.Eq(1).Text()
	match, err := regexp.MatchString(`(?m)\d* \w*`, text)
	if !match || err != nil {
		return views
	}
	text = strings.Split(text, " ")[0]
	views, err = strconv.Atoi(strings.ReplaceAll(text, ".", ""))
	if err != nil {
		return views
	}
	return views
}
