package main

import (
	"go/webscraping/helper"
	"go/webscraping/crawl"
	"log"
	"net/http"
	"regexp"
	"testing"
)

func TestUrl(t *testing.T) {
    domainUrl := "http://monzo.com/"
    workerCount := 10
    outFile := "out.csv"

    crawl.GetAllDomainLinks(domainUrl, workerCount, outFile)

    data := new(helper.File).ReadCsv(outFile)
    for _, row := range data {
        log.Println("link: ", row[0])
        link := row[0]
        resp, err := http.Get(link)

        // test link availability
	    if err != nil {
	    	t.Errorf("Failed to visit url %s", link)
	    }

        // test Failed response
        if resp.StatusCode > 400 {
            t.Errorf("ERROR: Not good %v %v\n", link, resp.Status)
        }
        // Test same domain
        r, _:= regexp.Compile(`^` + domainUrl)
        if !r.MatchString(link){
            t.Errorf("ERROR: Link is not in the root domain %s", link)
        }
    }
}
