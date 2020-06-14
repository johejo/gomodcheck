package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-github/v32/github"
	"github.com/hashicorp/go-version"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
	"golang.org/x/oauth2"
)

var (
	httpClient = &http.Client{}

	modFile string
)

func init() {
	flag.StringVar(&modFile, "modfile", "go.mod", "go.mod path")
}

func main() {
	flag.Parse()
	if err := _main(); err != nil {
		log.Fatal(err)
	}
}

func _main() error {
	ctx := context.Background()
	data, err := ioutil.ReadFile(modFile)
	if err != nil {
		return err
	}
	f, err := modfile.Parse(modFile, data, nil)
	if err != nil {
		return err
	}
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		httpClient = oauth2.NewClient(ctx, ts)
	}
	var wg sync.WaitGroup
	for _, r := range f.Require {
		r := r
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := checkMod(ctx, r.Mod); err != nil {
				log.Println(err)
			}
		}()
	}
	wg.Wait()
	return nil
}

func checkMod(ctx context.Context, m module.Version) error {
	p := m.Path
	if strings.HasPrefix(p, "gopkg.in") {
		return fmt.Errorf("gopkg.in is not supported, skip %s", p)
	}
	u := fmt.Sprintf("https://%s", p)
	source := u
	if !strings.HasPrefix(p, "github") {
		resp, err := httpClient.Get(u)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return err
		}
		doc.Find("head meta").Each(func(i int, selection *goquery.Selection) {
			name, ok := selection.Attr("name")
			if !ok {
				return
			}
			if name != "go-source" {
				return
			}
			c, ok := selection.Attr("content")
			if !ok {
				return
			}
			s := strings.Split(c, " ")
			source = s[1]
		})
	}
	_url, err := url.Parse(source)
	if err != nil {
		return err
	}
	h := _url.Hostname()
	if h != "github.com" {
		return fmt.Errorf("source from %s is not supported, skip %s", h, p)
	}
	_repo := strings.Split(strings.TrimPrefix(_url.Path, "/"), "/")
	owner := _repo[0]
	repo := _repo[1]

	client := github.NewClient(httpClient)
	tags, _, err := client.Repositories.ListTags(ctx, owner, repo, &github.ListOptions{Page: 1, PerPage: 1})
	if err != nil {
		return err
	}
	if len(tags) == 0 {
		return fmt.Errorf("%s/%s has no tags", owner, repo)
	}

	current, err := version.NewSemver(m.Version)
	if err != nil {
		return err
	}
	latest, err := version.NewSemver(tags[0].GetName())
	if err != nil {
		return err
	}
	if current.LessThan(latest) {
		log.Printf("%s is behind, latest=v%s", p, latest.String())
	} else {
		log.Printf("%s v%s is latest", p, current.String())
	}
	return nil
}
