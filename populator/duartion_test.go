package populator

import (
	"testing"
	"time"
)

type Config struct {
	Duration time.Duration `krakend:"duration,default=1m,required"`
}

func TestPopulate_TimeDuration_StringValue(t *testing.T) {
	input := map[string]any{
		"duration": "5s",
	}

	var cfg Config
	err := Populate(input, &cfg)
	if err != nil {
		t.Fatalf("unexpected populate error: %v", err)
	}

	expected := 5 * time.Second
	if cfg.Duration != expected {
		t.Errorf("wants %v, got %v", expected, cfg.Duration)
	}
}

func TestPopulate_TimeDuration_NumericValue(t *testing.T) {
	input := map[string]any{
		"duration": int64(5 * 1e9),
	}

	var cfg Config
	err := Populate(input, &cfg)
	if err != nil {
		t.Fatalf("unexpected populate error: %v", err)
	}

	expected := 5 * time.Second
	if cfg.Duration != expected {
		t.Errorf("wants %v, got %v", expected, cfg.Duration)
	}
}
