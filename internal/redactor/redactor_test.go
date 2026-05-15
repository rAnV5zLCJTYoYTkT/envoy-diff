package redactor_test

import (
	"encoding/json"
	"testing"

	"github.com/example/envoy-diff/internal/differ"
	"github.com/example/envoy-diff/internal/redactor"
)

func makeResult(before, after string) differ.Result {
	return differ.Result{
		Type:   "listener",
		Name:   "test-listener",
		Status: differ.Modified,
		Before: before,
		After:  after,
	}
}

func fieldValue(t *testing.T, raw, key string) string {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	v, _ := m[key].(string)
	return v
}

func TestApply_RedactsSensitiveField(t *testing.T) {
	before := `{"name":"l1","token":"abc123"}`
	after := `{"name":"l1","token":"xyz789"}`
	results := redactor.Apply([]differ.Result{makeResult(before, after)}, redactor.DefaultOptions())

	if got := fieldValue(t, results[0].Before, "token"); got != "[REDACTED]" {
		t.Errorf("Before token = %q, want [REDACTED]", got)
	}
	if got := fieldValue(t, results[0].After, "token"); got != "[REDACTED]" {
		t.Errorf("After token = %q, want [REDACTED]", got)
	}
}

func TestApply_PreservesNonSensitiveField(t *testing.T) {
	body := `{"name":"listener-a","port":8080}`
	results := redactor.Apply([]differ.Result{makeResult(body, body)}, redactor.DefaultOptions())
	if got := fieldValue(t, results[0].Before, "name"); got != "listener-a" {
		t.Errorf("name = %q, want listener-a", got)
	}
}

func TestApply_CaseInsensitiveKey(t *testing.T) {
	body := `{"TLS_CERTIFICATE":"cert-data"}`
	results := redactor.Apply([]differ.Result{makeResult(body, body)}, redactor.DefaultOptions())
	if got := fieldValue(t, results[0].Before, "TLS_CERTIFICATE"); got != "[REDACTED]" {
		t.Errorf("TLS_CERTIFICATE = %q, want [REDACTED]", got)
	}
}

func TestApply_CustomSensitiveKeys(t *testing.T) {
	body := `{"api_key":"secret-val","name":"svc"}`
	opts := redactor.Options{SensitiveKeys: []string{"api_key"}}
	results := redactor.Apply([]differ.Result{makeResult(body, body)}, opts)
	if got := fieldValue(t, results[0].Before, "api_key"); got != "[REDACTED]" {
		t.Errorf("api_key = %q, want [REDACTED]", got)
	}
	if got := fieldValue(t, results[0].Before, "name"); got != "svc" {
		t.Errorf("name = %q, want svc", got)
	}
}

func TestApply_InvalidJSONPassthrough(t *testing.T) {
	body := `not-json`
	results := redactor.Apply([]differ.Result{makeResult(body, body)}, redactor.DefaultOptions())
	if results[0].Before != body {
		t.Errorf("Before = %q, want %q", results[0].Before, body)
	}
}

func TestApply_EmptyBodyPassthrough(t *testing.T) {
	results := redactor.Apply([]differ.Result{makeResult("", "")}, redactor.DefaultOptions())
	if results[0].Before != "" || results[0].After != "" {
		t.Error("expected empty strings to pass through unchanged")
	}
}

func TestApply_OriginalUnmodified(t *testing.T) {
	body := `{"token":"secret"}`
	orig := makeResult(body, body)
	redactor.Apply([]differ.Result{orig}, redactor.DefaultOptions())
	if orig.Before != body {
		t.Error("original result should not be mutated")
	}
}
