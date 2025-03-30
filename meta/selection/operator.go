package selection

type Operator string

const (
	In           Operator = "in"
	NotIn        Operator = "notin"
	Exists       Operator = "exists"
	DoesNotExist Operator = "!"
	Equals       Operator = "="
	DoubleEquals Operator = "=="
	NotEquals    Operator = "!="
	LessThan     Operator = "lt"
	GreaterThan  Operator = "gt"
)
