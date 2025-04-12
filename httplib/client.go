package httplib

import (
	"crypto/tls"
	"github.com/go-resty/resty/v2"
	"time"
)

// Client 是 resty.Client 的别名，表示 HTTP 客户端
type Client = resty.Client

// New 创建并配置一个新的 HTTP 客户端实例
// 返回值为 *Client 类型，表示配置好的 HTTP 客户端
func New() *Client {
	cli := resty.New()
	cli = cli.SetTimeout(time.Minute)                                    // 设置请求超时时间为 1 分钟
	cli = cli.SetRetryCount(3)                                           // 设置最大重试次数为 3 次
	cli = cli.SetRetryWaitTime(time.Second * 3)                          // 设置重试等待时间为 3 秒
	cli = cli.SetRetryMaxWaitTime(time.Second * 15)                      // 设置最大重试等待时间为 15 秒
	cli = cli.SetHeader(`User-Agent`, `SimpleToolkit/1.0`)               // 设置默认 User-Agent 头
	return cli.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // 配置 TLS 跳过证书验证
}
