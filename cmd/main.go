package main

import (
	"crypto/tls"
	"github.com/Malyue/quotaguard/pkg/handle"
	"github.com/Malyue/quotaguard/pkg/server"
	"k8s.io/klog/v2"
	"net/http"
)

func main() {

	s, err := server.NewServer("")
	if err != nil {
		klog.Fatalf("server init err: %v", err)
	}

	cert, err := tls.LoadX509KeyPair(
		"/etc/webhook/certs/tls.crt",
		"/etc/webhook/certs/tls.key",
	)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}

	http.HandleFunc("/validate", handle.ValidHandler(s))
	http.HandleFunc("/get", handle.Get(s))
	http.HandleFunc("/all", handle.All(s))
	klog.Info("Start QuotaGuard Webhook...")
	if err := server.ListenAndServeTLS("", ""); err != nil {
		klog.Fatalf("服务器启动失败: %v", err)
	}
}
