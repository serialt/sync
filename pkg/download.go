package pkg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Download(src_url, dst_path string) (err error) {
	download_url := src_url

	request, err := http.NewRequest("GET", url.PathEscape(download_url), nil)
	if err != nil {
		return
	}
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36")

	// 设置不验证证书和设置超时时间
	client := &http.Client{
		// Transport: &http.Transport{
		// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		// },
		// Timeout: 30 * time.Second,
	}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("request web failed")
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	ioutil.WriteFile(dst_path, data, 0644)
	return
}
