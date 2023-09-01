package converter

import (
	"fmt"
	"math"

	"gopkg.in/inf.v0"
)

func NewInvalidOperatorErr(msg string) error {
	return fmt.Errorf("invalid operator: '%s'", msg)
}

func NewUnconvertableValueErr(msg string) error {
	return fmt.Errorf("unable to convert value '%s' to decimal", msg)
}

func NewUnsupportedIntervalCombinationErr(op []string) error {
	return fmt.Errorf("unsupported interval combination '%v'", op)
}

func NewEmptyOperatorErr(op []string) error {
	return fmt.Errorf("empty operators: '%v'", op)
}

func NewUnconvertableOperatorCombinationErr(op1, op2 string) error {
	return fmt.Errorf("unconvertable combination of operators: '%s', '%s'", op1, op2)
}

const MaxInt = math.MaxInt
const MinInt = -MaxInt - 1

type Operator struct {
	Value     *inf.Dec
	Operation string
}

type Interval struct {
	Start *inf.Dec
	End   *inf.Dec
}

func isGreaterOrEqual(op string) bool {
	return op == ">" || op == ">="
}

func isLessOrEqual(op string) bool {
	return op == "<" || op == "<="
}
