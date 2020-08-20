package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/cavaliercoder/grab"
	"github.com/google/go-github/v24/github"
	"go.transparencylog.com/tl/config"
	"go.transparencylog.com/tl/sumdb"
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

	var assets []string

	for _, v := range e.Release.Assets {
		assets = append(assets, *v.BrowserDownloadURL)
	}
	assets = append(assets, *e.Release.ZipballURL)
	assets = append(assets, *e.Release.TarballURL)

	archives := archiveURLs(*e.Repo.Owner.Login, *e.Repo.Name, *e.Release.TagName)
	for _, v := range archives {
		assets = append(assets, v)
	}

	for _, v := range assets {
		get(v)
	}

	fmt.Printf("::set-output name=verified::%v\n", assets)
}

func get(durl string) {
	u, err := url.Parse(durl)
	if err != nil {
		panic(err)
	}
	key := u.Host + u.Path

	cache := config.ClientCache()

	// create download request
	req, err := grab.NewRequest("", durl)
	if err != nil {
		fmt.Printf("failed to create grab request: %v\n", err)
		os.Exit(1)
	}
	req.NoCreateDirectories = true
	req.SkipExisting = true

	req.AfterCopy = func(resp *grab.Response) (err error) {
		var f *os.File
		f, err = os.Open(resp.Filename)
		if err != nil {
			return
		}
		defer func() {
			f.Close()
		}()

		h := sha256.New()
		_, err = io.Copy(h, f)
		if err != nil {
			return err
		}

		fileSum := h.Sum(nil)

		// Download the tlog entry for the URL
		want := "h1:" + base64.StdEncoding.EncodeToString(fileSum)
		client := sumdb.NewClient(cache)
		_, data, err := client.LookupOpts(key, sumdb.LookupOpts{Digest: want})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("fetched note: %s/lookup/%s\n", config.ServerURL, key)

		for _, line := range strings.Split(string(data), "\n") {
			if line == want {
				break
			}
			if strings.HasPrefix(line, "h1:") {
				log.Fatalf("file digest %x != log digest %x", fileSum, line)
			}
		}

		fmt.Printf("validated file sha256sum: %x\n", fileSum)

		req.SetChecksum(sha256.New(), fileSum, true)

		return
	}

	// download and validate file
	resp := grab.DefaultClient.Do(req)
	if err := resp.Err(); err != nil {
		fmt.Printf("Failed to grab: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Download validated and saved to", resp.Filename)
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
