package crawl

import (
	"go/webscraping/helper"
	"log"
	"net/http"
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
    if !contains(c.Links, name) {
        c.Links = append(c.Links, name)
        return true
    }
    return false
}


func worker(jobs chan workerJob, results chan<- workerResult, wg *sync.WaitGroup, c *SyncObj, domainUrl string) {
    for j := range jobs {
        hrefs := new(helper.Html).GetLinksFromSinglePage( domainUrl, j.Cursor)

        // Set some buffer time to avoid high traffic to the server
        time.Sleep(time.Millisecond * 100)
        for _, href := range hrefs{

            // Get the absolute link
            //---------------------------------------------------
            u, _ := url.Parse(domainUrl)
            u1, _ := url.Parse(href)
            u.Path = path.Join(u.Path, u1.Path)
            absLink := u.String()

            if !c.add(href){
                continue
            }

            // Remove Error / invalid link
            //---------------------------------------------------
            resp, err := http.Get(absLink)

            // test link availability
	        if err != nil {
	        	log.Printf("Failed to visit url %s\n", absLink)
                continue
	        }

            // test Failed response
            if resp.StatusCode > 400 {
                log.Printf("ERROR: Not good %v %v\n", absLink, resp.Status)
                continue
            }

            // Send worker result to result channel
            //---------------------------------------------------
            r := workerResult{
                Value: absLink,
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
        Links: []string{},
    }

    jobs := make(chan workerJob, workerCount)

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
                f.WriteString(res.Value +"\n")
                log.Printf(`result=%#v`, res.Value)
            case <-isVisitied:
                log.Printf(`Finished`)
                close(jobs)
                break loop
            }
        }
}

