package main

import (
	"fmt"
	"net/http"
	"time"
)

func checkSite(url string, ch chan string) {
	// 네이버는 의도적으로 3초 지연
	if url == "https://www.naver.com/" {
		time.Sleep(3 * time.Second)
	}

	_, err := http.Get(url)
	if err != nil {
		ch <- url + " ❌ 접속 불가"
		return
	}
	ch <- url + " ✅ 정상 작동"
}

func main() {
	ch := make(chan string)
	sites := []string{
		"https://www.google.com/", // 빠른 응답
		"https://www.naver.com/",  // 느린 응답 (3초)
	}

	// 각 사이트 체크를 고루틴으로 실행
	for _, url := range sites {
		go checkSite(url, ch)
	}

	// select로 타임아웃 처리
	for i := 0; i < len(sites); i++ {
		select {
		case result := <-ch: // 사이트 응답이 도착하면
			fmt.Println(result)

		case <-time.After(2 * time.Second): // 2초 타임아웃
			fmt.Println("⏰ 타임아웃! 응답이 너무 느립니다.")
		}
	}
}
