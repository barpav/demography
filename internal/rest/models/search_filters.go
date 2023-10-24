package models

type SearchFilters struct {
	Surname    string
	Name       string
	Patronymic string
	Age        int
	Gender     string
	Country    string
	After      int64
	Limit      int
}
