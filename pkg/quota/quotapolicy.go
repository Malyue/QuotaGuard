package quota

import "k8s.io/apimachinery/pkg/api/resource"

type QuotaPolicy struct {
	Team      string
	Namespace string
	MaxCPU    string
	MaxMemory string
	MaxGPU    string
}

func parseCPU(cpuStr string) (float64, error) {
	// 使用 Kubernetes 的 resource 包解析
	q, err := resource.ParseQuantity(cpuStr)
	if err != nil {
		return 0, err
	}
	// 转换为毫核后计算核数
	return float64(q.MilliValue()) / 1000, nil
}

func parseMemory(memoryStr string) (int64, error) {
	q, err := resource.ParseQuantity(memoryStr)
	if err != nil {
		return 0, err
	}
	return q.Value(), nil // 转换为字节（如 "1Gi" -> 1073741824 bytes）
}

func (qp *QuotaPolicy) Valid(key string, cpu, memory string) (bool, error) {
	if qp.Team != key && qp.Namespace != key {
		return false, nil
	}

	limitCPU, err := parseCPU(qp.MaxCPU)
	if err != nil {
		return false, err
	}

	reqCPU, err := parseCPU(cpu)
	if err != nil {
		return false, err
	}

	limitMemory, err := parseMemory(qp.MaxMemory)
	if err != nil {
		return false, err
	}

	reqMemory, err := parseMemory(memory)
	if err != nil {
		return false, err
	}

	return limitCPU <= reqCPU && limitMemory <= reqMemory, nil
}
