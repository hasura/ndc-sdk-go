package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hasura/ndc-sdk-go/schema"
)

var baseFixtureURL = fmt.Sprintf(
	"https://raw.githubusercontent.com/hasura/ndc-spec/refs/tags/v%s/ndc-reference",
	schema.NDCVersion,
)

type Article struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	PublishedDate string `json:"published_date"`
	AuthorID      int    `json:"author_id"`
}

type Author struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type InstitutionLocation struct {
	CountryID int      `json:"country_id"`
	City      string   `json:"city"`
	Country   string   `json:"country"`
	Campuses  []string `json:"campuses"`
}

type InstitutionStaff struct {
	FirstName     string   `json:"first_name"`
	LastName      string   `json:"last_name"`
	Specialities  []string `json:"specialities"`
	BornCountryID int      `json:"born_country_id"`
}

type Institution struct {
	ID          int                 `json:"id"`
	Name        string              `json:"name"`
	Location    InstitutionLocation `json:"location"`
	Staff       []InstitutionStaff  `json:"staff"`
	Departments []string            `json:"departments"`
}

type Country struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	AreaKm2 int    `json:"area_km2"`
	Cities  []City `json:"cities"`
}

type City struct {
	Name string `json:"city"`
}

func fetchAuthors() ([]Author, error) {
	return fetchRemoteFixtures[Author]("authors.jsonl")
}

func fetchArticles() ([]Article, error) {
	return fetchRemoteFixtures[Article]("articles.jsonl")
}

func fetchInstitutions() ([]Institution, error) {
	return fetchRemoteFixtures[Institution]("institutions.jsonl")
}

func fetchCountries() ([]Country, error) {
	return fetchRemoteFixtures[Country]("countries.jsonl")
}

func fetchRemoteFixtures[T any](name string) ([]T, error) {
	res, err := http.Get(baseFixtureURL + "/" + name) //nolint:noctx
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	var results []T

	decoder := json.NewDecoder(res.Body)
	for decoder.More() {
		var r T

		err := decoder.Decode(&r)
		if err != nil {
			return nil, err
		}

		results = append(results, r)
	}

	return results, nil
}
