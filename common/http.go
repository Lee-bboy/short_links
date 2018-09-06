package common

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"
)

func ContentType(inner http.Handler, contentType string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		inner.ServeHTTP(w, r)
	})
}

func Auth(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appKey := r.FormValue("appKey")
		timestamp := r.FormValue("timestamp")
		signature := r.FormValue("signature")

		//获取用户秘钥
		appSecret, err := GetAppSecret(appKey)
		if err != nil {
			w.Write([]byte(`{"errcode": 403, "errmsg": "API认证失败"}`))
			return
		}

		//验证签名是否正确
		sign := md5.Sum([]byte(appKey + "+" + appSecret + "+" + timestamp))
		if hex.EncodeToString(sign[:]) != signature {
			w.Write([]byte(`{"errcode": 403, "errmsg": "API认证失败"}`))
			return
		}

		//请求是否超时，请求最大时间为5秒
		t, err := strconv.Atoi(timestamp)
		if err != nil || t+60 <= int(time.Now().Unix()) {
			w.Write([]byte(`{"errcode": 403, "errmsg": "API认证失败"}`))
			return
		}

		inner.ServeHTTP(w, r)
	})
}
