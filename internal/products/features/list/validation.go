package list

import (
	"strings"

	apperrors "github.com/quinn9x/go-vertical-slice/internal/shared/errors"
)

func validateQuery(query ListerProductsQuery) error {
	if err := validateSortBy(query.SortBy); err != nil {
		return err
	}

	return validateSortOrder(query.SortOrder)
}

func validateSortBy(value string) error {
	switch value {
	case SortByCreatedAt, SortByName, SortByPrice:
		return nil

	default:
		return apperrors.NewBadRequest(
			"Invalid sortBy",
		)
	}
}

func validateSortOrder(value string) error {
	switch strings.ToLower(value) {
	case SortOrderAsc, SortOrderDesc:
		return nil

	default:
		return apperrors.NewBadRequest(
			"Invalid sortOrder",
		)
	}
}
