package main

import (
	"github.com/Malyue/quotaguard/pkg/handle"
	"github.com/Malyue/quotaguard/pkg/server"
	"k8s.io/klog/v2"
	"net/http"
)

func main() {

	s := server.NewServer()

	http.HandleFunc("/validate", handle.NewHandler(s))
	klog.Info("Start QuotaGuard Webhook...")
	if err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil); err != nil {
		klog.Fatalf("服务器启动失败: %v", err)
	}
}
