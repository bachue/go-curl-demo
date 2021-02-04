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
		Transport: &http.Transport{},
		//CAPath:     "/Users/senyang/Desktop/QiNiu/Test/Go/test/examples/http-transport/curl/lib/resource/cacert.pem",
		ForceHTTP3: true,
	}
	client.DefaultClient.Client = &http.Client{Transport: transport}

	count := 1
	source := make(chan int, 100)
	go func() {
		for i := 0; i < count; i++ {
			source <- i + 1
		}
		close(source)
	}()

	upload(source, 50)

	fmt.Println("======= Done =======")
}

func upload(source <-chan int, goroutineCount int) {

	filePath := "/Users/senyang/Desktop/QiNiu/pycharm.dmg"
	filePath = "/Users/senyang/Desktop/QiNiu/UploadResource_49M.zip"

	key := "http3_test1"
	token := "HwFOxpYCQU6oXoZXFOTh1mq5ZZig6Yyocgk3BTZZ:6MoNfPe6Tj6LaZXwSmRoY5PqcCA=:eyJzY29wZSI6ImtvZG8tcGhvbmUtem9uZTAtc3BhY2UiLCJkZWFkbGluZSI6MTYxNzUwNzUxMiwgInJldHVybkJvZHkiOiJ7XCJjYWxsYmFja1VybFwiOlwiaHR0cDpcL1wvY2FsbGJhY2suZGV2LnFpbml1LmlvXCIsIFwiZm9vXCI6JCh4OmZvbyksIFwiYmFyXCI6JCh4OmJhciksIFwibWltZVR5cGVcIjokKG1pbWVUeXBlKSwgXCJoYXNoXCI6JChldGFnKSwgXCJrZXlcIjokKGtleSksIFwiZm5hbWVcIjokKGZuYW1lKX0ifQ=="

	wait := &sync.WaitGroup{}
	wait.Add(goroutineCount)

	for i := 0; i < goroutineCount; i++ {
		go func(source <-chan int, filePath, key, token string) {
			defer wait.Done()

			for {
				index, ok := <-source
				if !ok && len(source) == 0 {
					break
				} else {
					key := key + fmt.Sprintf("%d", rand.Int())
					err := uploadFile(filePath, key, token)
					fmt.Printf("index:%d key:%s error:%v \n", index, key, err)
				}
			}
		}(source, filePath, key, token)
	}

	wait.Wait()
}

func uploadFile(filePath, key, token string) error {

	config := &storage.Config{
		Zone: &storage.Region{
			SrcUpHosts: []string{"up.qiniu.com"},
		},
		Region:   nil,
		UseHTTPS: true,
	}
	uploader := storage.NewResumeUploader(config)
	ctx := context.Background()

	return uploader.PutFile(ctx, nil, token, key, filePath, nil)
}
