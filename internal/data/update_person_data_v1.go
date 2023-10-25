package data

import (
	"context"
	"fmt"

	"github.com/barpav/demography/internal/rest/models"
)

type queryUpdatePersonDataV1 struct{}

func (q queryUpdatePersonDataV1) text() string {
	return `
	UPDATE people SET
		surname = $1,
		person_name = $2,
		patronymic = NULLIF($3, ''),
		age = NULLIF($4, 0),
		gender = NULLIF($5, '')::gender,
		country = NULLIF($6, '')
	WHERE id = $7;
	`
}

type ErrPersonDataNotFound struct{}

func (s *Storage) UpdatePersonDataV1(ctx context.Context, id int64, data *models.EditedPersonDataV1) error {
	result, err := s.queries[queryUpdatePersonDataV1{}].ExecContext(ctx,
		data.Surname, data.Name, data.Patronymic, data.Age, data.Gender, data.Country, id)

	var updated int64
	if err == nil {
		updated, err = result.RowsAffected()
	}

	if err != nil {
		return fmt.Errorf("failed to update person data (v1): %w", err)
	}

	if updated == 0 {
		return ErrPersonDataNotFound{}
	}

	return nil
}

func (e ErrPersonDataNotFound) Error() string {
	return "person data not found"
}

func (e ErrPersonDataNotFound) ImplementsPersonDataNotFoundError() {
}
