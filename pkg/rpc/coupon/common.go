package coupon

import (
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/protos"
)

const BaseTimeFormat = "2006-01-02 15:04:05"

func ChangeModelToProtobuf(info *model.Coupon) *protos.CouponInfo {
	return &protos.CouponInfo{
		ShopId:       info.ShopId,
		Type:         uint64(info.Type),
		Name:         info.Name,
		ValidityType: uint64(info.ValidityType),
		RelativeTime: uint64(info.RelativeTime),
		StartTime:    info.StartTime.Format(BaseTimeFormat),
		EndTime:      info.EndTime.Format(BaseTimeFormat),
		Instructions: info.Instructions,
		CouponId:     info.CouponId,
		CreatedAt:    info.CreatedAt.Format(BaseTimeFormat),
	}
}
