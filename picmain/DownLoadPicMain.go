package main

import (
	"GoSpire/http"
	"fmt"
	"GoSpire/entity"
)

func main() {
	url := "https://www.dd242.com/htm/piclist9/1.htm"
	baseUrl := "https://www.dd242.com/htm/piclist9/"
	//saveDir := "C:/Users/linhao/Desktop/goPicture/"
	pageNum, err := http.GetPageCountNum(url)
	if err != nil {
		fmt.Println("get All Picture Fail.....")
		return
	}
	fmt.Println("page-------",pageNum)
	parseEntity := entity.New()
	parseEntity.DownAllPageUrl(baseUrl, pageNum)
}