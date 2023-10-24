package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/barpav/demography/internal/rest/models"
)

type queryGetEnrichedPersonDataV1 struct{}

func (q queryGetEnrichedPersonDataV1) text() string {
	return `
	SELECT
		surname,
		person_name,
		COALESCE(patronymic, ''),
		COALESCE(age, 0),
		COALESCE(gender::varchar, ''),
		COALESCE(country, '')
	FROM people
	WHERE id = $1;
	`
}

// Returns nil, nil if data is not found.
func (s *Storage) EnrichedPersonDataV1(ctx context.Context, id int64) (*models.EnrichedPersonDataV1, error) {
	row := s.queries[queryGetEnrichedPersonDataV1{}].QueryRowContext(ctx, id)
	err := row.Err()

	if err != nil {
		return nil, fmt.Errorf("failed to execute sql statement (enrichedPersonData.v1): %w", err)
	}

	data := &models.EnrichedPersonDataV1{Id: id}
	err = row.Scan(
		&data.Surname,
		&data.Name,
		&data.Patronymic,
		&data.Age,
		&data.Gender,
		&data.Country,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to scan sql result (enrichedPersonData.v1): %w", err)
	}

	return data, nil
}
