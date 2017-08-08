package routers

import (
	"beego_web/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/notify", &controllers.NotifyController{}, "*:Notify")
	beego.Router("/order", &controllers.OrderController{}, "*:Order")
}
