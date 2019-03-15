package plugin

import (
	"github.com/advancedlogic/GoOse"
	"github.com/advancedlogic/goms/pkg/interfaces"
	"github.com/advancedlogic/goms/pkg/tools"
	"time"
)

type SpiderSource struct {
	Url string `json:"url"`
}

type SpiderDescriptor struct {
	ID          string `json:"id"`
	Url         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
	Content     string `json:"content"`
	Image       string `json:"image"`
	Timestamp   int64  `json:"timestamp"`
}

type Spider struct {
	interfaces.Service
}

func (s *Spider) Init(service interfaces.Service) error {
	s.Service = service
	return nil
}

func (s *Spider) Close() error {
	return nil
}

func (s *Spider) Process(data interface{}) (interface{}, error) {
	var source SpiderSource

	switch data.(type) {
	case string:
		source.Url = data.(string)
	case SpiderSource:
		source = data.(SpiderSource)
	}

	g := goose.New()

	article, err := g.ExtractFromURL(source.Url)
	if err != nil {
		return nil, err
	}

	return &SpiderDescriptor{
		ID:          tools.UUID(),
		Title:       article.Title,
		Description: article.MetaDescription,
		Keywords:    article.MetaKeywords,
		Content:     article.CleanedText,
		Url:         article.FinalURL,
		Image:       article.TopImage,
		Timestamp:   time.Now().Unix(),
	}, nil
}
