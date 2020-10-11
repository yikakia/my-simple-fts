package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	pageTop = `<!DOCTYPE HTML><html><head>
<style>.error{color:#FF0000;}</style></head><title>简易全文搜索</title>
<body><h3>作者：Yika</h3>
<p>简易的全文搜索系统</p>
<a href=https://github.com/yikakia/simplefts>源码地址</a><br />`
	form = `<form action="/" method="POST">
<label for="words">搜索词 (单个或多个英文的关键词，通过空格或逗号（英文）隔开):</label><br />
<input type="text" name="query" size="30"><br />
<input type="submit" value="搜索">
</form>`
	pageBottom = `</body></html>`
	anError    = `<p class="error">%s</p>`
)

type resultlist struct {
	docs []Document
	num  int
}

// StartWebService 用于开启web服务 参数需要是 int 类型 范围是0~65535
func StartWebService(port int) {
	if port < 0 || port > 65535 {
		log.Fatal("port out of range(0,65535)")
	}
	strport := fmt.Sprintf(":%d", port)
	http.HandleFunc("/", homePage)
	log.Printf("Initialization service at http://localhost%s\n", strport)
	if err := http.ListenAndServe(strport, nil); err != nil {
		log.Fatal("failed to start server", err)
	}
}
func homePage(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm() // Must be called before writing response
	fmt.Fprint(writer, pageTop, form)
	if err != nil {
		fmt.Fprintf(writer, anError, err)
	} else {
		if keywords, message, ok := processRequest(request); ok {
			fmt.Fprint(writer, formatResult(keywords))
		} else if message != "" {
			fmt.Fprintf(writer, anError, message)
		}
	}
	fmt.Fprint(writer, pageBottom)
}

// 用于过滤请求
func processRequest(request *http.Request) (string, string, bool) {
	keywords := ""
	message := ""
	if slice, found := request.Form["query"]; found && len(slice) > 0 {
		text := strings.Replace(slice[0], ",", " ", -1)
		for _, keyword := range text {
			keywords = keywords + string(keyword)
		}
	}
	if len(keywords) == 0 {
		return keywords, "", false // 没有搜索请求的时候显示的样子
	}
	return keywords, message, true
}

func formatResult(que string) string {
	query <- que
	starttime := time.Now()
	re := <-result
	returnstr := ""
	for _, doc := range re.docs {
		if doc.URL != "" {
			returnstr += fmt.Sprintf("<p>标题：%s</p>"+
				"简介：%s<br />"+
				"链接：<a href=\"%s\">%s</a><br />",
				doc.Title, doc.Text, doc.URL, doc.URL)
		}
	}
	returnstr = fmt.Sprintf(
		"搜索词为：%s<br />"+
			"搜索时间为：%v<br />"+
			"总共搜索到了%d个页面<br />",
		que, time.Since(starttime), re.num) +
		returnstr
	return returnstr
}
