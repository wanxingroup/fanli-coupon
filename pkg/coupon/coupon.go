package coupon

import (
	"fmt"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	baseDb "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/databases"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jinzhu/gorm"

	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/model"
)

const UserCouponStatusUnused = 0 // 未使用
const UserCouponStatusUsed = 1   // 已使用
const UserCouponStatusExpire = 2 // 已过期

func GetCouponListByCondition(conditions map[string]interface{}, pageData map[string]uint64) ([]*model.Coupon, uint64, error) {
	db := database.GetDB(constant.DatabaseConfigKey)

	db = db.Where("shopId = ?", conditions["shopId"]).
		Table(model.TableNameCoupon).
		Order("createdAt desc")
	if !validation.IsEmpty(conditions["name"]) {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", conditions["name"]))
	}

	if !validation.IsEmpty(conditions["startTime"]) && !validation.IsEmpty(conditions["endTime"]) {
		db = db.Where("createdAt BETWEEN ? AND ?", conditions["startTime"], conditions["endTime"])
	}

	results := make([]*model.Coupon, 0, pageData["pageSize"])

	var count uint64
	err := baseDb.FindPage(db, pageData, &results, &count)
	return results, count, err
}

// 根据主键 ID 获取记录
func GetCouponListByCouponIds(couponIds []uint64) ([]*model.Coupon, error) {

	if len(couponIds) <= 0 {
		return nil, nil
	}

	results := make([]*model.Coupon, 0)

	err2 := database.GetDB(constant.DatabaseConfigKey).Where("couponId in (?)", couponIds).Find(&results).Error

	if err2 != nil {
		return nil, err2
	}

	if gorm.IsRecordNotFoundError(err2) {
		return nil, nil
	}

	return results, nil
}

func GetUserCouponListByCondition(conditions map[string]interface{}, pageSize uint32) ([]*model.UserCoupon, uint64, error) {
	db := database.GetDB(constant.DatabaseConfigKey)

	db = db.Where("shopId = ? and userId = ? ", conditions["shopId"], conditions["userId"]).
		Model(&model.UserCoupon{}).
		Order("createdAt desc")

	if !validation.IsEmpty(conditions["status"]) || conditions["status"] == UserCouponStatusUnused {
		db = db.Where("status = ?", conditions["status"])
	}

	if !validation.IsEmpty(conditions["lastCode"]) {
		db = db.Where("code < ?", conditions["lastCode"])
	}

	results := make([]*model.UserCoupon, 0, pageSize)

	err := db.Limit(pageSize).Find(&results).Error

	if err != nil {
		return nil, 0, err
	}

	var count uint64
	err = db.Count(&count).Error

	return results, count, nil

}
