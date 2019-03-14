package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"github.com/heyuanchao/wxpay-alipay-demo-go/alipay"
	"github.com/heyuanchao/wxpay-alipay-demo-go/wxpay"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/alipay", http.HandlerFunc(handleAliPay))
	mux.Handle("/wxpay", http.HandlerFunc(handleWXPay))
	err := http.ListenAndServe("", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func handleAliPay(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		total_amount := r.URL.Query().Get("total_amount")
		if total_amount == "" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "%v", "no total_amount")
			return
		}
		request := alipay.NewAlipayTradeAppPayRequest(total_amount)
		data, err := alipay.DoRequest(request)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		}
		fmt.Fprintf(w, "%s", data)
	case "POST":
		result, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "%v", "failure")
			return
		}
		log.Printf("result: %s\n", result)
		m, err := url.ParseQuery(string(result))
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "%v", "failure")
			return
		}
		ok, err := alipay.Check(m)
		if err != nil {
			fmt.Fprintf(w, "%v", "failure")
			return
		}
		if ok {
			// 需要验证 out_trade_no 和 total_amount
			log.Println("success")
			fmt.Fprintf(w, "%v", "success")
		} else {
			fmt.Fprintf(w, "%v", "failure")
		}
	}
}

func handleWXPay(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		total_fee := r.URL.Query().Get("total_fee")
		if total_fee == "" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "%v", "no total_amount")
			return
		}
		ip := strings.Split(r.RemoteAddr, ":")[0]
		p := wxpay.NewWXTradeAppPayParameter(total_fee, ip)
		data, err := json.Marshal(p)
		if err != nil {
			log.Printf("marshal message %v error: %v\n", reflect.TypeOf(p), err)
			fmt.Fprintf(w, "%v", err)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, "%s", data)
	case "POST":
		result, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "%v", err)
			return
		}
		log.Printf("result: %s\n", result)
		payResult := new(wxpay.WXPayResult)
		err = xml.Unmarshal(result, &payResult)
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "%v", err)
			return
		}
		if wxpay.VerifyPayResult(payResult) {
			// 需要验证 out_trade_no 和 total_fee
			fmt.Fprintf(w, "%v", wxpay.ReturnWXSuccess)
		} else {
			fmt.Fprintf(w, "%v", wxpay.ReturnWXFail)
		}
	}
}
