package user_coupon

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	context "golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/coupon"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/utils/log"
)

func (_ Controller) GetUserCouponList(ctx context.Context, req *protos.GetUserCouponListRequest) (reply *protos.GetUserCouponListReply, err error) {

	logger := log.GetLogger().WithField("requestData", req)

	if req == nil {
		log.GetLogger().Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}

	conditions := map[string]interface{}{
		"shopId": req.GetShopId(),
		"userId": req.GetUserId(),
	}

	if !validation.IsEmpty(req.GetStatus()) {
		conditions["status"] = req.GetStatus()
	} else {
		conditions["status"] = coupon.UserCouponStatusUnused
	}

	if !validation.IsEmpty(req.GetLastCode()) {
		conditions["lastCode"] = req.GetLastCode()
	}

	var pageSize uint32 = model.GetUserDefaultPageSize
	if !validation.IsEmpty(req.GetPageSize()) {
		pageSize = req.GetPageSize()
	}

	userCoupons, count, err := coupon.GetUserCouponListByCondition(conditions, pageSize)

	if err != nil {
		logger.WithError(err).Error("get user coupons list failed")
		return &protos.GetUserCouponListReply{
			Err: &protos.Error{
				Code:    errorcode.InternalError,
				Message: err.Error(),
			},
		}, nil
	}
	var couponIds []uint64

	for _, userCoupon := range userCoupons {
		couponIds = append(couponIds, userCoupon.CouponId)
	}

	couponsList, err2 := coupon.GetCouponListByCouponIds(couponIds)
	if err2 != nil {
		logger.WithError(err2).Error("get coupons list by id failed")
	}

	CouponListByIds := make(map[uint64]*model.Coupon)

	if len(couponsList) > 0 {
		for _, m := range couponsList {
			CouponListByIds[m.CouponId] = m
		}
	}

	returnCouponList := make([]*protos.UserCouponInfo, len(userCoupons))

	for key, userCoupon := range userCoupons {
		if _, ok := CouponListByIds[userCoupon.CouponId]; ok {
			returnCouponList[key] = ChangeModelToProtobuf(CouponListByIds[userCoupon.CouponId], userCoupon)
		}
	}

	return &protos.GetUserCouponListReply{
		UserCouponInfo: returnCouponList,
		Count:          uint32(count),
	}, nil

}
