package main

import (
	"net/http"
	"fmt"
	"crypto/x509"
	"io/ioutil"
	"crypto/tls"
	"github.com/kataras/iris"
)

type myhandler struct {
}

func (h *myhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL, "method:", r.Method)
	//在此绑定路由
	switch r.URL.String() {
	case "/test":
		{
			test(w, r)
		}
	default:
		{
			fmt.Println("hello world")
			fmt.Fprintln(w, "hello world yy")
		}
	}
}
func main() {
	irisServer()
	//httpsSingle()
	//httpsServer()
}
//单向认证
func httpsSingle()  {
	//单向认证
	http.HandleFunc("/test",test)
	http.ListenAndServeTLS(":8081","server.crt","server.key",nil)
}
//双向认证
func httpsServer() {
	//创建一个CertPool，用来读取根证书CA
	pool := x509.NewCertPool()
	//ca证书路径
	caCertPath := "ca.crt"
	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)
	s := &http.Server{
		Addr:    ":8081",
		Handler: &myhandler{},
		TLSConfig: &tls.Config{
			//证书签发者
			ClientCAs:  pool,
			//认证方式，双向认证
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}
	//开启一个Https服务
	/*
	certFile: 服务端证书路径
	keyFile:服务端私钥路径
	*/
	err = s.ListenAndServeTLS("server.crt", "server.key")
	if err != nil {
		fmt.Println("ListenAndServeTLS err:", err)
	}
}
func irisServer() {
	app := iris.Default()
	app.Configure()
	irisTest(app)
	app.Run(iris.TLS(":8081", "./server.crt", "./server.key"))
}
func test(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello world5555")
	w.Write([]byte("hello world"))
}
func irisTest(app *iris.Application) {
	app.Handle("GET", "/test", func(ctx iris.Context) {
		fmt.Println("iris https")
		ctx.Write([]byte("iris https"))
	})
}
