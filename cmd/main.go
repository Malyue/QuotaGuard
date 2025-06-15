package main

import (
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

	http.HandleFunc("/validate", handle.ValidHandler(s))
	klog.Info("Start QuotaGuard Webhook...")
	if err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil); err != nil {
		klog.Fatalf("服务器启动失败: %v", err)
	}
}
