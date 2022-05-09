package helper

import (
	"log"
	"net/http"
	//"net/url"
	//"path"
	"regexp"

	"golang.org/x/net/html"
)

func getHref(t html.Token) (ok bool, href string) {
    // Iterate over token attributes until we find an "href"
    for _, a := range t.Attr {
        if a.Key == "href" {
            href = a.Val
            ok = true
        }
    }
    // "bare" return will return the variables (ok, href) as 
    // defined in the function definition
    return
}

func TestUrl(url string) bool{
	resp, _ := http.Get(url)
	if resp.StatusCode > 400 {
		return false
	}
    return true
}

func testStatus(url string, ch chan string, chFinished chan bool) {
	resp, err := http.Get(url)

	defer func() {
		// Notify that we're done after this function
		chFinished <- true
	}()

	if err != nil {
		log.Printf("ERROR: Failed to check  %v %v\n", url, err)
		return
	}
	if resp.StatusCode > 400 {
		//fmt.Printf("ERROR: Not good %v %v\n", url , resp.Status)
		return
	}

	ch <- url
}

func testUrls(urls []string) []string {
	testUrls := make(chan string, 5)
	testFinished := make(chan bool, 5)

	for _, url := range urls {
		go func(v string) {
			testStatus(v, testUrls, testFinished)
		}(url)
	}

	var goodUrls []string

	for c := 0; c < len(urls); {
		select {
		case url := <-testUrls:
			//goodUrls[url] = true
            goodUrls = append(goodUrls, url)
		case <-testFinished:
			c++
		}
	}
	log.Println("\nFound", len(goodUrls), "good urls:\n")
	close(testUrls)
	close(testFinished)
	return goodUrls
}


func GetLinksFromSinglePage(domainUrl string, cursor string) []string{
    absUrl := domainUrl + cursor
    resp, err := http.Get(absUrl)
    var linkList []string

    if err != nil {
        log.Println("ERROR: Failed to crawl:", absUrl)
        return linkList
    }

    b := resp.Body
    defer b.Close() // close Body when the function completes

    z := html.NewTokenizer(b)

    for {
        tt := z.Next()

        switch {
        case tt == html.ErrorToken:
            // End of the document, we're done
            return linkList
        case tt == html.StartTagToken:
            t := z.Token()

            // Check if the token is an <a> tag
            isAnchor := t.Data == "a"
            if !isAnchor {
                continue
            }

            // Extract the href value, if there is one
            ok, href := getHref(t)
            if !ok {
                continue
            }

            // Make sure it sticks at the root
            // use compile to save time / processing power
            //r, _ := regexp.Compile(`^/`)
            r, _ := regexp.Compile(`^/|^` + domainUrl + `.+`)
            isSameDomain:= r.MatchString(href)
            if isSameDomain {
                //u, _ := url.Parse(rootUrl)
                //u.Path = path.Join(u.Path, relUrl)
                //linkList = append(linkList, u.String())
                linkList = append(linkList, href)
            }
        }
    }
}

