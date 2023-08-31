package converter

import "gopkg.in/inf.v0"

const InvalidOperatorErrMsg = "invalid operator: '%s'"
const UnableConvertValueErrMsg = "unable to convert value '%s' to decimal"
const UnsupportedIntervalCombinationErrMsg = "unsupported interval combination '%s'"
const EmptyOperatorsErrMsg = "empty operators: '%v'"
const UnconvertableOperatorsCombinationErrMsg = "unconvertable combination of operators: '%s', '%s'"

const MaxInt = int(^uint(0) >> 1)
const MinInt = -MaxInt - 1

type Operator struct {
	Value     *inf.Dec
	Operation string
}

type Interval struct {
	Start *inf.Dec
	End   *inf.Dec
}
