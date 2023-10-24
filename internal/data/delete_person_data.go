package data

import (
	"context"
	"fmt"
)

type queryDeletePersonData struct{}

func (q queryDeletePersonData) text() string {
	return `
	DELETE FROM people WHERE id = $1;
	`
}

func (s *Storage) DeletePersonData(ctx context.Context, id int64) error {
	result, err := s.queries[queryDeletePersonData{}].ExecContext(ctx, id)

	var deleted int64
	if err == nil {
		deleted, err = result.RowsAffected()
	}

	if err != nil {
		return fmt.Errorf("failed to delete person data: %w", err)
	}

	if deleted == 0 {
		return ErrPersonDataNotFound{}
	}

	return nil
}
