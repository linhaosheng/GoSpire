package entity

import (
	"sync"
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"strconv"
	"net/http"
	"strings"
	"os"
	"io"
)

type ParseEntity struct {
	pageWg      sync.WaitGroup
	picDown     sync.WaitGroup
	picUrl      sync.WaitGroup
	lock        *sync.Mutex
	currPage    int
	entityMap   map[string]DownLoadEntity
	pageUrlChan chan DownLoadEntity
}

type DownLoadEntity struct {
	picUrl string
	title  string
}

func New() *ParseEntity {
	return &ParseEntity{
		currPage:1,
		entityMap:make(map[string]DownLoadEntity),
		pageUrlChan:make(chan DownLoadEntity),
		lock:new(sync.Mutex),
	}
}
//下载所有页面的url ，并且存入entitys
func (pageEn *ParseEntity )DownAllPageUrl(baseUrl string, count int) {
	for i := 1; i <= count; i++ {
		pageEn.pageWg.Add(1)
		pageEn.lock.Lock()
		url := baseUrl + strconv.Itoa(i) + ".htm"
		pageEn.lock.Unlock()
		go DownCurrPage(pageEn, url)
	}
	pageEn.pageWg.Wait()
}

//下载当前页面的所有的图片的url
func DownCurrPage(pageEn*ParseEntity, url string) {
	defer pageEn.pageWg.Done()
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println("parse CurrPage error....")
		return
	}
	doc.Find("div[class=\"box list channel\"]").Find("ul").Find("li").Each(func(i int, s*goquery.Selection) {

		picUrl, exit := s.Find("a").Attr("href")
		if !exit {
			fmt.Println("picUrl error----", exit)
			return
		}
		picUrl = "https://www.dd242.com" + picUrl
		title := s.Find("a").Text()
		fmt.Println("picUrl----" + picUrl)
		fmt.Println("title------" + title)
		/*downLoadEntity := DownLoadEntity{
			picUrl:picUrl,
			title:title,
		}*/
		//pageEn.pageUrlChan <- picUrl
		saveDir := "C:/Users/linhao/Desktop/goPicture/"
		pageEn.picUrl.Add(1)
		go downPictire(pageEn, picUrl, saveDir, title)
		//StartDownLoadPic(pageEn, saveDir)
		/*pageEn.lock.Lock()
		defer pageEn.lock.Unlock()
		pageEn.entityMap[picUrl] = downLoadEntity*/
	})
	pageEn.picUrl.Wait()
}

//开始下载图片
func StartDownLoadPic(pageEn*ParseEntity, baseDir string) {

	for key, downLoadEntity := range pageEn.entityMap {
		url := downLoadEntity.picUrl
		dir := baseDir + downLoadEntity.title + "/"
		<-pageEn.pageUrlChan
		pageEn.pageWg.Add(1)
		go downPictire(pageEn, url, dir, key)
	}
	pageEn.pageWg.Wait()
}

//下载每一个url 对应的图片
func downPictire(pageEn*ParseEntity, url string, dir string, key string) {
	defer pageEn.picUrl.Done()
	/*client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil);
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36)")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate,sdch")
	req.Header.Set("content-type","text/html")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("client downLoad error...." + err.Error())
		return
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)*/
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println("parse PicturePage error...." + err.Error())
		return
	}
	/*str,_:=doc.Html()
	fmt.Println("doc------"+str)
	return*/

	doc.Find(".pics").Find("img").Each(func(i int, s*goquery.Selection) {
		picUrl, exit := s.Attr("src")
		if !exit {
			fmt.Println("picUrl error----", exit)
			return
		}
		pageEn.picDown.Add(1)
		go downLoadPicture(pageEn, picUrl, dir)
	})
	pageEn.picDown.Wait()
	/*pageEn.lock.Lock()
	defer pageEn.lock.Unlock()
	delete(pageEn.entityMap, key)*/
}

//dir  "C:/Users/linhao/Desktop/goPicture/"

func downLoadPicture(pageEn*ParseEntity, picUrl string, dir string) {
	defer pageEn.picDown.Done()
	res, err := http.Get(picUrl)
	if err != nil {
		fmt.Println("down fail----", err)
		return
	}
	i := strings.LastIndex(picUrl, "/")
	name := []rune(picUrl)
	pathName := string(name[i:len(picUrl)])

	error := os.MkdirAll(dir, os.ModeDir)
	if error != nil {
		fmt.Println("create dir fail", error)
	}

	file, err := os.Create(dir + pathName)
	if err != nil {
		fmt.Println("create File fail---", err)
		return
	}
	fmt.Println("save path----", file.Name())
	io.Copy(file, res.Body)
	fmt.Println("done----")
}

