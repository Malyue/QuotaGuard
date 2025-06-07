package quota

import (
	"errors"
	"sync"
)

type QuotaManager struct {
	sync.Mutex
	policys map[string]QuotaPolicy
}

func (m *QuotaManager) AddPolicy(rule QuotaPolicy) {
	m.Lock()
	defer m.Unlock()
	m.policys[rule.Team] = rule
}

func (m *QuotaManager) DeletePolicy(rule QuotaPolicy) {
	m.Lock()
	defer m.Unlock()
	delete(m.policys, rule.Team)
}

func NewQuotaManager() *QuotaManager {
	return &QuotaManager{
		policys: make(map[string]QuotaPolicy),
	}
}

func (m *QuotaManager) GetPolicy(team string) (QuotaPolicy, bool) {
	m.Lock()
	defer m.Unlock()
	rule, exists := m.policys[team]
	return rule, exists
}

func (m *QuotaManager) Validate(team, namespace string, cpu, memory string) (bool, error) {
	if team != "" {
		teamPolicy, exists := m.GetPolicy(team)
		if exists {
			valid := teamPolicy.ValidTeam(team, cpu, memory)
			if !valid {
				return false, errors.New("the pod's resources is over the team quota policy")
			}
		}
	}

	if namespace != "" {
		nsPolicy, exists := m.GetPolicy(namespace)
		if !exists {
			return true, nil
		}

		valid := nsPolicy.ValidNamespace(namespace, cpu, memory)
		if !valid {
			return false, errors.New("the pod's resources is over the namespace quota policy")
		}
	}

	return true, nil
}
