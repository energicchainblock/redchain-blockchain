package beego_web

import (
	_ "beego_web/routers"
	"github.com/astaxie/beego"
)

func Main() {
	beego.Run()
}
