// Code generated by "stringer -linecomment -type numberOrderResult"; DO NOT EDIT.

package types

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[doubleNegativeZero-1]
	_ = x[doubleDT-2]
	_ = x[int32DT-3]
	_ = x[int64DT-4]
}

const _numberOrderResult_name = "doubleNegativeZerodoubleDTint32DTint64DT"

var _numberOrderResult_index = [...]uint8{0, 18, 26, 33, 40}

func (i numberOrderResult) String() string {
	i -= 1
	if i >= numberOrderResult(len(_numberOrderResult_index)-1) {
		return "numberOrderResult(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _numberOrderResult_name[_numberOrderResult_index[i]:_numberOrderResult_index[i+1]]
}
