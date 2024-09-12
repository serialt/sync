package main

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"log/slog"

	"github.com/google/go-github/v64/github"
	"github.com/serialt/crab"
)

var ctx = context.Background()

// Run 下载release
// 目前有两种方法结合处理长时间下载导致的限流
// 1) 使用多个token进行切换
// 2）随机 sleep n 秒,减小请求频率
func Run() {
	c := &GithubClient{}
	c.NextClient()
	var releases []*Release
	for _, v := range config.GithubRelease {
		time.Sleep(time.Second * time.Duration(rand.Intn(int(config.RandomSleep))))
		c.NextClient()
		slog.Info("Get release data", "repo", v)
		myMonitorRepo := strings.Split(v, "/")
		if config.LastNum <= 1 {
			releasesC, err := c.GetLastestRelease(myMonitorRepo[0], myMonitorRepo[1])
			if err != nil {
				continue
			}
			releases = crab.SliceMerge(releases, releasesC)
		} else {
			releasesC, err := c.ListRelease(myMonitorRepo[0], myMonitorRepo[1], int(config.LastNum))
			if err != nil {
				continue
			}
			releases = crab.SliceMerge(releases, releasesC)
		}

	}

	for _, v := range releases {
		time.Sleep(time.Second * time.Duration(rand.Intn(int(config.RandomSleep))))
		c.NextClient()
		c.DownloadReleaseAsset(config.MirrorRoot, v)
	}

}

func (c *GithubClient) NextClient() {
	if len(config.GithubToken) == 1 {
		if c.Client == nil {
			c.Client = github.NewClient(nil).WithAuthToken(config.GithubToken[0])
		}
	} else {
		index := int(c.TokenIndex) % (len(config.GithubToken) - 1)
		c.Client = github.NewClient(nil).WithAuthToken(config.GithubToken[index])
		index++
		c.TokenIndex = int64(index)
	}

}

// GetLastestRelease 获取最新的稳定release
func (c *GithubClient) GetLastestRelease(owner, repo string) (items []*Release, err error) {
	release, _, err := c.Client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		slog.Info("cat not get the lastest release",
			"owner", owner,
			"repo", repo,
			"error", err)
		return
	}

	for _, v := range release.Assets {
		// 如果名字中包含过滤字符
		if ExcludeTxt(*v.Name) {
			continue
		}
		items = append(items, &Release{
			Owner:     owner,
			Repo:      repo,
			Version:   *release.TagName,
			AssetName: *v.Name,
			AssetID:   *v.ID,
		})

	}
	return
}

// ListRelease 获取最近 lastNum 个数的release,可能包括beta和pre-release
func (c *GithubClient) ListRelease(owner, repo string, lastNum int) (items []*Release, err error) {
	opt := &github.ListOptions{Page: 1, PerPage: lastNum}
	releaseList, _, err := c.Client.Repositories.ListReleases(ctx, owner, repo, opt)
	if err != nil {
		slog.Error("cat not get the release",
			"owner", owner,
			"repo", repo,
			"lastNum", lastNum,
			"error", err)
		return
	}
	for _, release := range releaseList {
		for _, v := range release.Assets {
			if ExcludeTxt(*v.Name) {
				continue
			}
			items = append(items, &Release{
				Owner:     owner,
				Repo:      repo,
				Version:   *release.TagName,
				AssetName: *v.Name,
				AssetID:   *v.ID,
			})

		}
	}
	return
}

// DownloadReleaseAsset 根据assetID 下载release到指定的目录
func (c *GithubClient) DownloadReleaseAsset(dir string, item *Release) {
	rFilePath := fmt.Sprintf("%s/%s/%s/%s", dir, item.Repo, item.Version, item.AssetName)
	rDir := fmt.Sprintf("%s/%s/%s", dir, item.Repo, item.Version)

	if !IsDirExists(rDir) {
		Mkdir(rDir)
	}
	if IsDirExists(rFilePath) {
		return
	}
	reader, _, err := c.Client.Repositories.DownloadReleaseAsset(ctx, item.Owner, item.Repo, item.AssetID, http.DefaultClient)
	if err != nil {
		slog.Error(
			"Download release asset failed",
			"owner", item.Owner,
			"repo", item.Repo,
			"version", item.Version,
			"name", item.AssetName,
			"error", err,
		)
	}
	context, err := io.ReadAll(reader)
	if err != nil {
		slog.Error("Readall failed", "error", err)
		return
	}
	if err = os.WriteFile(rFilePath, context, 0644); err != nil {
		slog.Error("write to file failed",
			"owner", item.Owner,
			"repo", item.Repo,
			"version", item.Version,
			"name", item.AssetName,
			"error", err)
	} else {
		slog.Info("write to file succeed", "file_path", rFilePath)
	}

}

// ExcludeTxt 判断字符串A中是否含有字符串B
func ExcludeTxt(txt string) bool {
	excludeTxt := config.ExcludeTxt
	for _, v := range excludeTxt {
		if strings.Contains(txt, v) {
			return true
		}
	}
	return false
}

func Mkdir(path string) {
	if !IsDirExists(path) {
		CreateDir(path)
	}
}

// IsDirExists 判断文件目录否存在
func IsDirExists(path string) bool {
	if path == "" {
		return false
	}
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}

}

// CreateDir 创建目录
func CreateDir(dirs ...string) (err error) {
	for _, v := range dirs {
		exist := IsDirExists(v)

		if !exist {

			if err := os.MkdirAll(v, os.ModePerm); err != nil {
				return err
			}
		}
	}
	return err
}

// Download 普通下载
func Download(src_url, dst_path string) (err error) {

	request, err := http.NewRequest("GET", url.PathEscape(src_url), nil)
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
		slog.Error("request web failed", "error", err)
		return
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	os.WriteFile(dst_path, data, 0644)
	return
}

// // MonitorDownload 下载文件到 workspace 的 monitor 目录中
// func MonitorDownload(owner, repo string) {
// 	down := Githubclient.NewGitHubRelease(owner, repo, config.MirrorRoot+"/monitor")
// 	down.Download(Githubclient)
// }

// // OtherDownload 其他下载
// func OtherDownload(owner, repo string) {
// 	down := Githubclient.NewGitHubRelease(owner, repo, config.MirrorRoot)
// 	down.Download(Githubclient)

// }
