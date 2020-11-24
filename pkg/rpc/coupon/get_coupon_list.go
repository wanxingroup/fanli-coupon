package coupon

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	context "golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/coupon"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/utils/log"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

func (_ Controller) GetCouponList(ctx context.Context, req *protos.GetCouponListRequest) (reply *protos.GetCouponListReply, err error) {

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
		conditions["name"] = req.GetName()
	}

	if !validation.IsEmpty(req.GetStartTime()) && !validation.IsEmpty(req.GetEndTime()) {
		conditions["startTime"] = req.GetStartTime()
		conditions["endTime"] = req.GetEndTime()
	}

	pageData := map[string]uint64{
		"page":     page,
		"pageSize": pageSize,
	}

	data, count, err := coupon.GetCouponListByCondition(conditions, pageData)

	if err != nil {
		logger.WithError(err).Error("get user list failed")
		return &protos.GetCouponListReply{
			Err: &protos.Error{
				Code:    errorcode.InternalError,
				Message: err.Error(),
			},
		}, nil
	}

	couponInfoList := make([]*protos.CouponInfo, len(data))
	for key, info := range data {
		couponInfoList[key] = ChangeModelToProtobuf(info)
	}

	return &protos.GetCouponListReply{
		CouponList: couponInfoList,
		Count:      count,
	}, nil

}
