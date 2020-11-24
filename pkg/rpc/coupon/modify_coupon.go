package coupon

import (
	"fmt"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	context "golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/utils/log"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/utils/validate"

	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/protos"
)

func (_ Controller) ModifyCoupon(ctx context.Context, req *protos.CreateCouponRequest) (reply *protos.CouponInfoReply, err error) {

	if req == nil {
		log.GetLogger().Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}

	data, errMessage := validate.ValidateCreateCouponInfo(req, false)

	if errMessage != nil {
		return &protos.CouponInfoReply{
			Err: &protos.Error{
				Code:    errMessage.GetCode(),
				Message: errMessage.GetMessage(),
				Stack:   errMessage.GetStack(),
			},
		}, nil
	}

	couponId, errCreate := ModifyCoupon(data)
	if errCreate != nil {

		log.GetLogger().WithError(err).Error("update record error")
		return &protos.CouponInfoReply{
			Err: &protos.Error{
				Code:    errorcode.EditCouponFailed,
				Message: errorcode.EditCouponFailedMessage,
				Stack:   nil,
			},
		}, nil
	}

	return &protos.CouponInfoReply{
		CouponId: couponId,
	}, nil

}

func ModifyCoupon(data *model.Coupon) (uint64, error) {

	var err error

	uploadData := map[string]interface{}{
		"type":         data.Type,
		"name":         data.Name,
		"validityType": data.ValidityType,
		"relativeTime": data.RelativeTime,
		"startTime":    data.StartTime,
		"endTime":      data.EndTime,
		"instructions": data.Instructions,
	}

	err = database.GetDB(constant.DatabaseConfigKey).Model(&data).Where("couponId = ? and shopId = ?", data.CouponId, data.ShopId).Updates(uploadData).Error

	if err != nil {

		log.GetLogger().WithField("shop", data).WithError(err).Error("update record error")
		return 0, err
	}

	return data.CouponId, nil
}
