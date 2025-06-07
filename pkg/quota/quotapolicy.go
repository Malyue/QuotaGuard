package quota

type QuotaPolicy struct {
	Team      string
	Namespace string
	MaxCPU    string
	MaxMemory string
	MaxGPU    string
}

func (qp *QuotaPolicy) ValidTeam(team string, cpu, memory string) bool {
	if qp.Team != team {
		return false
	}

	// TODO how to compare cpu and memory

	return true
}

func (qp *QuotaPolicy) ValidNamespace(namespace string, cpu, memory string) bool {
	if qp.Namespace != namespace {
		return false
	}

	// TODO how to compare cpu and memory

	return true
}
