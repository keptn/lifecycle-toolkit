package converter

import (
	"fmt"
	"math"

	"gopkg.in/inf.v0"
)

const InvalidOperatorErrMsg = "invalid operator: '%s'"
const UnableConvertValueErrMsg = "unable to convert value '%s' to decimal"
const UnsupportedIntervalCombinationErrMsg = "unsupported interval combination '%v'"
const EmptyOperatorsErrMsg = "empty operators: '%v'"
const UnconvertableOperatorsCombinationErrMsg = "unconvertable combination of operators: '%s', '%s'"

func NewInvalidOperatorErr(msg string) error {
	return fmt.Errorf(InvalidOperatorErrMsg, msg)
}

func NewUnconvertableValueErr(msg string) error {
	return fmt.Errorf(UnableConvertValueErrMsg, msg)
}

func NewUnsupportedIntervalCombinationErr(op []string) error {
	return fmt.Errorf(UnsupportedIntervalCombinationErrMsg, op)
}

func NewEmptyOperatorErr(op []string) error {
	return fmt.Errorf(UnsupportedIntervalCombinationErrMsg, op)
}

func NewUnconvertableOperatorCombinationErr(op1, op2 string) error {
	return fmt.Errorf(UnconvertableOperatorsCombinationErrMsg, op1, op2)
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
