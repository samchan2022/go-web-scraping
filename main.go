package main

import (
	"fmt"
	"go/webscraping/helper"
	"log"
	"net/url"
	"os"
	"path"
	"sync"
	"time"
)

// Job for worker
type workerJob struct {
    Cursor string
}

// Result of a worker
type workerResult struct {
    Value string
}

type SyncObj struct {
    mu       sync.Mutex
    Links []string
}

//type LinkObj struct {
    //Parent string
    //Link string
    //Depth int
//}
//type LinkObj struct {
    //Parent string
    //Link string
    //Depth int
//}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func (c *SyncObj) add(name string) bool{
    c.mu.Lock()
    defer c.mu.Unlock()

    //log.Println("===================================================")
    //log.Println(c.Links)
    //log.Println(name)
    if !contains(c.Links, name) {
        c.Links = append(c.Links, name)
        //log.Println("should add")
        return true
    }
    //log.Println("should not add")
    return false
}


func worker(jobs chan workerJob, results chan<- workerResult, wg *sync.WaitGroup, c *SyncObj, domainUrl string) {
    for j := range jobs {
        hrefs := helper.GetLinksFromSinglePage( domainUrl, j.Cursor)

        time.Sleep(time.Millisecond * 100)
        //log.Println("len: links",len(hrefs))
        for _, href := range hrefs{
            //var linkObj LinkObj
            //linkObj.Parent = j.Root
            //linkObj.Link = link

            //u, _ := url.Parse(j.Cursor)
            //u.Path = path.Join(u.Path, href)
            //linkList = append(linkList, u.String())
            //newlink := u.String()

            //log.Println("===================================================")
            //time.Sleep(time.Second* 1)
            //log.Println("root", j.Cursor)
            //log.Println("href", href)
            //log.Println( c.coun00ters)

            //if !helper.TestUrl(newlink){
                //continue
            //}

            u, _ := url.Parse(domainUrl)
            u1, _ := url.Parse(href)
            u.Path = path.Join(u.Path, u1.Path)
            absLink := u.String()

            fmt.Println("absLink", absLink)
            //r1, _ := regexp.Compile(`/`)
            //if r1.MatchString(absLink){
                ////log.Println("abs: ", absLink)
            //}

            if !c.add(href){
                continue
            }

            // Send worker result to result channel
            //---------------------------------------------------
            //u, _ := url.Parse(domainUrl)
            //u.Path = path.Join(u.Path, href)
            r := workerResult{
                Value: u.String(),
            }
            results <- r

            // Create a new job
            //---------------------------------------------------

            newJob := workerJob{
                Cursor: href,
            }

            // Increment the wait group count
            wg.Add(1)
            // Invoke jobs
            go func() {
                jobs <- newJob
            }()
        }
        // Once the job is finished, decrement the wait group count
        wg.Done()
    }
}

func GetAllDomainLinks( rootUrl string, workerCount int, filename string){

    c := SyncObj{
        //Links: map[string]int{"a": 0, "b": 0},
        //Links: []string{rootUrl},
        Links: []string{},
    }

    jobs := make(chan workerJob, workerCount)

    //var csv [][]string
    //header := []string{"Links"}
    //csv = append(csv, header)

    f, err := os.Create(filename)
    if err != nil {
        panic(err)
    }

    // result channel
    results := make(chan workerResult)

    isVisitied := make(chan bool)
    wg := &sync.WaitGroup{}

    // Number of worker count
    for i := 0; i < workerCount; i++ {
        go worker(jobs, results, wg, &c, rootUrl)
    }

    // Initialise the first job
    wg.Add(1)
    go func() {
        jobs <- workerJob{
            Cursor: "",
        }
    }()

    // Wait for all jobs to finish 
    go func() {
        wg.Wait()
        isVisitied <- true
    }()

    loop:
        for {
            select {
            case res := <-results:
                var data []string
                data = append(data, res.Value)
                //csv = append(csv, data)
                f.WriteString(res.Value +"\n")

                //TODO
                //log.Printf(`result=%#v`, res.Value)

            case <-isVisitied:
                log.Printf(`Finished`)
                close(jobs)
                break loop
            }
        }
    //helper.WriteCsv( csv, filename)
    log.Println("count", c.Links)
}

func main(){
    //GetAllDomainLinks("http://localhost:3000", 1,"test.csv")
    sTime := time.Now()
    //GetAllDomainLinks("https://monzo.com", 3, "monzo.csv")
    GetAllDomainLinks("http://go-colly.org/", 10, "colly.csv")
    //GetAllDomainLinks("", 1, "colly.csv")

    elapsedTime := time.Since(sTime)
    log.Println("elapsed Time: ", elapsedTime)
}
