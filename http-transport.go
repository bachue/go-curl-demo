package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
	//      "crypto/tls"

	"github.com/YangSen-qn/go-curl/v2/curl"
	"github.com/qiniu/go-sdk/v7/client"
	"github.com/qiniu/go-sdk/v7/storage"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	transport := &curl.Transport{
		Transport:      &http.Transport{},
		ForceHTTP3:     true,
		HTTP3LogEnable: true,
	}
	//      transport := http.DefaultTransport.(*http.Transport)
	//      transport.ForceAttemptHTTP2 = false
	//      transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client.DefaultClient.Client = &http.Client{Transport: transport, Timeout: 180 * time.Second}

	upload(1, 1)
	//	upload(10000, 16)

	fmt.Println("======= Done =======")
}

func upload(uploadCount int, goroutineCount int) {

	source := make(chan int, 100)
	go func() {
		for i := 0; i < uploadCount; i++ {
			source <- i + 1
		}
		close(source)
	}()

	wait := &sync.WaitGroup{}
	wait.Add(goroutineCount)

	for i := 0; i < goroutineCount; i++ {
		go func(source <-chan int, goroutineIndex int) {
			defer func(goroutineIndex int) {
				fmt.Printf("== goroutineIndex:%d done \n", goroutineIndex)
				wait.Done()
			}(goroutineIndex)

			for {
				index, ok := <-source
				if !ok {
					break
				} else {
					key := "http3_test_" + time.Now().Format("2006/01/02 15:04:05.999999")
					fmt.Printf("goroutineIndex:%d, index:%d, key:%v\n", goroutineIndex, index, key)
					beginTime := time.Now()
					response, err := uploadFileToQiniu(key)
					elapsed := time.Since(beginTime)
					if err != nil {
						fmt.Printf("goroutineIndex:%d, index:%d, key:%v, response:%v, elapsed:%v, err:%v\n", goroutineIndex, index, key, response, elapsed, err)
					} else {
						fmt.Printf("goroutineIndex:%d, index:%d, key:%v, response:%v, elapsed:%v\n", goroutineIndex, index, key, response, elapsed)
					}
				}
			}
		}(source, i)
	}

	wait.Wait()
}

func uploadFileToQiniu(key string) (response storage.PutRet, err error) {
	filePath := "/tmp/1m"

	// token := "0tf5awMVxwf8WrEvrjtbiZrdRZRJU-91JgCqTOC8:6DVqykXJUW8YBrB1FA90NpT5ybc=:eyJkZWFkbGluZSI6MTY2MzA1MTUwOCwic2NvcGUiOiJ6aG91cm9uZy1odHRwMy10ZXN0In0="
	token := "0p6YmsTFa-MLL4cZhe2Yj8n7nXZ3N_wk0ZmuZMnu:9Ws_xYpJN-4DQLqyAgKI0iemEw0=:eyJkZWFkbGluZSI6MTY2Njg3MjY2Mywic2NvcGUiOiJ6MC1idWNrZXQifQ=="

	config := &storage.Config{
		Zone: &storage.Region{
			SrcUpHosts: []string{"up-z0.qbox.me"},
		},
		Region:   nil,
		UseHTTPS: true,
	}
	ctx := context.Background()

	if err = storage.NewResumeUploader(config).PutFile(ctx, &response, token, key, filePath, nil); err != nil {
		return
	}

	if err = storage.NewResumeUploaderV2(config).PutFile(ctx, &response, token, key, filePath, nil); err != nil {
		return
	}

	if err = storage.NewFormUploader(config).PutFile(ctx, &response, token, key, filePath, nil); err != nil {
		return
	}

	return
}
