package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	"github.com/google/go-github/v24/github"
)

func main() {
	p := os.Getenv("GITHUB_EVENT_PATH")
	if len(p) == 0 {
		log.Fatalf("GITHUB_EVENT_PATH must be set")
	}

	f, err := os.Open(p)
	if err != nil {
		log.Fatalf("error opening event file %s: %v\n", p, err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("error reading file %s: %v\n", p, err)
	}

	e := github.ReleaseEvent{}
	json.Unmarshal(b, &e)

	var verified []string

	for _, v := range e.Release.Assets {
		verified = append(verified, *v.BrowserDownloadURL)
	}
	verified = append(verified, *e.Release.ZipballURL)
	verified = append(verified, *e.Release.TarballURL)

	archives := archiveURLs(*e.Repo.Owner.Login, *e.Repo.Name, *e.Release.TagName)
	for _, v := range archives {
		verified = append(verified, v)
	}

	fmt.Printf("::set-output name=verified::%v\n", verified)
}

// archiveURLs generates source archive URLs for a GitHub repo tag
// e.g. https://github.com/philips/releases-test/archive/v1.0.zip and
// https://github.com/philips/releases-test/archive/v1.0.tar.gz
func archiveURLs(owner, repo, tag string) (urls []string) {
	u := url.URL{
		Scheme: "https",
		Host:   "github.com",
		Path:   fmt.Sprintf("/%s/%s/archive/%s", owner, repo, tag),
	}
	urls = append(urls, u.String()+".tar.gz")
	urls = append(urls, u.String()+".zip")

	return
}
