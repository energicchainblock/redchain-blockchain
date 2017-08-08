package controllers

import (
	"beego_web/db"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
)

type NotifyController struct {
	beego.Controller
}

func (this *NotifyController) Notify() {
	sender := this.GetString("invoice_id")
	hash := this.GetString("transaction_hash")
	sv := this.GetString("value")
	test := this.GetString("test")
	val, _ := strconv.Atoi(sv)

	if test == "true" {
		return
	}

	fmt.Printf("sender=%v, hash=%v, sv=%v, test=%v\r\n", sender, hash, sv, test)

	// err := db.InsertOrder(sender, hash, val)
	// if err != nil {
	// 	this.Ctx.Output.Body([]byte("error"))
	// 	return
	// }

	this.Ctx.Output.Body([]byte("*ok*"))
}
