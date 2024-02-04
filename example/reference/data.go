package main

import (
	_ "embed"
	"encoding/csv"
	"io"
	"sort"
	"strconv"
	"strings"
)

//go:embed articles.csv
var csvArticles string

//go:embed authors.csv
var csvAuthors string

func readAuthors() (map[int]Author, error) {
	r := csv.NewReader(strings.NewReader(csvAuthors))
	results := make(map[int]Author)
	// skip the title row
	_, err := r.Read()
	if err == io.EOF {
		return results, nil
	}
	if err != nil {
		return nil, err
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		id, err := strconv.ParseInt(record[0], 10, 32)
		if err != nil {
			return nil, err
		}
		results[int(id)] = Author{
			ID:        int(id),
			FirstName: record[1],
			LastName:  record[2],
		}
	}

	return results, nil
}

func readArticles() (map[int]Article, error) {
	r := csv.NewReader(strings.NewReader(csvArticles))
	results := make(map[int]Article)
	// skip the title row
	_, err := r.Read()
	if err == io.EOF {
		return results, nil
	}
	if err != nil {
		return nil, err
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		id, err := strconv.ParseInt(record[0], 10, 32)
		if err != nil {
			return nil, err
		}
		authorID, err := strconv.ParseInt(record[2], 10, 32)
		if err != nil {
			return nil, err
		}
		results[int(id)] = Article{
			ID:       int(id),
			Title:    record[1],
			AuthorID: int(authorID),
		}
	}

	return results, nil
}

func getMapValues[K comparable, V any](input map[K]V) []V {
	results := make([]V, 0, len(input))
	for _, v := range input {
		results = append(results, v)
	}
	return results
}

func getMapKeys[K comparable, V any](input map[K]V) []K {
	results := make([]K, 0, len(input))
	for k := range input {
		results = append(results, k)
	}
	return results
}

func sortArticles(input map[int]Article, key string, descending bool) []Article {
	ids := getMapKeys(input)
	sort.Ints(ids)
	results := make([]Article, 0, len(ids))
	if !descending {
		for _, id := range ids {
			results = append(results, input[id])
		}
		return results
	}

	for i := range ids {
		id := ids[len(ids)-1-i]
		results = append(results, input[id])
	}
	return results
}
