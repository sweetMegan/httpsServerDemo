package main

import (
	"crypto/x509"
	"io/ioutil"
	"fmt"
	"crypto/tls"
	"net/http"
)

func main() {
	//单向认证
	//verifySingle()
	//双向认证
	verifyEachOther()
}
func verifySingle() {
	//跳过验证
	tr := &http.Transport{
		TLSClientConfig:&tls.Config{InsecureSkipVerify:true},
	}
	client := &http.Client{Transport:tr}
	resp,err := client.Get("https://localhost:8081/test")
	//resp, err := http.Get("https://localhost:8081/test")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

//双向认证
func verifyEachOther()  {
	pool := x509.NewCertPool()
	caCertPath := "./ca.crt"
	caCrt,err := ioutil.ReadFile(caCertPath)
	if err != nil {
		fmt.Println("ReadFile err: ",err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)
	cliCrt,err := tls.LoadX509KeyPair("client/client.crt","client/client.key")
	if err != nil {
		fmt.Println("LoadX509keypair err:",err)
		return
	}
	tr := &http.Transport{
		TLSClientConfig:&tls.Config{
			RootCAs:pool,
			Certificates:[]tls.Certificate{cliCrt},
		},
	}
	client := &http.Client{
		Transport:tr,
	}
	resp,err := client.Get("https://localhost:8081/test")
	if err != nil {
		fmt.Println("Get error:",err)
		return
	}
	defer resp.Body.Close()
	body,err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}