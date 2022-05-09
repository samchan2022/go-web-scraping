package helper

import (
    "encoding/csv"
    "log"
    "os"
)

type File struct{}

func (*File) ReadCsv(path string) [][]string {
    csvFile, err := os.Open( path )
    if err != nil {
        log.Println(err)
    }
    defer csvFile.Close()
    csvReader := csv.NewReader(csvFile)
    data, err := csvReader.ReadAll()
    if err != nil {
        log.Println(err)
    }
    return data
}

func (*File) WriteCsv(rows [][]string, filepath string){
    csvfile, err := os.Create(filepath)
 
    if err != nil {
        log.Printf("failed creating file: %s\n", err)
    }
 
    csvwriter := csv.NewWriter(csvfile)
 
    for _, row := range rows {
        _ = csvwriter.Write(row)
    }
 
    csvwriter.Flush()
 
    csvfile.Close()
}

