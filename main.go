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

    // Max depth of the recursion search
    // 0 means no max depth limit
    maxDepth := 0

    log.Println("Start crawling")
    log.Println("---------------------------------------------------")
    log.Println("domainUrl: ", domainUrl)
    log.Println("workerCount: ",workerCount)
    log.Println("outFile: ", outFile)
    log.Println("maxDepth: ", maxDepth)

    crawl.GetAllDomainLinks(domainUrl, workerCount, outFile, maxDepth)
}
