package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
	"net/url"
	//"strings"
)

type OrderErr struct {
	Code int    `json:"code"`
	Addr string `json:"address"`
	Msg  string `json:"msg"`
}

type OrderResponse struct {
	Addr string `json:"address"`
}

type OrderController struct {
	beego.Controller
}

func (this *OrderController) genError(code int, addr, msg string) string {

	orderErr := OrderErr{
		Code: code,
		Addr: addr,
		Msg:  msg,
	}
	ret, _ := json.Marshal(orderErr)
	return string(ret)
}

func (this *OrderController) Order() {
	invoiceId := this.GetString("invoice_id")
	//value := this.GetString("value")

	callBack := CALLBACK_ADDR + "notify?invoice_id=" + invoiceId + "&secret=" + SECRET
	urlCallBack := url.QueryEscape(callBack)

	fmt.Printf("callback=%v\r\n", callBack)
	fmt.Printf("callback111=%v\r\n", urlCallBack)

	url := BLOCKCHAIN_RECEIVE_ROOT + "v2/receive?" + "key=" + API_KEY + "&callback=" + urlCallBack + "&xpub=" + XPUB
	resp, err := http.Get(url)
	if err != nil {
		this.Ctx.Output.Body([]byte(this.genError(1, "", err.Error())))
		return
	}

	defer resp.Body.Close()

	if err != nil {
		this.Ctx.Output.Body([]byte(this.genError(2, "", err.Error())))
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		this.Ctx.Output.Body([]byte(this.genError(3, "", err.Error())))
		return
	}

	fmt.Printf("body=%v\r\n", string(body))

	orderRes := &OrderResponse{}
	err = json.Unmarshal(body, orderRes)
	if err != nil {
		this.Ctx.Output.Body([]byte(this.genError(4, "", err.Error())))
		return
	}
	if orderRes.Addr == "" {
		this.Ctx.Output.Body([]byte(this.genError(5, "", "addr is nil.")))
		return
	}
	this.Ctx.Output.Body([]byte(this.genError(0, orderRes.Addr, "success")))
}
