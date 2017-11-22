package alipay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"wxpay-alipay-demo-go/common"
)

var (
	rsaPrivateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
XXX
-----END RSA PRIVATE KEY-----`)

	alipayRSAPublicKey = []byte(`-----BEGIN PUBLIC KEY-----
XX
-----END PUBLIC KEY-----`)

	partnerID  = "2088111111116894"
	appID      = "2014072300007148"
	gatewayUrl = "https://openapi.alipay.com/gateway.do"
	notifyUrl  = "http://xxx.xxx.xxx.xxx/alipay"
)

func DoRequest(req *AlipayTradeAppPayRequest) []byte {
	p := url.Values{}
	p.Add("app_id", appID)
	p.Add("biz_content", req.BizContent)
	p.Add("charset", "utf-8")
	p.Add("method", req.Method)
	p.Add("notify_url", req.NotifyUrl)
	p.Add("sign_type", "RSA2")
	p.Add("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	p.Add("version", req.Version)
	p.Add("sign", generateSign(p))

	r, err := http.NewRequest("POST", gatewayUrl, strings.NewReader(p.Encode()))
	if err != nil {
		log.Print(err)
		return []byte{}
	}
	defer r.Body.Close()
	result, _ := ioutil.ReadAll(r.Body)
	return result
}

func rsaCheck(params url.Values) bool {
	sign := params.Get("sign")
	params.Del("sign")
	params.Del("sign_type")
	return verify([]byte(common.GetSignContent(params)), sign)
}

func Check(params url.Values) bool {
	tradeStatus := params.Get("trade_status")
	if appID == params.Get("app_id") && partnerID == params.Get("seller_id") && (tradeStatus == "TRADE_SUCCESS" || tradeStatus == "TRADE_FINISHED") {
		return rsaCheck(params)
	}
	return false
}

func generateSign(params url.Values) string {
	return sign([]byte(common.GetSignContent(params)))
}

func sign(data []byte) string {
	block, _ := pem.Decode(rsaPrivateKey)
	if block == nil {
		return ""
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return ""
	}
	h := sha256.New()
	h.Write(data)
	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, h.Sum(nil))
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(sign)
}

func verify(data []byte, sign string) bool {
	sig, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false
	}
	block, _ := pem.Decode(alipayRSAPublicKey)
	if block == nil {
		return false
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false
	}
	h := sha256.New()
	h.Write(data)
	err = rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA256, h.Sum(nil), sig)
	if err == nil {
		return true
	}
	return false
}
