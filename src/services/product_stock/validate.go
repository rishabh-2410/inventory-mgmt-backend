package services

import (
	"errors"
	"strings"
)

func ValidateLimit(limit *int) (*int, error) {
	if limit == nil {
		return nil, nil
	}

	if *limit < 1 {
		return nil, errors.New("limit must be greater than 0")
	}

	if *limit > 100 {
		return nil, errors.New("limit must be less than 100")
	}

	return limit, nil
}

func ValidatePage(page *int) (*int, error) {
	if page == nil {
		return nil, nil
	}

	if *page < 1 {
		return nil, errors.New("page must be greater than 0")
	}

	if *page > 100 {
		return nil, errors.New("page must be less than 100")
	}

	return page, nil
}

func ValidateSearchParam(search *string) (*string, error) {
	if search == nil {
		return nil, nil
	}

	searchParam := strings.TrimSpace(*search)
	if searchParam == "" {
		return nil, errors.New("search parameter must be a non-empty string")
	}

	return &searchParam, nil
}
