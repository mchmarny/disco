package types

const (
	TargetTypeUndefined TargetType = iota
	TargetTypeBigQuery
)

type TargetType int64

func (t TargetType) String() string {
	switch t {
	case TargetTypeBigQuery:
		return "bigquery"
	default:
		return "undefined"
	}
}

func ParseTargetTypeOrDefault(s string) TargetType {
	switch s {
	case "bigquery":
		return TargetTypeBigQuery
	default:
		return TargetTypeUndefined
	}
}
