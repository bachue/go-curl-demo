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


	upload(1000, 50)

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
				if !ok && len(source) == 0 {
					break
				} else {
					uploadFileToQiniu(index, goroutineIndex)
				}
			}
		}(source, i)
	}

	wait.Wait()
}

func uploadFileToQiniu(index int, goroutineIndex int) {

	filePath := "/Users/senyang/Desktop/QiNiu/pycharm.dmg"
	filePath = "/Users/senyang/Desktop/QiNiu/UploadResource_49M.zip"
	filePath = "/Users/senyang/Desktop/QiNiu/Image/image.png"

	key := "http3_test1" + fmt.Sprintf("%d", rand.Int())
	token := "HwFOxpYCQU6oXoZXFOTh1mq5ZZig6Yyocgk3BTZZ:6MoNfPe6Tj6LaZXwSmRoY5PqcCA=:eyJzY29wZSI6ImtvZG8tcGhvbmUtem9uZTAtc3BhY2UiLCJkZWFkbGluZSI6MTYxNzUwNzUxMiwgInJldHVybkJvZHkiOiJ7XCJjYWxsYmFja1VybFwiOlwiaHR0cDpcL1wvY2FsbGJhY2suZGV2LnFpbml1LmlvXCIsIFwiZm9vXCI6JCh4OmZvbyksIFwiYmFyXCI6JCh4OmJhciksIFwibWltZVR5cGVcIjokKG1pbWVUeXBlKSwgXCJoYXNoXCI6JChldGFnKSwgXCJrZXlcIjokKGtleSksIFwiZm5hbWVcIjokKGZuYW1lKX0ifQ=="

	config := &storage.Config{
		Zone: &storage.Region{
			SrcUpHosts: []string{"up.qiniu.com"},
		},
		Region:   nil,
		UseHTTPS: true,
	}
	uploader := storage.NewResumeUploader(config)
	ctx := context.Background()

	var response storage.PutRet
	err := uploader.PutFile(ctx, &response, token, key, filePath, nil)
	fmt.Printf("goroutineIndex:%d, index:%d key:%s error:%v response:%v \n", goroutineIndex, index, key, err, response)
}
