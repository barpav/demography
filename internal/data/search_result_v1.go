package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/barpav/demography/internal/rest/models"
	"github.com/doug-martin/goqu/v9"
)

func (s *Storage) SearchResultV1(ctx context.Context, filters *models.SearchFilters) (result *models.SearchResultV1, err error) {
	builder := goqu.Select(
		"id",
		"surname",
		"person_name",
		goqu.COALESCE(goqu.C("patronymic"), ""),
		goqu.COALESCE(goqu.C("age"), 0),
		goqu.COALESCE(goqu.C("gender").Cast("varchar"), ""),
		goqu.COALESCE(goqu.C("country"), ""),
	).From("people")

	if filters.After != 0 {
		builder = builder.Where(goqu.C("id").Gt(filters.After))
	}

	if filters.Surname != "" {
		builder = builder.Where(goqu.C("surname").Eq(filters.Surname))
	}

	if filters.Name != "" {
		builder = builder.Where(goqu.C("person_name").Eq(filters.Name))
	}

	if filters.Patronymic != "" {
		builder = builder.Where(goqu.C("patronymic").Eq(filters.Patronymic))
	}

	if filters.Age != 0 {
		builder = builder.Where(goqu.C("age").Eq(filters.Age))
	}

	if filters.Gender != "" {
		builder = builder.Where(goqu.C("gender").Eq(filters.Gender))
	}

	if filters.Country != "" {
		builder = builder.Where(goqu.C("country").Eq(filters.Country))
	}

	builder = builder.Order(goqu.C("id").Asc())
	builder = builder.Limit(uint(filters.Limit))

	var query string
	query, _, err = builder.ToSQL()

	if err != nil {
		return nil, fmt.Errorf("failed to build sql query text for search result (v1): %w", err)
	}

	var rows *sql.Rows
	rows, err = s.db.QueryContext(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("failed to execute sql query for search result (v1): %w", err)
	}

	defer rows.Close()

	result = &models.SearchResultV1{Data: make([]*models.EnrichedPersonDataV1, 0, filters.Limit)}

	for rows.Next() {
		info := &models.EnrichedPersonDataV1{}
		err = rows.Scan(
			&info.Id,
			&info.Surname,
			&info.Name,
			&info.Patronymic,
			&info.Age,
			&info.Gender,
			&info.Country,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to process sql query result for search result (v1): %w", err)
		}

		result.Data = append(result.Data, info)
	}

	err = rows.Err()

	if err != nil {
		return nil, fmt.Errorf("failed to process sql query results for search result (v1): %w", err)
	}

	result.Total = len(result.Data)

	return result, err
}
