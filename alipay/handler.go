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
	"net/http"
	"net/url"
	"strings"
	"time"
	"github.com/heyuanchao/wxpay-alipay-demo-go/common"
	"log"
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

func DoRequest(req *AlipayTradeAppPayRequest) ([]byte, error) {
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
		log.Println(err)
		return []byte{}, err
	}
	result, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println(err)
	}
	return result, err
}

func rsaCheck(params url.Values) (bool, error) {
	sign := params.Get("sign")
	params.Del("sign")
	params.Del("sign_type")
	return verify([]byte(common.GetSignContent(params)), sign)
}

func Check(params url.Values) (bool, error) {
	tradeStatus := params.Get("trade_status")
	if appID == params.Get("app_id") && partnerID == params.Get("seller_id") && (tradeStatus == "TRADE_SUCCESS" || tradeStatus == "TRADE_FINISHED") {
		return rsaCheck(params)
	}
	return false, nil
}

func generateSign(params url.Values) string {
	s, _ := sign([]byte(common.GetSignContent(params)))
	return s
}

func sign(data []byte) (string, error) {
	block, _ := pem.Decode(rsaPrivateKey)
	if block == nil {
		return "", nil
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Println(err)
		return "", err
	}
	h := sha256.New()
	_, err = h.Write(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, h.Sum(nil))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sign), nil
}

func verify(data []byte, sign string) (bool, error) {
	sig, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		log.Println(err)
		return false, err
	}
	block, _ := pem.Decode(alipayRSAPublicKey)
	if block == nil {
		log.Println(err)
		return false, nil
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Println(err)
		return false, err
	}
	h := sha256.New()
	_, err = h.Write(data)
	if err != nil {
		log.Println(err)
		return false, err
	}
	err = rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA256, h.Sum(nil), sig)
	if err != nil {
		log.Println(err)
	}
	return err == nil , err
}
