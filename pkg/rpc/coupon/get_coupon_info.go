package coupon

import (
	"fmt"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	"github.com/jinzhu/gorm"
	context "golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/utils/log"
)

func (_ Controller) GetCouponInfo(ctx context.Context, req *protos.GetCouponInfoRequest) (reply *protos.GetCouponInfoReply, err error) {

	if req == nil {
		log.GetLogger().Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}

	if req.GetShopId() == 0 {
		return &protos.GetCouponInfoReply{
			Err: &protos.Error{
				Code:    errorcode.ShopIdInvalid,
				Message: errorcode.ShopIdInvalidMessage,
				Stack:   nil,
			},
		}, nil
	}

	if req.GetCouponId() == 0 {
		return &protos.GetCouponInfoReply{
			Err: &protos.Error{
				Code:    errorcode.CouponIdInvalid,
				Message: errorcode.CouponIdInvalidMessage,
				Stack:   nil,
			},
		}, nil
	}
	couponInfo, err2 := GetCouponInfo(req.GetCouponId())

	if err2 != nil {
		return &protos.GetCouponInfoReply{
			Err: &protos.Error{
				Code:    err2.Code,
				Message: err2.Message,
				Stack:   err2.Stack,
			},
		}, nil
	}
	returnCouponInfo := ChangeModelToProtobuf(couponInfo)

	return &protos.GetCouponInfoReply{
		CouponInfo: returnCouponInfo,
	}, nil

}

func GetCouponInfo(couponId uint64) (*model.Coupon, *protos.Error) {

	couponInfo := &model.Coupon{}
	err2 := database.GetDB(constant.DatabaseConfigKey).First(couponInfo, couponId).Error

	if gorm.IsRecordNotFoundError(err2) {
		log.GetLogger().Info("coupon rocord not found")
		return nil, nil
	}

	if err2 != nil {
		log.GetLogger().Error("get coupons info by code failed")
		return nil, &protos.Error{
			Code:    errorcode.GetCouponInfoError,
			Message: errorcode.GetCouponInfoErrorMessage,
			Stack:   nil,
		}
	}

	return couponInfo, nil
}
