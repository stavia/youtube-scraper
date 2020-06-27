package scraper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBestResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		query := req.URL.String()
		expected := "/results?search_query=Arctic+Monkeys+-+Do+I+Wanna+Know%3F"
		if query != expected {
			t.Errorf("Expected %s, got %s", expected, query)
		}
		data, _ := ioutil.ReadFile("test-fixtures/test1.html")
		rw.Write([]byte(data))
	}))
	defer server.Close()

	scraper := Scraper{server.URL}
	query := "Arctic Monkeys - Do I Wanna Know?"
	results, _ := scraper.Search(query)
	expected := "https://www.youtube.com/watch?v=bpOSxM0rNPM"
	bestResult := scraper.GetBestResult(query, results)
	if expected != bestResult.Link {
		t.Errorf("Expected %s, got %s", expected, bestResult.Link)
	}
}

func TestNotFoundError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		data, _ := ioutil.ReadFile("test-fixtures/test2.html")
		rw.WriteHeader(404)
		rw.Write([]byte(data))
	}))
	defer server.Close()

	scraper := Scraper{server.URL}
	_, err := scraper.Search("foobar")
	fmt.Println(err.Error())
	expected := "status code error: 404 404 Not Found"
	if expected != err.Error() {
		t.Errorf("Expected %s, got %s", expected, err.Error())
	}
}
