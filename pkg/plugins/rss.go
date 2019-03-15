package plugins

import (
	"encoding/json"
	"fmt"
	"github.com/advancedlogic/goms/pkg/interfaces"
	"github.com/advancedlogic/goms/pkg/tools"
	"github.com/mmcdole/gofeed"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

type RSS struct {
	interfaces.Service
}

func NewRSS() *RSS {
	return &RSS{}
}

func (r *RSS) Reload(descriptor Descriptor) error {
	files, err := ioutil.ReadDir(descriptor.Folder)
	if err != nil {
		return err
	}
	urls := make([]string, 0)
	for _, file := range files {
		name := file.Name()
		if strings.HasSuffix(name, ".rss") {
			blines, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", descriptor.Folder, file.Name()))
			if err != nil {
				r.Error(err.Error())
				continue
			}
			lines := strings.Split(string(blines), "\n")
			for _, line := range lines {
				urls = append(urls, strings.TrimSuffix(line, "\r\n"))
			}
		}
	}

	descriptor.Urls = urls
	return nil
}

type Descriptor struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Folder  string   `json:"folder, omitempty"`
	Urls    []string `json:"urls, omitempty"`
	Timeout int      `json:"timeout"`
}

func (r *RSS) Init(service interfaces.Service) error {
	r.Service = service

	folder := r.Config("rss.folder", "").(string)

	if folder != "" {
		loopString := r.Config("rss.loop", "0s").(string)
		loop := tools.String2Milliseconds(loopString)
		load := func() (Descriptor, error) {
			descriptor := Descriptor{}
			files, err := ioutil.ReadDir(folder)
			if err != nil {
				return descriptor, err
			}
			for _, file := range files {
				b, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", folder, file.Name()))
				if err != nil {
					r.Error(err.Error())
					continue
				}
				urls := strings.Split(string(b), "\n")
				for _, url := range urls {
					url = strings.TrimSpace(url)
					descriptor.Urls = append(descriptor.Urls, url)
				}
			}

			return descriptor, nil
		}
		routine := func(descriptor Descriptor) {
			feeds := r.Download(descriptor)
			for _, feed := range feeds {
				id := tools.SHA1(feed)
				r.Infof("[%s] %s", id, feed)
				if err := r.Store().Create(id, feed); err != nil {
					r.Error(err.Error())
				}
			}
		}
		final := func() {
			descriptor, err := load()
			if err != nil {
				r.Fatal(err.Error())
			}
			routine(descriptor)
		}
		if loop > 0 {
			go func() {
				for r.Running() {
					final()
					r.Infof("Round is over. Waiting now %s", loopString)
					time.Sleep(time.Duration(loop) * time.Millisecond)
				}
				r.Warnf("RSS loop is over")
			}()
		} else {
			go final()
		}
	}
	return nil
}

func (r *RSS) Close() error { return nil }

func (r *RSS) Process(data interface{}) (interface{}, error) {
	descriptor := data.(Descriptor)
	err := r.Reload(descriptor)
	if err != nil {
		return nil, err
	}
	rss := r.Download(descriptor)
	return rss, nil
}

func Shuffle(vals []string) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(vals) > 0 {
		n := len(vals)
		randIndex := r.Intn(n)
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1]
		vals = vals[:n-1]
	}
}

func (r *RSS) Download(descriptor Descriptor) []string {
	fp := gofeed.NewParser()
	Shuffle(descriptor.Urls)
	feeds := make([]string, 0)
	for _, url := range descriptor.Urls {
		feed, err := fp.ParseURL(url)
		if err != nil {
			r.Error(err.Error())
			continue
		}

		bstr, err := json.Marshal(feed)
		if err != nil {
			r.Error(err.Error())
			continue
		}
		feeds = append(feeds, string(bstr))
	}

	return feeds
}
