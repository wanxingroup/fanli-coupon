package cronjob

import (
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/coupon"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/utils/log"
)

func ExpireCoupon() {
	err := coupon.ExpireCoupon()
	if err != nil {
		log.GetLogger().WithError(err).Error("expire coupon error")
		return
	}

	return
}
