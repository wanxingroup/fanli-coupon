package constant

import (
	"github.com/jinzhu/gorm"
)

const DatabaseConfigKey = "common"

type PageStruct struct {
	PageNum, PageSize uint64
}

func FindPage(db *gorm.DB, pageStruct PageStruct, out interface{}) (uint64, error) {

	if pageStruct.PageNum < 0 || pageStruct.PageSize < 0 {
		return 0, nil
	}

	if pageStruct.PageNum > 0 && pageStruct.PageSize > 0 {
		db = db.Offset((pageStruct.PageNum - 1) * pageStruct.PageSize)
	}

	if pageStruct.PageSize > 0 {
		db = db.Limit(pageStruct.PageSize)
	}

	err := db.Find(out).Error

	if err != nil {
		return 0, err
	}

	var count uint64
	err = db.Count(&count).Error

	return count, err

}
