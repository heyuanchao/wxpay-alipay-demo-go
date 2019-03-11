package wxpay

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/url"
	"time"
	"../common"
)

var (
	appid      = "wxd678efh567hg6787"
	body       = "游戏充值"
	key        = "E10ADC3949BA59ABBE56E057F20F883E"
	mch_id     = "1230000109"
	notify_url = "http://xxx.xxx.xxx.xxx/wxpay"

	ReturnWXSuccess = "<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg></xml>"
	ReturnWXFail    = "<xml><return_code><![CDATA[FAIL]]></return_code><return_msg><![CDATA[签名失败/参数格式校验错误]]></return_msg></xml>"
)

type WXPayResult struct {
	AppID         string `xml:"appid"`
	BankType      string `xml:"bank_type"` // 付款银行
	CashFee       string `xml:"cash_fee"`  // 现金支付金额
	FeeType       string `xml:"fee_type"`  // 货币种类
	IsSubscribe   string `xml:"is_subscribe"`
	MchID         string `xml:"mch_id"`       // 商户号
	NonceStr      string `xml:"nonce_str"`    // 随机字符串
	OpenID        string `xml:"openid"`       // 用户标识
	OutTradeNo    string `xml:"out_trade_no"` // 商户订单号
	ResultCode    string `xml:"result_code"`
	ReturnCode    string `xml:"return_code"`
	Sign          string `xml:"sign"`           // 签名
	TimeEnd       string `xml:"time_end"`       // 支付完成时间
	TotalFee      string `xml:"total_fee"`      // 总金额
	TradeType     string `xml:"trade_type"`     // 交易类型
	TransactionID string `xml:"transaction_id"` // 微信支付订单号
}

func NewWXTradeAppPayParameter(total_fee, ip string) map[string]string {
	m := map[string]string{}
	m["appid"] = appid
	m["body"] = body
	m["key"] = key
	m["mch_id"] = mch_id
	m["notify_url"] = notify_url
	m["out_trade_no"] = common.GetOutTradeNo()
	m["total_fee"] = total_fee
	m["spbill_create_ip"] = ip
	return m
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getSignParams(r *WXPayResult) url.Values {
	p := url.Values{}
	p.Add("appid", r.AppID)
	p.Add("bank_type", r.BankType)
	p.Add("cash_fee", r.CashFee)
	p.Add("fee_type", r.FeeType)
	p.Add("is_subscribe", r.IsSubscribe)
	p.Add("mch_id", r.MchID)
	p.Add("nonce_str", r.NonceStr)
	p.Add("openid", r.OpenID)
	p.Add("out_trade_no", r.OutTradeNo)
	p.Add("result_code", r.ResultCode)
	p.Add("return_code", r.ReturnCode)
	p.Add("time_end", r.TimeEnd)
	p.Add("total_fee", r.TotalFee)
	p.Add("trade_type", r.TradeType)
	p.Add("transaction_id", r.TransactionID)
	return p
}

func VerifyPayResult(result *WXPayResult) bool {
	if result.ReturnCode == "SUCCESS" && verify(generateSign(getSignParams(result)), result.Sign) {
		return true
	}
	return false
}

func generateSign(params url.Values) string {
	return sign(common.GetSignContent(params) + "&key=" + key)
}

func sign(data string) string {
	log.Println(data)
	m := md5.New()
	io.WriteString(m, data)
	return fmt.Sprintf("%X", m.Sum(nil))
}

func verify(data string, sign string) bool {
	log.Println(data)
	return data == sign
}
