package application

import (
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"

	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/model"
)

func autoMigration() {

	db := database.GetDB(constant.DatabaseConfigKey)
	db.AutoMigrate(model.Coupon{})
	db.AutoMigrate(model.UserCoupon{})
	db.AutoMigrate(model.UserVerificationCoupon{})
}
