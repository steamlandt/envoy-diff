package process

import (
	"os"
	"testing"
)

func TestParseStrings_Basic(t *testing.T) {
	input := []string{"FOO=bar", "BAZ=qux"}
	got := parseStrings(input)
	if got["FOO"] != "bar" {
		t.Errorf("FOO: want %q, got %q", "bar", got["FOO"])
	}
	if got["BAZ"] != "qux" {
		t.Errorf("BAZ: want %q, got %q", "qux", got["BAZ"])
	}
}

func TestParseStrings_EmptyValue(t *testing.T) {
	got := parseStrings([]string{"EMPTY="})
	if v, ok := got["EMPTY"]; !ok || v != "" {
		t.Errorf("EMPTY: want empty string, got %q (ok=%v)", v, ok)
	}
}

func TestParseStrings_NoEquals(t *testing.T) {
	got := parseStrings([]string{"STANDALONE"})
	if v, ok := got["STANDALONE"]; !ok || v != "" {
		t.Errorf("STANDALONE: want empty string, got %q (ok=%v)", v, ok)
	}
}

func TestParseStrings_ValueWithEquals(t *testing.T) {
	got := parseStrings([]string{"URL=http://example.com?a=1&b=2"})
	want := "http://example.com?a=1&b=2"
	if got["URL"] != want {
		t.Errorf("URL: want %q, got %q", want, got["URL"])
	}
}

func TestParseStrings_SkipsBlanks(t *testing.T) {
	got := parseStrings([]string{"A=1", "", "  ", "B=2"})
	if len(got) != 2 {
		t.Errorf("expected 2 entries, got %d", len(got))
	}
}

func TestParseEnviron_NulDelimited(t *testing.T) {
	data := []byte("KEY1=hello\x00KEY2=world\x00")
	got := parseEnviron(data)
	if got["KEY1"] != "hello" {
		t.Errorf("KEY1: want %q, got %q", "hello", got["KEY1"])
	}
	if got["KEY2"] != "world" {
		t.Errorf("KEY2: want %q, got %q", "world", got["KEY2"])
	}
}

func TestReadSelf_ContainsKnownVar(t *testing.T) {
	const key = "_ENVOY_DIFF_TEST_VAR"
	const val = "hello_test"
	t.Setenv(key, val)

	env := ReadSelf()
	if env[key] != val {
		t.Errorf("%s: want %q, got %q", key, val, env[key])
	}
}

func TestParsePID_Valid(t *testing.T) {
	pid, err := ParsePID("1234")
	if err != nil || pid != 1234 {
		t.Errorf("want 1234, nil; got %d, %v", pid, err)
	}
}

func TestParsePID_Invalid(t *testing.T) {
	for _, bad := range []string{"0", "-1", "abc", ""} {
		_, err := ParsePID(bad)
		if err == nil {
			t.Errorf("ParsePID(%q): expected error, got nil", bad)
		}
	}
}

func TestReadPID_NonExistent(t *testing.T) {
	if os.Getenv("CI") == "" {
		t.Skip("skipping /proc test outside CI")
	}
	_, err := ReadPID(999999999)
	if err == nil {
		t.Error("expected error for non-existent PID, got nil")
	}
}
