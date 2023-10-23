package data

import (
	"context"

	"github.com/barpav/demography/internal/rest/models"
)

type queryCreateNewPersonDataV1 struct{}

func (q queryCreateNewPersonDataV1) text() string {
	return `
	INSERT INTO people (surname, person_name, patronymic, age, gender, country)
	VALUES ($1, $2, NULLIF($3, ''), NULLIF($4, 0), NULLIF($5, '')::gender, NULLIF($6, ''))
	RETURNING id;
	`
}

func (s *Storage) CreateNewPersonDataV1(ctx context.Context, data *models.EnrichedPersonDataV1) error {
	row := s.queries[queryCreateNewPersonDataV1{}].QueryRowContext(ctx,
		data.Surname, data.Name, data.Patronymic, data.Age, data.Gender, data.Country)
	return row.Scan(&data.Id)
}
