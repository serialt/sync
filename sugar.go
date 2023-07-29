package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v45/github"
	"github.com/serialt/lancet/cryptor"
	"golang.org/x/exp/slog"
	"golang.org/x/oauth2"
)

func Run() {
	Githubclient = &GithubClient{Token: config.GithubToken}

	// LogSugar.Info("info log")
	// LogSugar.Info(ConfigPath)

	// LogSugar.Info(config.LogFile)
	// service.GetLastestRelease("fatedier", "frp")
	// service.DownloadReleaseAsset("fatedier", "frp", 56250083)

	// 减少段时间内请求github的次数

	for _, v := range config.GithubRelease {
		time.Sleep(time.Second * 3)
		myMonitorRepo := strings.Split(v, "/")
		OtherDownload(myMonitorRepo[0], myMonitorRepo[1])

	}

	// monitor
	for _, v := range config.Monitor {
		time.Sleep(time.Second * 3)
		myMonitorRepo := strings.Split(v, "/")
		MonitorDownload(myMonitorRepo[0], myMonitorRepo[1])
	}
	// down := service.NewGitHubRelease(myM, "frp", "/tmp")
	// down.Download()

}

// Decrypt 解密
func Decrypt(data, key string) string {

	_data, _ := base64.StdEncoding.DecodeString(data)
	_key := []byte(key)

	tmpByte := cryptor.AesCbcDecrypt(_data, _key)
	return string(tmpByte)
}

// 加密
func Encrypt(data, key string) string {
	_data := []byte(data)
	_key := []byte(key)

	tmpByte := cryptor.AesCbcEncrypt(_data, _key)

	// fmt.Println(Decrypt(base64.StdEncoding.EncodeToString(tmpByte), key))
	return base64.StdEncoding.EncodeToString(tmpByte)
}

// 解密配置文件
func (c *Config) DecryptConfig() {
	if c.Encrypt {
		c.GithubToken = Decrypt(c.GithubToken, AesKey)
		slog.Debug(c.GithubToken)
	}
}

// 功能函数区域

func NewClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client
}

// GetLastestReleaseAsset 获取最新的 ReleaseAsset
func (c *GithubClient) GetLastestReleaseAsset() {
	// fatedier / frp
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	reader, _, err := client.Repositories.DownloadReleaseAsset(ctx, "fatedier", "frp", 1, http.DefaultClient)

	if err != nil {
		slog.Error("get latest release failed", "error", err)
	}
	content, _ := ioutil.ReadAll(reader)
	ioutil.WriteFile("/tmp/frp.tar.gz", content, 0644)
}

// func GetLastestReleaseID() {
// 	ctx := context.Background()
// 	client := github.NewClient(nil)

// }

// GetLastestRelease 获取最新的稳定release
func (c *GithubClient) GetLastestRelease(owner, repo string) (release *github.RepositoryRelease) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	release, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		slog.Info("cat not get the lastest release",
			"owner", owner,
			"repo", repo,
			"error", err)
	}
	// slog.Info("release msg", "release", release)
	return
}

// ListRelease 获取最近 lastNum 个数的release,可能包括beta和pre-release
func (c *GithubClient) ListRelease(owner, repo string, lastNum int) (releaseList []*github.RepositoryRelease) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	opt := &github.ListOptions{Page: 1, PerPage: lastNum}
	release, _, err := client.Repositories.ListReleases(ctx, owner, repo, opt)
	if err != nil {
		slog.Error("cat not get the lastest release",
			"owner", owner,
			"repo", repo,
			"error", err)
	}
	slog.Debug("release msg", "release", release)
	return
}

// DownloadReleaseAsset 根据assetID 下载release到指定的目录
func (c *GithubClient) DownloadReleaseAsset(owner, repo string, assetID int, filepath string) {
	// 如果目录里文件存在则不操作
	if IsDirExists(filepath) {
		return
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	reader, _, err := client.Repositories.DownloadReleaseAsset(ctx, owner, repo, int64(assetID), http.DefaultClient)
	if err != nil {
		slog.Error("Download release asset failed",
			"owner", owner,
			"repo", repo,
			"error", err)
	}
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		slog.Error("Readall failed", "error", err)

	}
	if err = ioutil.WriteFile(filepath, content, 0644); err != nil {
		slog.Error("write to file failed", "error", err)
	} else {
		slog.Info("write to file succeed", "file_path", filepath)
	}

}

// NewGitHubRelease 获取单个版本的release信息
func (c *GithubClient) NewGitHubRelease(owner, repo, path string) *GithubRelease {
	GR := &GithubRelease{
		Owner: owner,
		Repo:  repo,
	}

	myRelease := c.GetLastestRelease(owner, repo)
	GR.Version = *myRelease.TagName
	slog.Info("get a release",
		"owner", owner,
		"repo", repo,
		"version", GR.Version,
	)

	for _, v := range myRelease.Assets {
		slog.Info("get asset name", "name", *v.Name)
		if ExcludeTxt(*v.Name) {
			continue
		}
		// if ExcludeTxt(*v.BrowserDownloadURL) {
		// 	break
		// }
		GR.AssetID = append(GR.AssetID, int(*v.ID))
		GR.AssetName = append(GR.AssetName, *v.Name)
		// GR.BrowserDownloadUrl = append(GR.BrowserDownloadUrl, *v.BrowserDownloadURL)
		slog.Debug("get release msg",
			"version", *myRelease.TagName,
			"asset_id", int(*v.ID),
			"asset_name", *v.Name,
			"download_url", *v.BrowserDownloadURL,
		)

	}
	GR.Path = fmt.Sprintf("%s/%s/%s", path, GR.Repo, GR.Version)
	slog.Info("release path", "path", GR.Path)

	Mkdir(GR.Path)
	slog.Info("mkdir path", "path", GR.Path)

	return GR
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

// Download 下载
func (g *GithubRelease) Download(client *GithubClient) {
	for k, v := range g.AssetID {
		filename := fmt.Sprintf("%s/%s", g.Path, g.AssetName[k])
		client.DownloadReleaseAsset(g.Owner, g.Repo, v, filename)
	}

}

// Download 普通下载
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
		slog.Error("request web failed", "error", err)
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

// MonitorDownload 下载文件到 workspace 的 monitor 目录中
func MonitorDownload(owner, repo string) {
	down := Githubclient.NewGitHubRelease(owner, repo, config.MirrorRoot+"/monitor")
	down.Download(Githubclient)
}

// OtherDownload 其他下载
func OtherDownload(owner, repo string) {
	down := Githubclient.NewGitHubRelease(owner, repo, config.MirrorRoot)
	down.Download(Githubclient)

}
