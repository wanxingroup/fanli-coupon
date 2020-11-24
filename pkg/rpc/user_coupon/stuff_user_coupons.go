package user_coupon

import (
	"fmt"
	"time"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	context "golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/idcreator/verifiable"

	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/constant"
	pkgCoupon "dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/coupon"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/utils/log"
)

func (ctr Controller) StuffUserCoupons(ctx context.Context, req *protos.StuffUserCouponRequest) (reply *protos.StuffUserCouponReply, err error) {

	logger := log.GetLogger().WithField("requestData", req)

	if req == nil {
		log.GetLogger().Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}

	if len(req.StuffCouponStuck) <= 0 {
		log.GetLogger().Error("request StuffCouponStuck is nil")
		return nil, fmt.Errorf("request StuffCouponStuck is nil")
	}

	var couponIds []uint64
	for _, scs := range req.StuffCouponStuck {
		couponIds = append(couponIds, scs.CouponId)
	}

	couponsList, err2 := pkgCoupon.GetCouponListByCouponIds(couponIds)
	if err2 != nil {
		logger.WithError(err2).Error("get coupons list by id failed")
	}
	CouponListByIds := make(map[uint64]*model.Coupon)

	if len(couponsList) > 0 {
		for _, m := range couponsList {
			CouponListByIds[m.CouponId] = m
		}
	}

	idcConfig := verifiable.Settings{
		StartTime: time.Time{},
		SecretKey: RandStringRunes(5),
	}

	ids, err3 := createUserCouponData(req, CouponListByIds, idcConfig)

	if err3 != nil {
		return &protos.StuffUserCouponReply{
			Err: &protos.Error{
				Code:    errorcode.StuffCouponFiled,
				Message: errorcode.StuffCouponFiledMessage,
				Stack:   nil,
			},
		}, nil
	}

	return &protos.StuffUserCouponReply{
		Codes: ids,
	}, nil

}

func createUserCouponData(req *protos.StuffUserCouponRequest, coupon map[uint64]*model.Coupon, idcConfig verifiable.Settings) ([]uint64, error) {

	stuffData := req.GetStuffCouponStuck()
	if len(stuffData) <= 0 {
		return nil, nil
	}

	idC := verifiable.NewIDCreator(idcConfig)

	tx := database.GetDB(constant.DatabaseConfigKey).Begin()

	var returnIds []uint64

	for _, datum := range stuffData {
		for datum.Nums > 0 {
			record := &model.UserCoupon{
				Code:             idC.NextID(),
				Salt:             idcConfig.SecretKey,
				ShopId:           coupon[datum.CouponId].ShopId,
				CouponId:         coupon[datum.CouponId].CouponId,
				UserId:           req.GetUserId(),
				ExpireTime:       CalcExpireTime(coupon[datum.CouponId]),
				Status:           pkgCoupon.UserCouponStatusUnused,
				VerificationTime: time.Unix(0, 0),
			}

			if err := tx.Create(&record).Error; err != nil {
				log.GetLogger().WithField("createVerifier", record).WithError(err).Error("create createVerifier error")
				tx.Rollback()
				return nil, err
			}
			datum.Nums--
			returnIds = append(returnIds, record.CouponId)
		}
	}

	tx.Commit()
	return returnIds, nil
}

// 获取优惠券的过期时间
func CalcExpireTime(couponInfo *model.Coupon) time.Time {
	if couponInfo.ValidityType == 0 {
		timestamp := time.Now().Unix() + int64(couponInfo.RelativeTime*3600)
		return time.Unix(timestamp, 0)
	} else {
		return couponInfo.EndTime
	}
}
