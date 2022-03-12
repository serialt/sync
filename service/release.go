package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-github/v42/github"
	"github.com/serialt/sync/config"
	"github.com/serialt/sync/pkg"
)

type GithubRelease struct {
	Owner              string
	Repo               string
	Version            string
	AssetName          []string
	AssetID            []int
	BrowserDownloadUrl []string
	Path               string
}

func GetLastestReleaseAsset() {
	// fatedier / frp
	ctx := context.Background()

	client := github.NewClient(nil)

	reader, _, err := client.Repositories.DownloadReleaseAsset(ctx, "fatedier", "frp", 1, http.DefaultClient)

	if err != nil {
		pkg.Sugar.Infof("get latest release failed: %v", err)
	}
	content, _ := ioutil.ReadAll(reader)
	ioutil.WriteFile("/tmp/frp.tar.gz", content, 0644)
}

// func GetLastestReleaseID() {
// 	ctx := context.Background()
// 	client := github.NewClient(nil)

// }

// 获取最新的稳定release
func GetLastestRelease(owner, repo string) (release *github.RepositoryRelease) {
	ctx := context.Background()
	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		pkg.Sugar.Infof("cat not get the lastest release, owner: %v,repo: %v, err: %v", owner, repo, err)
	}
	pkg.Sugar.Debugf("release msg: %v", release)
	return
}

// 获取最近 lastNum 个数的release,可能包括beta和pre-release
func ListRelease(owner, repo string, lastNum int) (releaseList []*github.RepositoryRelease) {
	ctx := context.Background()
	client := github.NewClient(nil)
	opt := &github.ListOptions{Page: 1, PerPage: lastNum}
	release, _, err := client.Repositories.ListReleases(ctx, owner, repo, opt)
	if err != nil {
		pkg.Sugar.Infof("cat not get the lastest release, owner: %v,repo: %v, err: %v", owner, repo, err)
	}
	pkg.Sugar.Debugf("release msg: %v", release)
	return
}

// 根据assetID 下载release到指定的目录
func DownloadReleaseAsset(owner, repo string, assetID int, filepath string) {
	// 如果目录里文件存在则不操作
	if IsDirExists(filepath) {
		return
	}
	ctx := context.Background()
	client := github.NewClient(nil)
	reader, _, err := client.Repositories.DownloadReleaseAsset(ctx, owner, repo, int64(assetID), http.DefaultClient)
	if err != nil {
		pkg.Sugar.Infof("Download release asset failed, owner: %v,repo: %v,asset_id: %v,err: %v", owner, repo, int64(assetID), err)
	}
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		pkg.Sugar.Infof("Readall failed, err: ", err)
	}
	ioutil.WriteFile(filepath, content, 0644)
	pkg.Sugar.Infof("Release download file succeed, file path: %v", filepath)
}

// 获取单个版本的release信息
func NewGitHubRelease(owner, repo, path string) *GithubRelease {
	GR := &GithubRelease{
		Owner: owner,
		Repo:  repo,
	}

	myRelease := GetLastestRelease(owner, repo)
	GR.Version = *myRelease.TagName
	pkg.Sugar.Infof("get a release, owner: %s ,repo: %s ,version: %s", owner, repo, GR.Version)
	for _, v := range myRelease.Assets {
		if ExcludeTxt(*v.Name) {
			continue
		}
		// if ExcludeTxt(*v.BrowserDownloadURL) {
		// 	break
		// }
		GR.AssetID = append(GR.AssetID, int(*v.ID))
		GR.AssetName = append(GR.AssetName, *v.Name)
		// GR.BrowserDownloadUrl = append(GR.BrowserDownloadUrl, *v.BrowserDownloadURL)
		pkg.Sugar.Infof("get release msg, version: %s ,asset id: %d ,asset_name: %s ,download_url: %s ", *myRelease.TagName, int(*v.ID), *v.Name, *v.BrowserDownloadURL)
	}
	GR.Path = fmt.Sprintf("%s/%s/%s", path, GR.Repo, GR.Version)
	pkg.Sugar.Infof("release path: %s", GR.Path)

	Mkdir(GR.Path)
	pkg.Sugar.Infof("mkdir path: %s", GR.Path)
	return GR
}

// 判断字符串A中是否含有字符串B
func ExcludeTxt(txt string) bool {
	excludeTxt := config.Config.ExcludeTxt
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

// 判断文件目录否存在
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

func (g *GithubRelease) Download() {
	for k, v := range g.AssetID {
		filename := fmt.Sprintf("%s/%s", g.Path, g.AssetName[k])
		DownloadReleaseAsset(g.Owner, g.Repo, v, filename)
	}

}
