package user_coupon

import (
	"math/rand"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/coupon"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/protos"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func ChangeModelToProtobuf(couponInfo *model.Coupon, userCoupon *model.UserCoupon) *protos.UserCouponInfo {
	return &protos.UserCouponInfo{
		CouponInfo: &protos.CouponInfo{
			Type:         uint64(couponInfo.Type),
			Name:         couponInfo.Name,
			Instructions: couponInfo.Instructions,
			CouponId:     couponInfo.CouponId,
		},
		ExpireTime:       userCoupon.ExpireTime.Format(coupon.BaseTimeFormat),
		Status:           uint32(userCoupon.Status),
		VerificationTime: userCoupon.VerificationTime.Format(coupon.BaseTimeFormat),
		CreatedAt:        userCoupon.CreatedAt.Format(coupon.BaseTimeFormat),
		Code:             userCoupon.Code,
		Salt:             userCoupon.Salt,
		UserId:           userCoupon.UserId,
	}
}

func RandStringRunes(n int) string {

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GetVerifiedCouponListByCondition(conditions map[string]interface{}, pageNum, pageSize uint64) ([]*model.UserVerificationCoupon, uint64, error) {
	db := database.GetDB(constant.DatabaseConfigKey)

	db = db.Where("shopId = ?", conditions["shopId"]).Model(&model.UserVerificationCoupon{}).Order("createdAt desc")
	if !validation.IsEmpty(conditions["couponIds"]) {
		db = db.Where("couponIds in (?)", conditions["couponIds"])
	}

	if !validation.IsEmpty(conditions["startTime"]) && !validation.IsEmpty(conditions["endTime"]) {
		db = db.Where("createdAt BETWEEN ? AND ?", conditions["startTime"], conditions["endTime"])
	}

	if !validation.IsEmpty(conditions["verificationStatus"]) {
		db = db.Where("verificationStatus = ?", conditions["verificationStatus"])
	}

	if !validation.IsEmpty(conditions["verifierUserId"]) {
		db = db.Where("verifierUserId = ?", conditions["verifierUserId"])
	}

	results := make([]*model.UserVerificationCoupon, 0, pageSize)

	var count uint64

	pageStruct := constant.PageStruct{
		PageSize: pageSize,
		PageNum:  pageNum,
	}

	count, err := constant.FindPage(db, pageStruct, &results)
	return results, count, err

}
