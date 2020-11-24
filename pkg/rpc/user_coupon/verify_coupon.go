package user_coupon

import (
	"fmt"
	"time"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/idcreator/verifiable"
	context "golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/coupon"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/utils/log"
)

func (_ Controller) VerifyCoupon(ctx context.Context, req *protos.VerifyCouponRequest) (reply *protos.VerifyCouponReply, err error) {

	logger := log.GetLogger().WithField("requestData", req)

	if req == nil {
		log.GetLogger().Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}

	// 获取券码对应的 优惠券信息
	userCouponInfo, err2 := GetUserCouponInfoByCode(req.GetCode(), req.GetShopId())

	if err2 != nil {
		logger.Error("get user coupons info error")
		return &protos.VerifyCouponReply{
			Err: &protos.Error{
				Code:    errorcode.GetCouponInfoError,
				Message: errorcode.GetCouponInfoErrorMessage,
				Stack:   nil,
			},
		}, nil
	}

	// 校验优惠券是否正确
	if !checkCodeRight(userCouponInfo.Code, userCouponInfo.Salt) {
		logger.Info("check user coupons salt error")
		return &protos.VerifyCouponReply{
			Err: &protos.Error{
				Code:    errorcode.CouponCodeIsInvalid,
				Message: errorcode.CouponCodeIsInvalidMessage,
				Stack:   nil,
			},
		}, nil
	}

	switch userCouponInfo.Status {
	case coupon.UserCouponStatusUsed:
		logger.WithField("code", req.GetCode()).Info("user coupon is used")
		return &protos.VerifyCouponReply{
			Err: &protos.Error{
				Code:    errorcode.CouponIsUsed,
				Message: errorcode.CouponIsUsedMessage,
				Stack:   nil,
			},
		}, nil

	case coupon.UserCouponStatusExpire:
		logger.WithField("code", req.GetCode()).Info("user coupon is expire")
		return &protos.VerifyCouponReply{
			Err: &protos.Error{
				Code:    errorcode.CouponIsExpired,
				Message: errorcode.CouponIsCouponIsExpiredMessage,
				Stack:   nil,
			},
		}, nil
	}

	// 核销
	if ok, err := verifyCoupon(req, userCouponInfo.UserId, userCouponInfo.CouponInfo); !ok {
		return &protos.VerifyCouponReply{
			Err: &protos.Error{
				Code:    err.Code,
				Message: err.Message,
				Stack:   nil,
			},
		}, nil
	}

	// 重新获取优惠券记录
	userCouponInfo, _ = GetUserCouponInfoByCode(req.GetCode(), req.GetShopId())

	return &protos.VerifyCouponReply{
		UserCouponInfo: userCouponInfo,
		Err:            nil,
	}, nil
}

func checkCodeRight(code uint64, SecretKey string) bool {
	return verifiable.Verify(code, SecretKey)
}

func verifyCoupon(req *protos.VerifyCouponRequest, userId uint64, info *protos.CouponInfo) (bool, *protos.Error) {
	// 核销

	tx := database.GetDB(constant.DatabaseConfigKey).Begin()
	query := tx.Model(&model.UserCoupon{}).Where("code = ?", req.GetCode()).Updates(&model.UserCoupon{
		Status:           coupon.UserCouponStatusUsed,
		VerificationTime: time.Now(),
	})

	if query.RowsAffected == 0 {
		log.GetLogger().WithField("code", req.GetCode()).Error("verify user code error")
		tx.Rollback()
		return false, &protos.Error{
			Code:    errorcode.VerifyCouponUnSuccess,
			Message: errorcode.VerifyCouponUnSuccessMessage,
			Stack:   nil,
		}
	}

	// 落库
	recode := &model.UserVerificationCoupon{
		Code:           req.GetCode(),
		UserId:         userId,
		VerifierUserId: req.GetVerifierUserId(),
		ShopId:         req.GetShopId(),
		CouponId:       info.CouponId,
	}

	createError := tx.Create(recode).Error
	if createError != nil {
		tx.Rollback()
		log.GetLogger().WithField("createVerifier", req.GetVerifierUserId()).WithError(createError).Error("create Verifier record error")
		return false, &protos.Error{
			Code:    errorcode.VerifyCouponUnSuccess,
			Message: errorcode.VerifyCouponUnSuccessMessage,
			Stack:   nil,
		}
	}

	tx.Commit()
	return true, nil
}
