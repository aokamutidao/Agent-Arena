package user

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCustomAgentNeverSerializesCredentials(t *testing.T) {
	agent := CustomAgent{ID: "agent", Name: "Example", APIKey: "test-secret-must-not-leak"}
	data, err := json.Marshal(map[string]interface{}{"agents": []CustomAgent{agent}})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "api_key") || strings.Contains(string(data), agent.APIKey) {
		t.Fatalf("credential leaked in public response: %s", data)
	}
}

func TestLegacyJWTRequiresConfiguredSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "")
	if _, err := GenerateJWT("user", "0x0000000000000000000000000000000000000001"); err == nil {
		t.Fatal("an unconfigured deployment must not issue JWTs")
	}
	t.Setenv("JWT_SECRET", strings.Repeat("a", 32))
	token, err := GenerateJWT("user", "0x0000000000000000000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ValidateJWT(token); err != nil {
		t.Fatal(err)
	}
	t.Setenv("JWT_SECRET", strings.Repeat("b", 32))
	if _, err := ValidateJWT(token); err == nil {
		t.Fatal("rotating the secret must reject an old token")
	}
}
