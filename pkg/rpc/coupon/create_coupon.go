package coupon

import (
	"fmt"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	idcreator "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/idcreator/snowflake"
	context "golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/utils/log"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/utils/validate"

	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/protos"
)

func (_ Controller) CreateCoupon(ctx context.Context, req *protos.CreateCouponRequest) (reply *protos.CouponInfoReply, err error) {

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

	couponId, errCreate := CreateCoupon(data)
	if errCreate != nil {

		log.GetLogger().WithError(err).Error("create coupon error")
		return &protos.CouponInfoReply{
			Err: &protos.Error{
				Code:    errorcode.CreateCouponFailed,
				Message: errorcode.CreateCouponFailedMessage,
				Stack:   nil,
			},
		}, nil
	}

	return &protos.CouponInfoReply{
		CouponId: couponId,
	}, nil

}

func CreateCoupon(data *model.Coupon) (uint64, error) {

	var err error
	data.CouponId = idcreator.NextID()
	err = database.GetDB(constant.DatabaseConfigKey).Create(&data).Error

	if err != nil {

		log.GetLogger().WithField("shop", data).WithError(err).Error("create record error")
		return 0, err
	}

	return data.CouponId, nil
}
