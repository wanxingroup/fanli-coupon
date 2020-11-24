package model

import (
	"time"

	databases "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database/models"
)

const TableNameCoupon = "coupon"

type Coupon struct {
	CouponId     uint64    `gorm:"column:couponId;type:bigint unsigned;primary_key;index:idx_shop_coupon;comment:'优惠券 ID'"`
	ShopId       uint64    `gorm:"column:shopId;type:bigint unsigned;index:shopId,idx_shop_coupon;comment:'店铺 ID'"`
	Type         uint8     `gorm:"column:type;type:tinyint unsigned;not null;default:'0';comment:'类别(1代金券)'"`
	Name         string    `gorm:"column:name;type:varchar(11);not null;default:'';comment:'优惠券名称'"`
	ValidityType uint8     `gorm:"column:validityType;type:tinyint unsigned;not null;default:'0';comment:'有效期类型:0相对时间1绝对时间'"`
	RelativeTime uint16    `gorm:"column:relativeTime;type:smallint(5) unsigned;not null;default:'0';comment:'有效期相对时间 小时计算'"`
	StartTime    time.Time `gorm:"column:startTime;null;comment:'绝对时间开始'"`
	EndTime      time.Time `gorm:"column:endTime;null;comment:'绝对时间结束'"`
	Instructions string    `gorm:"column:instructions;type:varchar(150);not null;default:'';comment:'使用须知'"`
	databases.Time
}

func (coupon *Coupon) TableName() string {

	return TableNameCoupon
}
