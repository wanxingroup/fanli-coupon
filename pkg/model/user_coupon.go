package model

import (
	"time"

	databases "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database/models"
)

const TableNameUserCoupon = "user_coupon"
const GetUserDefaultPageSize = 20

type UserCoupon struct {
	Code             uint64    `gorm:"column:code;type:bigint unsigned;not null;primary_key;comment:'券码'"`
	Salt             string    `gorm:"column:salt;type:char(5);not null;comment:'券码的盐'"`
	ShopId           uint64    `gorm:"column:shopId;type:bigint unsigned;index:shopId;comment:'店铺 ID'"`
	CouponId         uint64    `gorm:"column:couponId;type:bigint unsigned;index:couponId;comment:'优惠券 ID'"`
	UserId           uint64    `gorm:"column:userId;type:bigint unsigned;index:userId;comment:'用户 ID'"`
	ExpireTime       time.Time `gorm:"column:expireTime;not null;index:idx_status_expireTime;comment:'优惠券失效时间'"`
	Status           uint8     `gorm:"column:status;type:tinyint unsigned;not null;index:idx_status_expireTime;default:'0';comment:'状态  0表示未使用，1表示已使用 2已过期'"`
	VerificationTime time.Time `gorm:"column:verificationTime;null;comment:'优惠券核销时间'"`
	databases.Time
}

func (coupon *UserCoupon) TableName() string {

	return TableNameUserCoupon
}
