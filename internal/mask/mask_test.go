package mask_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/mask"
)

func TestValue_EmptyString(t *testing.T) {
	opts := mask.DefaultOptions()
	if got := mask.Value("", opts); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestValue_ShortValue_MaskedEntirely(t *testing.T) {
	opts := mask.DefaultOptions() // MinLength = 8
	v := "abc" // len 3 < 8
	got := mask.Value(v, opts)
	if got != "***" {
		t.Errorf("expected %q, got %q", "***", got)
	}
}

func TestValue_LongValue_ShowsLeading(t *testing.T) {
	opts := mask.DefaultOptions() // VisibleChars = 4
	v := "supersecretvalue"
	got := mask.Value(v, opts)
	expected := "supe" + "************"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestValue_ExactMinLength(t *testing.T) {
	opts := mask.DefaultOptions() // MinLength = 8, VisibleChars = 4
	v := "12345678" // len == MinLength
	got := mask.Value(v, opts)
	expected := "1234****"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestValue_CustomMaskChar(t *testing.T) {
	opts := mask.Options{VisibleChars: 2, MaskChar: '#', MinLength: 4}
	v := "password"
	got := mask.Value(v, opts)
	expected := "pa######"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestValue_VisibleCharsExceedsLength(t *testing.T) {
	opts := mask.Options{VisibleChars: 20, MaskChar: '*', MinLength: 1}
	v := "short"
	got := mask.Value(v, opts)
	// visible capped to len-1 = 4, last char masked
	expected := "shor*"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestMap_DoesNotMutateInput(t *testing.T) {
	opts := mask.DefaultOptions()
	input := map[string]string{
		"API_KEY": "supersecretvalue",
		"DB_PASS": "anotherpassword1",
	}
	_ = mask.Map(input, opts)
	if input["API_KEY"] != "supersecretvalue" {
		t.Error("Map mutated the input map")
	}
}

func TestMap_AllValuesObscured(t *testing.T) {
	opts := mask.DefaultOptions()
	input := map[string]string{
		"TOKEN": "abcdefghij",
		"SECRET": "xyz123456789",
	}
	out := mask.Map(input, opts)
	for k, v := range out {
		if v == input[k] {
			t.Errorf("key %s: value was not masked", k)
		}
	}
}
