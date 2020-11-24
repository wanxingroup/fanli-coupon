package model

import (
	databases "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database/models"
)

const TableNameUserCheckCoupon = "user_check_coupon"
const GetUserCheckCouponDefaultPageSize = 20

type UserVerificationCoupon struct {
	Code               uint64 `gorm:"column:code;type:bigint unsigned;not null;primary_key;comment:'券码'"`
	CouponId           uint64 `gorm:"column:couponId;type:bigint unsigned;index:couponId;comment:'优惠券 ID'"`
	ShopId             uint64 `gorm:"column:shopId;type:bigint unsigned;index:shopId;comment:'店铺 ID'"`
	UserId             uint64 `gorm:"column:userId;type:bigint unsigned;not null;index:userId;comment:'券码拥有者 UID'"`
	VerifierUserId     uint64 `gorm:"column:verifierUserId;type:bigint unsigned;not null;index:checkUserId;comment:'核销员 UID'"`
	VerificationType   uint8  `gorm:"column:verificationType;type:tinyint unsigned;not null;default:0;index:verifierType;comment:'核销方式: 0: 扫码核销 1: 下单核销'"`
	VerificationStatus uint8  `gorm:"column:verificationStatus;type:tinyint unsigned;not null;default:0;index:verificationStatus;comment:'核销时优惠券的状态 0: 正常核销 1:过期自动核销'"`
	databases.Time
}
