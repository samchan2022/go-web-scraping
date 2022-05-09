package helper

import (
	"log"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"golang.org/x/net/html"
)

type Html struct{}

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

func (*Html) GetLinksFromSinglePage(domainUrl string, relPath string) []string{
    u, _ := url.Parse(domainUrl)
    u.Path = path.Join(u.Path, relPath)
    absUrl := u.String()
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

            // Find the relative link
            href = filepath.ToSlash(href)
            r1, _ := regexp.Compile(`^/`)
            isSameDomain:= r1.MatchString(href)
            if isSameDomain {
                linkList = append(linkList, href)
            }

            // Find the link starting with http
            r2, _ := regexp.Compile(`^` + domainUrl)
            isAbsPath:= r2.MatchString(href)
            if isAbsPath{
                u, _ := url.Parse(domainUrl)
                relPath, _ := filepath.Rel(u.String(), href)
                if relPath == "." {
                    relPath = "/"
                }
                relPath = filepath.ToSlash(relPath)
                linkList = append(linkList, relPath)
                continue
            }

        }
    }
}

