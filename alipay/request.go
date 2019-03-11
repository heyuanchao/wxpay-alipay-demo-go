package alipay

import "../common"

type AlipayTradeAppPayRequest struct {
	BizContent string
	Method     string // alipay.trade.app.pay
	NotifyUrl  string
	Version    string // 1.0
}

func NewAlipayTradeAppPayRequest(total_amount string) *AlipayTradeAppPayRequest {
	bizContent := `{"product_code":"QUICK_MSECURITY_PAY","total_amount":"` + total_amount + `","subject":"游戏充值","out_trade_no":"` + common.GetOutTradeNo() + `"}`
	req := new(AlipayTradeAppPayRequest)
	req.BizContent = bizContent
	req.Method = "alipay.trade.app.pay"
	req.NotifyUrl = notifyUrl
	req.Version = "1.0"
	return req
}
