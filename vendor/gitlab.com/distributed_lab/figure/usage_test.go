package figure_test

import (
	"reflect"
	"testing"
	"time"

	"math/big"

	"gitlab.com/distributed_lab/figure"
)

func TestSimpleUsage(t *testing.T) {
	type Config struct {
		SomeInt                int
		SomeString             string
		Missing                int
		Default                int
		DurationStr            time.Duration
		DurationInt            time.Duration
		DurationPointer        *time.Duration
		DurationPointerMissing *time.Duration
		BigInt                 *big.Int
		BigStr                 *big.Int
	}
	c := Config{
		Default: 42,
	}
	err := figure.Out(&c).From(map[string]interface{}{
		"some_int":         1,
		"some_string":      "satoshi",
		"duration_str":     "1s",
		"duration_int":     1,
		"duration_pointer": "1h",
		"big_int":          42,
		"big_str":          "42",
	}).Please()
	if err != nil {
		t.Fatalf("expected nil error got %s", err)
	}
	durationPointer := 1 * time.Hour
	reference := Config{
		SomeInt:         1,
		SomeString:      "satoshi",
		Default:         42,
		DurationStr:     1 * time.Second,
		DurationInt:     1,
		DurationPointer: &durationPointer,
		BigInt:          big.NewInt(42),
		BigStr:          big.NewInt(42),
	}
	if !reflect.DeepEqual(c, reference) {
		t.Errorf("expected %#v got %#v", reference, c)
	}

}
