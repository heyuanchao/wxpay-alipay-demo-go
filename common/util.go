package common

import (
	"net/url"
	"bytes"
	"sort"
	"time"
	"fmt"
	"math/rand"
)


func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetSignContent(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		prefix := k + "="
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(v)
		}
	}
	return buf.String()
}

func GetOutTradeNo() string {
	return time.Now().Format("0102150405") + fmt.Sprintf("%05d", rand.Intn(100000))
}
