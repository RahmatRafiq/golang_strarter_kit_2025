package scopes

import (
	"golang_starter_kit_2025/app/requests"

	"gorm.io/gorm"
)

func Paginate(filter requests.FilterRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := 1
		if filter.Page != nil {
			page = *filter.Page
		}

		limit := 10
		if filter.Limit != nil {
			limit = *filter.Limit
		}
		switch {
		case limit > 100:
			limit = 100
		case limit <= 10:
			limit = 10
		}
		offset := 0
		if filter.Offset != nil {
			offset = *filter.Offset
		}
		if offset <= 0 {
			if filter.Page != nil {
				offset = (page - 1) * limit
			} else {
				offset = 0
			}
		}

		return db.Offset(offset).Limit(limit)
	}
}
