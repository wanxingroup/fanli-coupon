package user_coupon

import (
	"fmt"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	"github.com/jinzhu/gorm"

	context "golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/coupon"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/utils/log"
)

func (_ Controller) GetCouponInfoByCode(ctx context.Context, req *protos.CouponInfoByCodeRequest) (reply *protos.GetUserCouponInfoReply, err error) {

	if req == nil {
		log.GetLogger().Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}

	returnCouponInfo, err2 := GetUserCouponInfoByCode(req.GetCode(), req.GetShopId())
	if err2 != nil {
		return &protos.GetUserCouponInfoReply{
			Err: &protos.Error{
				Code:    err2.Code,
				Message: err2.Message,
				Stack:   err2.Stack,
			},
		}, nil
	}

	return &protos.GetUserCouponInfoReply{
		UserCouponInfo: returnCouponInfo,
		Err:            nil,
	}, nil
}

func GetUserCouponInfoByCode(code, shopId uint64) (*protos.UserCouponInfo, *protos.Error) {

	logger := log.GetLogger().WithField("GetUserCouponInfoByCode : ", code)
	results := &model.UserCoupon{}

	// 查询用户领劵数据
	err2 := database.GetDB(constant.DatabaseConfigKey).Where("code = ? and shopId = ?", code, shopId).
		First(results).Error

	if err2 == gorm.ErrRecordNotFound {
		logger.Info("user coupon record not found")
		return nil, &protos.Error{
			Code:    errorcode.GetCouponInfoError,
			Message: errorcode.GetCouponInfoErrorMessage,
			Stack:   nil,
		}
	}

	if err2 != nil {
		logger.Error("get user coupons info by code failed")
		return nil, &protos.Error{
			Code:    errorcode.GetCouponInfoError,
			Message: errorcode.GetCouponInfoErrorMessage,
			Stack:   nil,
		}
	}

	// 获取优惠券的数据
	couponInfo, err3 := coupon.GetCouponInfo(results.CouponId)
	if err3 != nil {
		return nil, &protos.Error{
			Code:    errorcode.GetCouponInfoError,
			Message: errorcode.GetCouponInfoErrorMessage,
			Stack:   nil,
		}
	}

	return ChangeModelToProtobuf(couponInfo, results), nil
}
