package http

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
	"errors"
)

func GetPageCountNum(url string) (int, error) {
	doc, error := goquery.NewDocument(url)
	if error != nil {
		fmt.Println("parse PageCount error....")
		return 0, error
	}
	pages := doc.Find(".pagination")
	page, exit := pages.Find("a").Last().Attr("href")
	if !exit {
		fmt.Println("get last page error.....")
		return 0, errors.New("cannot find a")
	}
	index := strings.LastIndex(page, ".htm")
	pageCount := [] rune(page)
	count := string(pageCount[0:index])
	countNum, err := strconv.Atoi(count)
	if err != nil {
		fmt.Println("string to int fail.....")
		return 0, error
	}
	return countNum, nil
}


