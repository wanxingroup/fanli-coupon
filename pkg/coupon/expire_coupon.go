package coupon

import (
	"time"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"

	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/utils/log"
)

// 处理优惠券过期任务
func ExpireCoupon() error {
	now := time.Now()
	dbError := database.GetDB(constant.DatabaseConfigKey).Model(&model.UserCoupon{}).Where("status = ? and expireTime < ?", UserCouponStatusUnused, now).Updates(&model.UserCoupon{
		Status:           UserCouponStatusExpire,
		VerificationTime: now,
	}).Error

	if dbError != nil {
		log.GetLogger().WithField("Auto expire Time:", now).WithError(dbError).Error("Auto expire Time error")

		return dbError
	}

	return nil
}
