package errorcode

const (
	ShopIdInvalid                  = 406001
	ShopIdInvalidMessage           = "店铺 ID 无效"
	CreateCouponInvalid            = 406002
	CreateCouponInvalidMessage     = "创建优惠券参数有误"
	CreateCouponFailed             = 406003
	CreateCouponFailedMessage      = "创建优惠券错误"
	CouponIdInvalid                = 406004
	CouponIdInvalidMessage         = "优惠券 ID 无效"
	GetCouponInfoError             = 406005
	GetCouponInfoErrorMessage      = "获取优惠券详情错误"
	EditCouponFailed               = 406006
	EditCouponFailedMessage        = "更新优惠券错误"
	StuffCouponFiled               = 406007
	StuffCouponFiledMessage        = "添加优惠券错误"
	CouponIsUsed                   = 406008
	CouponIsUsedMessage            = "优惠券已经被使用"
	CouponIsExpired                = 406008
	CouponIsCouponIsExpiredMessage = "优惠券已经过期了"
	CouponCodeIsInvalid            = 406009
	CouponCodeIsInvalidMessage     = "优惠劵码校验失败"
	VerifyCouponUnSuccess          = 406010
	VerifyCouponUnSuccessMessage   = "优惠劵码校验失败, 或许已经校验成功"
	CronExpireCouponError          = 406011
	CronExpireCouponErrorMessage   = "自动过期优惠券异常"
)
