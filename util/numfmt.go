package util

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Unit struct {
	Symbol string
	OoM    int
}

// https://egg-inc.fandom.com/wiki/Order_of_Magnitude
var Units = []Unit{
	{"M", 6},
	{"B", 9},
	{"T", 12},
	{"q", 15},
	{"Q", 18},
	{"s", 21},
	{"S", 24},
	{"o", 27},
	{"N", 30},
	{"d", 33},
	{"U", 36},
	{"D", 39},
	{"Td", 42},
	{"qd", 45},
	{"Qd", 48},
	{"sd", 51},
	{"Sd", 54},
	{"Od", 57},
	{"Nd", 60},
	{"V", 63},
	{"uV", 66},
	{"dV", 69},
	{"tV", 72},
	{"qV", 75},
	{"QV", 78},
	{"sV", 81},
	{"SV", 84},
	{"OV", 87},
	{"NV", 90},
	{"tT", 93},
}

var (
	_oomMap map[int]string
	_minOoM int
	_maxOoM int
)

func init() {
	_oomMap = make(map[int]string)
	for _, u := range Units {
		_oomMap[u.OoM] = u.Symbol
	}
	_minOoM = Units[0].OoM
	_maxOoM = Units[len(Units)-1].OoM
}

func Numfmt(x float64) string {
	return numfmt(x, 3, false)
}

func NumfmtWhole(x float64) string {
	return numfmt(x, 2, true)
}

func numfmt(x float64, decimalDigits uint, trimTrailingZeros bool) string {
	if x < 0 {
		return "-" + numfmt(-x, decimalDigits, trimTrailingZeros)
	}
	if x == 0 {
		return "0"
	}
	oom := math.Log10(x)
	if oom < float64(_minOoM) {
		// Always round small number to an integer.
		return fmt.Sprintf("%.0f", x)
	}
	oomFloor := int(oom)
	if oom+1e-9 >= float64(oomFloor+1) {
		// Fix problem of 1q being displayed as 1000T, 1N displayed as 1000o, etc,
		// where the floor is one integer down due to floating point imprecision.
		oomFloor += 1
	}
	oomFloor -= oomFloor % 3
	if oomFloor > _maxOoM {
		oomFloor = _maxOoM
	}
	principal := x / math.Pow10(oomFloor)
	numpart := fmt.Sprintf("%."+strconv.Itoa(int(decimalDigits))+"f", principal)
	if trimTrailingZeros && strings.Contains(numpart, ".") {
		numpart = strings.TrimRight(numpart, "0")
		numpart = strings.TrimRight(numpart, ".")
	}
	return numpart + _oomMap[oomFloor]
}

func FormatPercentage(x float64, decimalDigits uint) string {
	if x <= 0 {
		return "0%"
	}
	if x >= 100 {
		return "100%"
	}
	return fmt.Sprintf("%."+strconv.Itoa(int(decimalDigits))+"f%%", x)
}
