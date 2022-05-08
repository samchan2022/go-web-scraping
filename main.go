package main

import (
    "fmt"
    "net/url"
    "path"
    "regexp"

    "net/http"
    "os"
    "golang.org/x/net/html"
)

// Helper function to pull the href attribute from a Token
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

// Extract all http** links from a given webpage
func crawl(rootUrl string, ch chan string, chFinished chan bool) {
    resp, err := http.Get(rootUrl)
    
    //http.Response.Body
    //The response body is streamed on demand as the Body field

    // what does defer means?
    // go routine
    fmt.Println("---------------------------------------------------")
    fmt.Println("rootUrl", rootUrl)
    defer func() {
        // Notify that we're done after this function
        chFinished <- true
    }()

    if err != nil {
        fmt.Println("ERROR: Failed to crawl:", rootUrl)
        return
    }

    b := resp.Body
    defer b.Close() // close Body when the function completes
    //fmt.Println("---------------------------------------------------")
    //bytes, _ := ioutil.ReadAll(b)
    //fmt.Println("body: ",string(bytes))
    //fmt.Println("---------------------------------------------------")

    z := html.NewTokenizer(b)

    for {
        tt := z.Next()

        switch {
        case tt == html.ErrorToken:
            // End of the document, we're done
            return
        case tt == html.StartTagToken:
            t := z.Token()


            // Check if the token is an <a> tag
            isAnchor := t.Data == "a"
            if !isAnchor {
                continue
            }
            //fmt.Println("---------------------------------------------------")
            //fmt.Println(t.Data)

            // Extract the href value, if there is one
            ok, relUrl := getHref(t)
            fmt.Println("url", relUrl)
            if !ok {
                continue
            }

            // Make sure it sticks at the root
            // use compile to save time / processing power
            r, _ := regexp.Compile(`^/`)
            isSameDomain:= r.MatchString(relUrl)
            if isSameDomain {
                //absoluteUrl := rootUrl + url
                //absoluteUrl := rootUrl + url
                u, _ := url.Parse(rootUrl)
                u.Path = path.Join(u.Path, relUrl)
                //absoluteUrl := path.Join(rootUrl, relUrl)
                //continue
                //ch <- url
                ch <- u.String()
            }

            // Make sure the url begines in http**
            //hasProto := strings.Index(url, "http") == 0
            //if hasProto {
                //ch <- url
            //}

        }
    }
}

func main() {
    foundUrls := make(map[string]bool)
    seedUrls := os.Args[1:]
    fmt.Print("url", seedUrls )
    //return

    // Channels
    chUrls := make(chan string)
    chFinished := make(chan bool) 

    // Kick off the crawl process (concurrently)
    for _, url := range seedUrls {
        go crawl(url, chUrls, chFinished)
    }

    // Subscribe to both channels
    for c := 0; c < len(seedUrls); {
        select {
        case url := <-chUrls:
            foundUrls[url] = true
        case <-chFinished:
            c++
        }
    }

    // We're done! Print the results...

    fmt.Println("\nFound", len(foundUrls), "unique urls:\n")

    for url, _ := range foundUrls {
        fmt.Println(" - " + url)
    }

    close(chUrls)
}
