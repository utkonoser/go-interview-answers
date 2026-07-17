// Решение с собеса Lamoda (исправлено: длина result, return, Close body).
package parallelcrawl

import (
	"fmt"
	"net/http"
	"sync"
)

var urls = []string{
	"https://www.lamoda.ru/p/mp002xw0lvkd/clothes-tomollyfromjames-plate/",
	"https://www.lamoda.ru/p/mp002xw14uf2/clothes-tomollyfromjames-plate/",
	"https://www.lamoda.ru/p/rtladr746901/clothes-iceberg-plate/",
	"https://www.lamoda.ru/p/mp002xw18h9d/clothes-victoriaveisbrut-plate/",
	"https://www.lamoda.ru/p/mp002xw004x4/clothes-clanvi-plate/",
	"https://www.lamoda.ru/p/mp002xw0zfxy/clothes-glvr-plate/",
	"https://www.lamoda.ru/p/mp002xw0slmg/clothes-snezhnayakoroleva-plate-kozhanoe/",
	"https://www.lamoda.ru/p/mp002xw132c3/clothes-auranna-plate/",
}

type Bundle struct {
	Index int
	Url   string
	Code  int
}

func crawl(urls []string, k int) []int {
	urlChannel := make(chan Bundle, len(urls))
	codeChannel := make(chan Bundle, len(urls))
	result := make([]int, len(urls))

	go func() {
		for i, url := range urls {
			urlChannel <- Bundle{Url: url, Index: i}
		}
		close(urlChannel)
	}()

	var wgCodes sync.WaitGroup
	wgCodes.Add(1)
	go func() {
		defer wgCodes.Done()
		for bundle := range codeChannel {
			result[bundle.Index] = bundle.Code
		}
	}()

	var wg sync.WaitGroup
	wg.Add(k)
	for range k {
		go func() {
			defer wg.Done()
			for url := range urlChannel {
				resp, err := http.Get(url.Url)
				if err != nil {
					fmt.Println(err.Error())
					codeChannel <- url
					continue
				}
				url.Code = resp.StatusCode
				_ = resp.Body.Close()
				codeChannel <- url
			}
		}()
	}
	wg.Wait()
	close(codeChannel)

	wgCodes.Wait()
	return result
}
