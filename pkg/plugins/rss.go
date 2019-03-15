package plugins

import (
	"encoding/json"
	"fmt"
	"github.com/advancedlogic/goms/pkg/interfaces"
	"github.com/advancedlogic/goms/pkg/tools"
	"github.com/mmcdole/gofeed"
	"io/ioutil"
	"strings"
	"time"
)

type RSSSource struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Folder  string   `json:"folder, omitempty"`
	Urls    []string `json:"urls, omitempty"`
	Timeout int      `json:"timeout"`
}

type RSS struct {
	interfaces.Service
}

func NewRSS() *RSS {
	return &RSS{}
}

func (r *RSS) reload(source RSSSource) error {
	files, err := ioutil.ReadDir(source.Folder)
	if err != nil {
		return err
	}
	urls := make([]string, 0)
	for _, file := range files {
		name := file.Name()
		if strings.HasSuffix(name, ".rss") {
			blines, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", source.Folder, file.Name()))
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

	source.Urls = urls
	return nil
}

func (r *RSS) Init(service interfaces.Service) error {
	r.Service = service

	folder := r.Config("rss.folder", "").(string)

	if folder != "" {
		loopString := r.Config("rss.loop", "0s").(string)
		loop := tools.String2Milliseconds(loopString)
		load := func() (RSSSource, error) {
			source := RSSSource{}
			source.ID = tools.UUID()
			source.Folder = folder

			files, err := ioutil.ReadDir(folder)
			if err != nil {
				return source, err
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
					source.Urls = append(source.Urls, url)
				}
			}

			return source, nil
		}
		routine := func(source RSSSource) {
			feeds := r.download(source)
			for _, feed := range feeds {
				id := tools.SHA1(feed)
				r.Infof("[%s] %s", id, feed)
				if err := r.Store().Create(id, feed); err != nil {
					r.Error(err.Error())
				}
			}
		}
		final := func() {
			source, err := load()
			if err != nil {
				r.Fatal(err.Error())
			}
			routine(source)
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
	var source RSSSource
	switch data.(type) {
	case RSSSource:
		source = data.(RSSSource)
	case string:

	}
	err := r.reload(source)
	if err != nil {
		return nil, err
	}
	rss := r.download(source)
	return rss, nil
}

func (r *RSS) download(source RSSSource) []string {
	fp := gofeed.NewParser()
	tools.Shuffle(source.Urls)
	feeds := make([]string, 0)
	for _, url := range source.Urls {
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
