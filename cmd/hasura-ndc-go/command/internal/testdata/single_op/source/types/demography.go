package types

type Demography struct {
	ID string
}

type DemographyResult struct {
	Demography []Demography `json:"demography"`
	ErrorCode  string       `json:"errorCode"`
}
