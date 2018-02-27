/*
爬取内涵段子
//明确目标
//第1页  https://www.neihanba.com/dz/index.html
//第2页  https://www.neihanba.com/dz/list_2.html
//第n页 https://www.neihanba.com/dz/list_n.html

步骤:1. hettpGet首页的url
	2. 解析某个段子的url的正则 `<h4> <a href="(.*?)"`
	3. 拼接完整的段子的url  https://www.neihanba.com + /dz/1092886.html
	4. hettpGet某个段子的url
	5. 解析某个段子标题的正则 `<h1>(?:s(.*?))</h1>`
	6. 解析这个段子正文的正则 `<td><p>(?s:(.*?))</p></td>`
*/
package main

import (
	"fmt"
	iconv "github.com/djimenez/iconv-go"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

type mySpider struct {
	m_pageNum   int
	m_pageUrl   string
	m_duanziUrl string
}

func (this *mySpider) Spider_page() {
	fmt.Println("正在爬取第" + strconv.Itoa(this.m_pageNum) + " 页")
	var content, pageUrl string
	if this.m_pageNum == 1 {
		pageUrl = "https://www.neihanba.com/dz/index.html"
	} else {
		pageUrl = fmt.Sprintf("https://www.neihanba.com/dz/list_%d.html", this.m_pageNum)

	}
	content, _ = this.HttpGet(pageUrl)
	duanziUrlReg := regexp.MustCompile(`<h4> <a href="(.*?)"`)
	for _, match := range duanziUrlReg.FindAllStringSubmatch(content, -1) {
		duanziUrl := "https://www.neihanba.com" + match[1]
		fmt.Printf("%s\n", duanziUrl)
		//this.m_duanziUrl = duanziUrl
		this.Spider_duanzi(duanziUrl)
	}
}

func (this *mySpider) Spider_duanzi(url string) {
	content, _ := this.HttpGet(url)
	titleReg := regexp.MustCompile(`<h1>(.*?)</h1>`)
	contentReg := regexp.MustCompile(`<td><p>(?s:(.*?))</p></td>`)
	for _, match := range titleReg.FindAllStringSubmatch(content, -1) {
		matchFile, _ := iconv.ConvertString(match[1], "gb2312", "utf-8")
		//fmt.Println("create file " + matchFile)
		fout, err := os.Create("./" + matchFile + ".txt")
		if err != nil {
			log.Println("create file err")
		}
		for _, duanziContent := range contentReg.FindAllStringSubmatch(content, -1) {
			matchContent, _ := iconv.ConvertString(duanziContent[1], "gb2312", "utf-8")
			fout.WriteString(matchContent)
		}
	}
}

func (this *mySpider) HttpGet(url string) (content string, status int) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("get url err ![URL: %s]\n", url)
		content = ""
		status = -100
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("read url err ![URL: %s]\n", url)
		content = ""
		status = resp.StatusCode
		return
	}
	//成功拿到内容 编码格式为源网页格式  -->转换成utf-8
	out := make([]byte, len(data))
	out = out[:]
	iconv.Convert(data, out, "gb2312", "utf-8")

	content = string(data)
	status = resp.StatusCode
	return
}

func main() {
	s := new(mySpider)
	s.m_pageNum = 1

	var cmd string
	for {
		fmt.Println("按任意键爬取下一页,按exit退出!")
		fmt.Scanf("%s", &cmd)
		if cmd == "exit" {
			break
		}
		s.Spider_page()
		s.m_pageNum++
	}
}
