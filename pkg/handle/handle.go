package handle

import (
	"encoding/json"
	"github.com/Malyue/quotaguard/pkg/server"
	"k8s.io/apimachinery/pkg/api/resource"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

const (
	AdmissionApiVersion = "admission.k8s.io/v1"
	AdmissionReview     = "AdmissionReview"
)

func NewHandler(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var admissionReview admissionv1.AdmissionReview

		if err := json.NewDecoder(r.Body).Decode(&admissionReview); err != nil {
			klog.Errorf("cannot decode request: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response := admissionv1.AdmissionReview{
			TypeMeta: metav1.TypeMeta{
				APIVersion: AdmissionApiVersion,
				Kind:       AdmissionReview,
			},
			Response: &admissionv1.AdmissionResponse{
				UID: admissionReview.Request.UID,
			},
		}

		var pod corev1.Pod
		if err := json.Unmarshal(admissionReview.Request.Object.Raw, &pod); err != nil {
			klog.Errorf("cannot unmarshal pod info: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		cpu, memory := calculateResources(&pod)
		team := pod.Labels["team"]
		namespace := pod.Namespace

		allowed, err := s.QuotaManager.Validate(team, namespace, cpu, memory)
		if err != nil {
			klog.Errorf("cannot validate quota: %v", err)
		}

		response.Response.Allowed = allowed

		if !allowed {
			response.Response.Result = &metav1.Status{
				Message: err.Error(),
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			klog.Errorf("cannot encode response: %v", err)
		}
	}
}

func calculateResources(pod *corev1.Pod) (totalCPU, totalMemory string) {
	cpu := resource.NewQuantity(0, resource.DecimalSI)
	memory := resource.NewQuantity(0, resource.BinarySI)

	for _, container := range pod.Spec.Containers {
		if container.Resources.Requests != nil {
			if containerCPU, ok := container.Resources.Requests[corev1.ResourceCPU]; ok {
				cpu.Add(containerCPU.DeepCopy()) // 直接操作指针
			}
			if containerMem, ok := container.Resources.Requests[corev1.ResourceMemory]; ok {
				memory.Add(containerMem.DeepCopy())
			}
		}
	}

	return cpu.String(), memory.String()
}
