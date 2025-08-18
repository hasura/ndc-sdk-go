package main

import (
	"github.com/hasura/ndc-sdk-go/v2/connector"
)

type Configuration struct{}

type State struct {
	Authors      []Author
	Articles     []Article
	Institutions []Institution
	Countries    []Country
	Telemetry    *connector.TelemetryState
}

func (s *State) GetLatestArticle() *Article {
	if len(s.Articles) == 0 {
		return nil
	}

	var latestArticle Article
	for _, article := range s.Articles {
		if latestArticle.ID < article.ID {
			latestArticle = article
		}
	}

	return &latestArticle
}

type Chunk struct {
	Dimensions []any
	Rows       []map[string]any
}
