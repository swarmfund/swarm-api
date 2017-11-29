package skrill

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrWrongFormat = errors.New("wrong amount format")
)

func pow10(n int) (result int64) {
	result = 1
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}

func rightPad(s string, padWith string, totalLen int) string {
	if len(padWith) != 1 {
		panic("invalid input")
	}
	return s + strings.Repeat(padWith, totalLen-len(s))
}

func leftPad(s string, padWith string, totalLen int) string {
	if len(padWith) != 1 {
		panic("invalid input")
	}
	return strings.Repeat(padWith, totalLen-len(s)) + s
}

func ParseAmount(value string, precision int) (int64, error) {
	negative := strings.HasPrefix(value, "-")
	var result int64
	var err error
	split := strings.Split(value, ".")
	switch len(split) {
	case 1: // exponent only
		exponent, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return 0, err
		}
		result = exponent * pow10(precision)
	case 2: // mantis, skrill does not force double-digit precision
		var exponent int64 = 0
		if split[0] != "" {
			exponent, err = strconv.ParseInt(split[0], 10, 64)
			if err != nil {
				return 0, err
			}
		}
		if len(split[1]) > precision {
			split[1] = split[1][:precision]
		}
		mantis, err := strconv.ParseInt(rightPad(split[1], "0", precision), 10, 64)
		if err != nil {
			return 0, err
		}
		result = exponent*pow10(precision) + mantis
	default:
		return 0, ErrWrongFormat
	}
	if negative {
		result *= -1
	}
	return result, nil
}

func AmountToString(value int64, precision int) string {
	var abs string
	negative := value < 0
	if negative {
		value *= -1
	}
	if value < pow10(precision) {
		// mantis only
		mantis := strings.TrimRight(leftPad(fmt.Sprintf("%d", value), "0", precision), "0")
		if mantis == "" {
			return "0"
		}
		abs = fmt.Sprintf("0.%s", mantis)
	} else {
		wodot := fmt.Sprintf("%d", value)
		abs = strings.TrimRight(fmt.Sprintf("%s.%s", wodot[:len(wodot)-precision], wodot[len(wodot)-precision:]), ".0")
	}
	if negative {
		return fmt.Sprintf("-%s", abs)
	}
	return abs
}
