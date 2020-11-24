package validate

import (
	"time"

	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/protos"
)

const couponNameMaxLength = 20
const baseFormat = "2006-01-02 15:04:05"

func ValidateCreateCouponInfo(req *protos.CreateCouponRequest, isUpdate bool) (*model.Coupon, *protos.Error) {

	if req.ShopId == 0 {
		return nil, &protos.Error{
			Code:    errorcode.ShopIdInvalid,
			Message: errorcode.ShopIdInvalidMessage,
			Stack:   nil,
		}
	}

	if req.Type == 0 {
		return nil, &protos.Error{
			Code:    errorcode.CreateCouponInvalid,
			Message: errorcode.CreateCouponInvalidMessage,
			Stack:   nil,
		}
	}

	if req.Name == "" || len([]rune(req.Name)) > couponNameMaxLength {
		return nil, &protos.Error{
			Code:    errorcode.CreateCouponInvalid,
			Message: errorcode.CreateCouponInvalidMessage,
			Stack:   nil,
		}
	}

	if req.ValidityType == 0 && req.RelativeTime <= 0 {
		return nil, &protos.Error{
			Code:    errorcode.CreateCouponInvalid,
			Message: errorcode.CreateCouponInvalidMessage,
			Stack:   nil,
		}
	}

	returnModel := &model.Coupon{
		ShopId:       req.ShopId,
		Type:         uint8(req.Type),
		Name:         req.Name,
		ValidityType: uint8(req.ValidityType),
		RelativeTime: uint16(req.RelativeTime),
		Instructions: req.Instructions,
		CouponId:     req.CouponId,
	}

	if req.ValidityType == 1 {

		startTime, err := time.Parse(baseFormat, req.StartTime)
		endTime, err2 := time.Parse(baseFormat, req.EndTime)

		if err != nil || err2 != nil {
			return nil, &protos.Error{
				Code:    errorcode.CreateCouponInvalid,
				Message: errorcode.CreateCouponInvalidMessage,
				Stack:   nil,
			}
		}

		returnModel.StartTime = startTime
		returnModel.EndTime = endTime
	} else {
		returnModel.StartTime = time.Unix(0, 0)
		returnModel.EndTime = time.Unix(0, 0)
	}

	return returnModel, nil
}
