package main

import (
	"go/webscraping/crawl"
	"log"
)


func main(){
    // target url
    domainUrl := "http://monzo.com/"
    // Number of threads
    workerCount := 10
    //out file path
    outFile := "out.csv"

    log.Println("Start crawling")
    log.Println("---------------------------------------------------")
    log.Println("domainUrl: ", domainUrl)
    log.Println("workerCount: ",workerCount)
    log.Println("outFile: ", outFile)

    crawl.GetAllDomainLinks(domainUrl, workerCount, outFile)
}
