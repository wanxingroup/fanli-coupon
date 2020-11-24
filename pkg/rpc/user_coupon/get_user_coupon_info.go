package user_coupon

import (
	"fmt"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	context "golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/coupon"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/utils/log"
)

func (_ Controller) GetUserCouponInfo(ctx context.Context, req *protos.GetUserCouponInfoRequest) (reply *protos.GetUserCouponInfoReply, err error) {

	logger := log.GetLogger().WithField("requestData", req)

	if req == nil {
		log.GetLogger().Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}

	conditions := map[string]interface{}{
		"shopId": req.GetShopId(),
		"userId": req.GetUserId(),
		"code":   req.GetCode(),
	}
	results := &model.UserCoupon{}

	// 查询用户领劵数据
	err = database.GetDB(constant.DatabaseConfigKey).Where("code = ? and shopId = ? and userId = ? ", conditions["code"], conditions["shopId"], conditions["userId"]).
		First(results).Error

	if err != nil {
		logger.WithError(err).Error("get user coupons list failed")
		return &protos.GetUserCouponInfoReply{
			Err: &protos.Error{
				Code:    errorcode.InternalError,
				Message: err.Error(),
			},
		}, nil
	}

	// 获取优惠券的数据
	couponInfo, err2 := coupon.GetCouponInfo(results.CouponId)
	if err2 != nil {
		return &protos.GetUserCouponInfoReply{
			Err: &protos.Error{
				Code:    err2.Code,
				Message: err2.Message,
				Stack:   err2.Stack,
			},
		}, nil
	}

	returnCouponInfo := ChangeModelToProtobuf(couponInfo, results)

	return &protos.GetUserCouponInfoReply{
		UserCouponInfo: returnCouponInfo,
		Err:            nil,
	}, nil
}
