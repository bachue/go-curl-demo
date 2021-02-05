module github.com/yangjiaxuan/http-transport/demo

go 1.15

require (
	github.com/YangSen-qn/go-curl/v2 v2.0.6
	github.com/qiniu/go-sdk/v7 v7.9.0
)

replace (
	github.com/YangSen-qn/go-curl/v2 => ../go-curl
	github.com/yangjiaxuan/http-transport/demo => ./
)
