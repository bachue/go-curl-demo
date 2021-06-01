package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

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
	client.DefaultClient.Client = &http.Client{Transport: transport, Timeout: 180 * time.Second}

	// upload(1, 1)
	upload(100000, 6)

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
					response, err := uploadFileToQiniu(key)
					if err != nil {
						fmt.Printf("goroutineIndex:%d, index:%d, key:%v, response:%v, err:%v\n", goroutineIndex, index, key, response, err)
					} else {
						fmt.Printf("goroutineIndex:%d, index:%d, key:%v, response:%v\n", goroutineIndex, index, key, response)
					}
				}
			}
		}(source, i)
	}

	wait.Wait()
}

func uploadFileToQiniu(key string) (response storage.PutRet, err error) {
	filePath := "/tmp/1m"

	token := "HwFOxpYCQU6oXoZXFOTh1mq5ZZig6Yyocgk3BTZZ:K4CV6KyJiU8anDG9czn-w999xQc=:eyJkZWFkbGluZSI6MTY1NDA1MDY0Niwic2NvcGUiOiIyMDIwLTA2LWNoZWNrYmlsbHMifQ=="

	config := &storage.Config{
		Zone: &storage.Region{
			SrcUpHosts: []string{"up.qiniu.com"},
		},
		Region:   nil,
		UseHTTPS: true,
	}
	ctx := context.Background()

	// uploader := storage.NewResumeUploader(config)
	uploader := storage.NewResumeUploaderV2(config)
	// uploader := storage.NewFormUploader(config);

	err = uploader.PutFile(ctx, &response, token, key, filePath, nil)
	return
}
