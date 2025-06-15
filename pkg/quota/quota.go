package quota

import (
	"errors"
	"fmt"
	quotav1 "github.com/Malyue/quotaguard/pkg/apis/quota/v1"
	"sync"
)

const (
	TeamKey      = "Team"
	NamespaceKey = "Namespace"
)

type QuotaManager struct {
	sync.Mutex
	policys map[string]QuotaPolicy
}

func (m *QuotaManager) AddQuotaPolicy(policy *quotav1.QuotaPolicy) error {
	for _, rule := range policy.Spec.Rule {
		if rule.Target.Kind == TeamKey {
			m.AddPolicy(QuotaPolicy{
				Team:      rule.Target.Key,
				MaxCPU:    rule.Limit.CPU,
				MaxMemory: rule.Limit.Memory,
			})
		} else if rule.Target.Kind == NamespaceKey {
			m.AddPolicy(QuotaPolicy{
				Namespace: rule.Target.Key,
				MaxCPU:    rule.Limit.CPU,
				MaxMemory: rule.Limit.Memory,
			})
		} else {
			return fmt.Errorf("invalid target key")
		}
	}

	return nil
}

func (m *QuotaManager) DeleteQuotaPolicy(policy *quotav1.QuotaPolicy) error {
	for _, rule := range policy.Spec.Rule {
		if rule.Target.Kind != TeamKey && rule.Target.Kind != NamespaceKey {
			return fmt.Errorf("invalid target key")
		}

		if rule.Target.Kind == TeamKey {
			m.DeletePolicy(generateKey(TeamKey, rule.Target.Key))
		} else {
			m.DeletePolicy(generateKey(NamespaceKey, rule.Target.Key))
		}
	}

	return nil
}

func (m *QuotaManager) UpdateQuotaPolicy(policy *quotav1.QuotaPolicy) error {
	// TODO 如果 Target.Kind 发生了变化，需要删除旧的再创建新的
	return m.AddQuotaPolicy(policy)
}

func (m *QuotaManager) AddPolicy(rule QuotaPolicy) {
	m.Lock()
	defer m.Unlock()

	prefix := TeamKey
	value := rule.Team
	if rule.Team == "" && rule.Namespace != "" {
		prefix = NamespaceKey
		value = rule.Namespace
	}

	m.policys[prefix+value] = rule
}

func (m *QuotaManager) DeletePolicy(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.policys, key)
}

func NewQuotaManager() *QuotaManager {
	return &QuotaManager{
		policys: make(map[string]QuotaPolicy),
	}
}

func (m *QuotaManager) GetPolicy(key string) (QuotaPolicy, bool) {
	m.Lock()
	defer m.Unlock()
	rule, exists := m.policys[key]
	return rule, exists
}

func (m *QuotaManager) Validate(team, namespace string, cpu, memory string) (bool, error) {
	if team != "" {
		teamPolicy, exists := m.GetPolicy(generateKey(TeamKey, team))
		if exists {
			valid, err := teamPolicy.Valid(team, cpu, memory)
			if err != nil {
				return false, err
			}
			if !valid {
				return false, errors.New("the pod's resources is over the team quota policy")
			}
		}
	}

	if namespace != "" {
		nsPolicy, exists := m.GetPolicy(generateKey(NamespaceKey, namespace))
		if !exists {
			return true, nil
		}

		valid, err := nsPolicy.Valid(namespace, cpu, memory)
		if err != nil {
			return false, err
		}
		if !valid {
			return false, errors.New("the pod's resources is over the namespace quota policy")
		}
	}

	return true, nil
}

func generateKey(prefix, key string) string {
	return prefix + "-" + key
}
