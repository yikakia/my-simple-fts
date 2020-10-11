package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

var query chan string
var result chan resultlist

func main() {
	query = make(chan string)
	result = make(chan resultlist)

	var dumpPath string
	flag.StringVar(&dumpPath, "p", "enwiki-latest-abstract1.xml.gz", "wiki abstract dump path")
	flag.Parse()

	log.Println("Full-text engine starting...")

	// 在端口 9001 开启 web 服务
	go StartWebService(9001)

	start := time.Now()
	docs, err := loadDocuments(dumpPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded %d documents in %v", len(docs), time.Since(start))

	start = time.Now()
	idx := make(Index)
	idx.add(docs)
	log.Printf("Indexed %d documents in %v", len(docs), time.Since(start))

	// 读取搜索词 并且处理搜索需求
	for {
		que := <-query
		fmt.Println(que)
		matchedIDs := idx.search(que)
		tmpresult := resultlist{docs: make([]Document, 10), num: 0}
		for _, id := range matchedIDs {
			doc := docs[id]
			tmpresult.docs = append(tmpresult.docs, doc)
		}
		tmpresult.num = len(matchedIDs)
		result <- tmpresult
	}

}
