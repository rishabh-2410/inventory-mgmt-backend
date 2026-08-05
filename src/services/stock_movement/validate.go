package services

import (
	"errors"
	"inventory-backend/src/constants"
	"strings"
	"time"
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

func ValidateSortBy(sortBy *string) (*string, error) {
	if sortBy == nil {
		return nil, nil
	}

	if strings.ToLower(*sortBy) != "created_at_desc" && strings.ToLower(*sortBy) != "created_at_asc" && strings.ToLower(*sortBy) != "movement_quantity_desc" && strings.ToLower(*sortBy) != "movement_quantity_asc" {
		return nil, errors.New("invalid sort by")
	}

	return sortBy, nil
}

func ValidateTimeRange(from *time.Time, to *time.Time) (*time.Time, *time.Time, error) {
	if from == nil || to == nil {
		return nil, nil, nil
	}

	if to.Before(*from) || from.After(*to) {
		return nil, nil, errors.New("invalid time range")
	}

	return from, to, nil
}

func ValidateMovementType(movementType *string) (*string, error) {
	if movementType == nil {
		return nil, nil
	}

	if *movementType != constants.MOVEMENT_SALE && *movementType != constants.MOVEMENT_RETURN && *movementType != constants.MOVEMENT_DAMAGE && *movementType != constants.MOVEMENT_RECEIVE {
		return nil, errors.New("invalid movement type")
	}

	return movementType, nil
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
