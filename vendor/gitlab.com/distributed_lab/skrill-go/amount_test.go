package skrill

import "testing"

func TestParseAmount(t *testing.T) {
	type Input struct {
		value     string
		precision int
		result    int64
		error     bool
	}
	inputs := []Input{
		{"1", 4, 10000, false},
		{"111", 4, 1110000, false},
		{"1,1", 4, 0, true},
		{"1.0", 4, 10000, false},
		{"0.0001", 4, 1, false},
		{"10.10", 4, 101000, false},
		{"0.000000001", 4, 0, false},
		{"10.1111", 4, 101111, false},
		{".326758", 4, 3267, false},
		{"-0.1", 4, -1000, false},
	}
	for i, input := range inputs {
		result, err := ParseAmount(input.value, input.precision)
		if input.error && err == nil {
			t.Errorf("%d should return error", i)
			continue
		}
		if result != input.result {
			t.Errorf("%d returned %d expected %d", i, result, input.result)
			continue
		}
	}
}

func TestAmountToString(t *testing.T) {
	type Input struct {
		value     int64
		precision int
		result    string
	}
	inputs := []Input{
		{0, 4, "0"},
		{1, 4, "0.0001"},
		{3345, 4, "0.3345"},
		{11112, 4, "1.1112"},
		{10000, 4, "1"},
		{-1000, 4, "-0.1"},
	}
	for i, input := range inputs {
		result := AmountToString(input.value, input.precision)
		if result != input.result {
			t.Errorf("%d returned %s expected %s", i, result, input.result)
			continue
		}
	}
}
