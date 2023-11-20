package monitor

var (
	CheckStatusEnqueued = CheckStatus{"enqueued"}
	CheckStatusPending  = CheckStatus{"pending"}
	CheckStatusChecked  = CheckStatus{"checked"}
)

type CheckStatus struct {
	value string
}

func (c CheckStatus) String() string {
	return c.value
}
