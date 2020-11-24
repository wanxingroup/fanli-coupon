package application

import (
	"math/rand"
	"time"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/launcher"
	idcreator "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/idcreator/snowflake"
	"github.com/spf13/viper"

	"github.com/robfig/cron/v3"

	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/config"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/cronjob"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/coupon"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/rpc/user_coupon"
	cronUtils "dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/utils/cron"
	"dev-gitlab.wanxingrowth.com/fanli/coupon/pkg/utils/log"
)

var c *cron.Cron

func Start() {

	app := launcher.NewApplication(
		launcher.SetApplicationDescription(
			&launcher.ApplicationDescription{
				ShortDescription: "users service",
				LongDescription:  "support mini program user data management function.",
			},
		),
		launcher.SetApplicationLogger(log.GetLogger()),
		launcher.SetApplicationEvents(
			launcher.NewApplicationEvents(
				launcher.SetOnInitEvent(func(app *launcher.Application) {

					unmarshalConfiguration()

					registerCronJobInit()

					registerUserRPCRouter(app)
					idcreator.InitCreator(app.GetServiceId())
					randInit()
				}),
				launcher.SetOnStartEvent(func(app *launcher.Application) {

					autoMigration()
				}),
			),
		),
	)

	app.Launch()
}

func registerUserRPCRouter(app *launcher.Application) {

	rpcService := app.GetRPCService()
	if rpcService == nil {

		log.GetLogger().WithField("stage", "onInit").Error("get rpc service is nil")
		return
	}

	protos.RegisterShopCouponControllerServer(rpcService.GetRPCConnection(), &coupon.Controller{})
	protos.RegisterUserCouponControllerServer(rpcService.GetRPCConnection(), &user_coupon.Controller{})
}

func unmarshalConfiguration() {
	err := viper.Unmarshal(config.Config)
	if err != nil {

		log.GetLogger().WithError(err).Error("unmarshal config error")
	}
}

func randInit() {
	rand.Seed(time.Now().UnixNano())
}

func registerCronJobInit() {

	c = cron.New(
		cron.WithParser(cron.NewParser(cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor)),
		cron.WithLogger(cronUtils.NewLogger(log.GetLogger())),
	)

	type Cron struct {
		Spec string
		Run  func()
	}

	jobs := []Cron{
		{Spec: "*/1 * * * *", Run: cronjob.ExpireCoupon},
	}

	for _, job := range jobs {
		_, err := c.AddFunc(job.Spec, job.Run)
		if err != nil {
			log.GetLogger().WithError(err).Error("register cron job error")
		}
	}
	c.Start()
}
