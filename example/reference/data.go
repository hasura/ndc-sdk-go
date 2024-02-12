package main

import (
	_ "embed"
	"encoding/csv"
	"encoding/json"
	"io"
	"strconv"
	"strings"
)

//go:embed data/articles.csv
var csvArticles string

//go:embed data/authors.csv
var csvAuthors string

//go:embed data/institutions.json
var jsonInstitutions []byte

func readAuthors() ([]Author, error) {
	r := csv.NewReader(strings.NewReader(csvAuthors))
	results := make([]Author, 0)
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
		results = append(results, Author{
			ID:        int(id),
			FirstName: record[1],
			LastName:  record[2],
		})
	}

	return results, nil
}

func readArticles() ([]Article, error) {
	r := csv.NewReader(strings.NewReader(csvArticles))
	results := make([]Article, 0)
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
		results = append(results, Article{
			ID:       int(id),
			Title:    record[1],
			AuthorID: int(authorID),
		})
	}

	return results, nil
}

func readInstitutions() ([]Institution, error) {
	var institutions []Institution
	if err := json.Unmarshal(jsonInstitutions, &institutions); err != nil {
		return nil, err
	}

	return institutions, nil
}
