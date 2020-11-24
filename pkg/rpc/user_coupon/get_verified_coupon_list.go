package user_coupon

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	context "golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/coupon"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/model"
	rpcCoupon "dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/coupon"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/utils/log"
)

const DefaultPage = 1
const DefaultPageSize = 20
const MaxPageSize = 100

func (_ Controller) GetVerificationCouponList(ctx context.Context, req *protos.GetVerificationCouponListRequest) (reply *protos.GetVerificationCouponListReply, err error) {

	logger := log.GetLogger().WithField("requestData", req)

	if req == nil {
		log.GetLogger().Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}
	var page = uint64(DefaultPage)
	var pageSize = uint64(DefaultPageSize)

	if req.GetPage() > 0 {
		page = req.GetPage()
	}

	if req.GetPageSize() > 0 {
		pageSize = req.GetPageSize()
	}

	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	conditions := map[string]interface{}{
		"shopId": req.GetShopId(),
	}

	if !validation.IsEmpty(req.GetName()) {

		//  查询优惠券的名称
		var couponIds []uint64

		couponConditions := map[string]interface{}{
			"shopId": req.GetShopId(),
			"name":   req.GetName(),
		}

		// 默认查询 100 条的优惠券出来
		pageData := map[string]uint64{
			"page":     DefaultPage,
			"pageSize": MaxPageSize,
		}

		couponInfoList, _, err := coupon.GetCouponListByCondition(couponConditions, pageData)
		if err != nil {
			logger.WithError(err).Error("get verified couponInfo with couponInfo by couponId list failed")
			return &protos.GetVerificationCouponListReply{
				Err: &protos.Error{
					Code:    errorcode.InternalError,
					Message: err.Error(),
				},
			}, nil
		}
		for _, couponInfo := range couponInfoList {
			couponIds = append(couponIds, couponInfo.CouponId)
		}
		conditions["couponIds"] = couponIds
	}

	if !validation.IsEmpty(req.GetVerificationStatus()) {
		conditions["verificationStatus"] = req.GetVerificationStatus()
	}

	if !validation.IsEmpty(req.GetVerifierUserId()) {
		conditions["verifierUserId"] = req.GetVerifierUserId()
	}

	if !validation.IsEmpty(req.GetStartTime()) && !validation.IsEmpty(req.GetEndTime()) {
		conditions["startTime"] = req.GetStartTime()
		conditions["endTime"] = req.GetEndTime()
	}

	verifiedCouponList, count, err := GetVerifiedCouponListByCondition(conditions, page, pageSize)

	if err != nil {
		logger.WithError(err).Error("get user list failed")
		return &protos.GetVerificationCouponListReply{
			Err: &protos.Error{
				Code:    errorcode.InternalError,
				Message: err.Error(),
			},
		}, nil
	}

	var couponIds []uint64

	for _, userCoupon := range verifiedCouponList {
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

	returnCouponList := make([]*protos.VerificationCouponInfo, 0, len(verifiedCouponList))

	for _, verifiedCoupon := range verifiedCouponList {
		if _, ok := CouponListByIds[verifiedCoupon.CouponId]; ok {
			returnCouponList = append(returnCouponList, &protos.VerificationCouponInfo{
				Code:               verifiedCoupon.Code,
				Name:               CouponListByIds[verifiedCoupon.CouponId].Name,
				VerificationTime:   verifiedCoupon.CreatedAt.Format(rpcCoupon.BaseTimeFormat),
				CouponRelativeTime: uint64(CouponListByIds[verifiedCoupon.CouponId].RelativeTime),
				CouponType:         uint32(CouponListByIds[verifiedCoupon.CouponId].Type),
				VerificationType:   uint32(verifiedCoupon.VerificationType),
				VerifierUserId:     verifiedCoupon.VerifierUserId,
				ValidityType:       uint32(CouponListByIds[verifiedCoupon.CouponId].ValidityType),
				CouponId:           verifiedCoupon.CouponId,
			})
		}
	}

	return &protos.GetVerificationCouponListReply{
		VerificationCouponInfoList: returnCouponList,
		Count:                      uint32(count),
	}, nil
}
