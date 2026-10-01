package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	configurationsFileName       = "configurations"
	clabeValidationRulesFileName = "clabe_validation_rules"
)

func TestNew_Success(t *testing.T) {
	mockReader := func(name string) ([]byte, error) {
		switch name {
		case configurationsFileName:
			return []byte(`{
				"cap_provider": "accumulators",
				"transaction_intent_api":        {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"search_transaction_intent_api": {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"account_balance_api":           {"base_url": "https://api.example.com", "timeout": 30, "circuit_breaker_ratio": 0.5},
				"accumulators_api":              {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5, "optional_configurations": {"cap_id": "CAP_ID"}},
				"px_checkout_initializer_api":   {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"flow_control_engine_api":       {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"metadata_api":                  {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"public_key_api":                {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"check_user_api":                {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"expected_dates_api":            {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"audit_api":                     {"name": "test", "timeout": 1000, "retries": 1},
				"kvs_config":                    {"container": "c", "segment": "s", "read_timeout": 1000, "write_timeout": 1000, "ttl": 86400000}
			}`), nil
		case clabeValidationRulesFileName:
			return []byte(`{
				"rules": [
					{
						"clabe_number": "012914002006413756",
						"reference_required": true,
						"concept_required": true,
						"concept_length": 8
					}
				]
			}`), nil
		default:
			return nil, errors.New("unknown config file")
		}
	}

	config, err := New(mockReader)

	assert.NoError(t, err)
	assert.NotEmpty(t, config.AccountBalanceAPI.BaseURL)
	assert.Len(t, config.ClabeValidationRules.Rules, 1)
	assert.Equal(t, "012914002006413756", config.ClabeValidationRules.Rules[0].ClabeNumber)
}

func TestNew_ClabeValidationRulesError(t *testing.T) {
	mockReader := func(name string) ([]byte, error) {
		switch name {
		case configurationsFileName:
			return []byte(`{
				"account_balance_api": {
					"base_url": "https://api.example.com",
					"timeout": 30,
					"circuit_breaker_ratio": 0.5
				}
			}`), nil
		case clabeValidationRulesFileName:
			return nil, errors.New("file not found")
		default:
			return nil, errors.New("unknown config file")
		}
	}

	config, err := New(mockReader)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not get clabe validation rules")
	assert.Empty(t, config)
}

func TestNew_ClabeValidationRulesUnmarshalError(t *testing.T) {
	mockReader := func(name string) ([]byte, error) {
		switch name {
		case configurationsFileName:
			return []byte(`{
				"account_balance_api": {
					"base_url": "https://api.example.com",
					"timeout": 30,
					"circuit_breaker_ratio": 0.5
				}
			}`), nil
		case clabeValidationRulesFileName:
			return []byte(`{
				"rules": [
					{
						"clabe_number": "012914002006413756",
						"reference_required": "invalid_boolean",  // This should be boolean, not string
						"concept_required": true,
						"concept_length": 8
					}
				]
			}`), nil
		default:
			return nil, errors.New("unknown config file")
		}
	}

	config, err := New(mockReader)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could unmarshal clabe validation rules")
	assert.Empty(t, config)
}

func TestNew_ConfigurationsError(t *testing.T) {
	mockReader := func(name string) ([]byte, error) {
		switch name {
		case configurationsFileName:
			return nil, errors.New("configurations file not found")
		default:
			return nil, errors.New("unknown config file")
		}
	}

	config, err := New(mockReader)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not get configuration")
	assert.Empty(t, config)
}

func TestNew_ConfigurationsUnmarshalError(t *testing.T) {
	mockReader := func(name string) ([]byte, error) {
		switch name {
		case configurationsFileName:
			return []byte(`{
				"account_balance_api": {
					"base_url": "https://api.example.com",
					"timeout": "invalid_timeout",
					"circuit_breaker_ratio": 0.5
				}
			}`), nil
		default:
			return nil, errors.New("unknown config file")
		}
	}

	config, err := New(mockReader)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could unmarshal configuration")
	assert.Empty(t, config)
}

// validFullConfig returns a minimal but valid JSON config for use in tests.
func validFullConfig() string {
	return `{
		"cap_provider": "accumulators",
		"transaction_intent_api":        {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
		"search_transaction_intent_api": {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
		"account_balance_api":           {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
		"accumulators_api":              {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5, "optional_configurations": {"cap_id": "CAP_ID"}},
		"px_checkout_initializer_api":   {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
		"flow_control_engine_api":       {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
		"metadata_api":                  {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
		"public_key_api":                {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
		"check_user_api":                {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
		"expected_dates_api":            {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
		"audit_api":                     {"name": "test", "timeout": 1000, "retries": 1},
		"kvs_config":                    {"container": "c", "segment": "s", "read_timeout": 1000, "write_timeout": 1000, "ttl": 86400000}
	}`
}

func validClabeRules() string {
	return `{"rules": [{"clabe_number": "012914002006413756", "reference_required": true, "concept_required": true, "concept_length": 8}]}`
}

func makeReader(cfg, clabe string) Reader {
	return func(name string) ([]byte, error) {
		switch name {
		case configurationsFileName:
			return []byte(cfg), nil
		case clabeValidationRulesFileName:
			return []byte(clabe), nil
		default:
			return nil, errors.New("unknown config file")
		}
	}
}

func TestValidate_MissingRequiredTimeout(t *testing.T) {
	cases := []struct {
		name    string
		missing string
		cfg     string
	}{
		{
			name:    "transaction_intent_api timeout zero",
			missing: "transaction_intent_api",
			cfg: `{
				"cap_provider": "accumulators",
				"transaction_intent_api":        {"base_url": "http://example.com", "timeout": 0, "circuit_breaker_ratio": 0.5},
				"search_transaction_intent_api": {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"account_balance_api":           {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"accumulators_api":              {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5, "optional_configurations": {"cap_id": "CAP_ID"}},
				"px_checkout_initializer_api":   {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"flow_control_engine_api":       {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"metadata_api":                  {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"public_key_api":                {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"check_user_api":                {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"expected_dates_api":            {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"audit_api":                     {"name": "test", "timeout": 1000, "retries": 1},
				"kvs_config":                    {"container": "c", "segment": "s", "read_timeout": 1000, "write_timeout": 1000, "ttl": 86400000}
			}`,
		},
		{
			name:    "audit_api timeout zero",
			missing: "audit_api",
			cfg: `{
				"cap_provider": "accumulators",
				"transaction_intent_api":        {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"search_transaction_intent_api": {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"account_balance_api":           {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"accumulators_api":              {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5, "optional_configurations": {"cap_id": "CAP_ID"}},
				"px_checkout_initializer_api":   {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"flow_control_engine_api":       {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"metadata_api":                  {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"public_key_api":                {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"check_user_api":                {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"expected_dates_api":            {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"audit_api":                     {"name": "test", "timeout": 0, "retries": 1},
				"kvs_config":                    {"container": "c", "segment": "s", "read_timeout": 1000, "write_timeout": 1000, "ttl": 86400000}
			}`,
		},
		{
			name:    "kvs read_timeout zero",
			missing: "kvs_config.read_timeout",
			cfg: `{
				"cap_provider": "accumulators",
				"transaction_intent_api":        {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"search_transaction_intent_api": {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"account_balance_api":           {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"accumulators_api":              {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5, "optional_configurations": {"cap_id": "CAP_ID"}},
				"px_checkout_initializer_api":   {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"flow_control_engine_api":       {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"metadata_api":                  {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"public_key_api":                {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"check_user_api":                {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"expected_dates_api":            {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
				"audit_api":                     {"name": "test", "timeout": 1000, "retries": 1},
				"kvs_config":                    {"container": "c", "segment": "s", "read_timeout": 0, "write_timeout": 1000, "ttl": 86400000}
			}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := New(makeReader(tc.cfg, validClabeRules()))
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.missing)
		})
	}
}

func TestOptionalConfigurations_DiscriminatorsScheduledKey(t *testing.T) {
	t.Run("mobile-scheduled key is parsed from discriminators when present", func(t *testing.T) {
		cfg := `{
			"cap_provider": "accumulators",
			"transaction_intent_api":        {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"search_transaction_intent_api": {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"account_balance_api":           {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"accumulators_api":              {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5, "optional_configurations": {"cap_id": "CAP_ID"}},
			"px_checkout_initializer_api":   {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5, "optional_configurations": {"discriminators": {"mobile": "onepx-mlm-omega", "web": "onepx-mlm-web-omega", "mobile-scheduled": "onepx-mlm-schedules-omega"}}},
			"flow_control_engine_api":       {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"metadata_api":                  {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"public_key_api":                {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"check_user_api":                {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"expected_dates_api":            {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"audit_api":                     {"name": "test", "timeout": 1000, "retries": 1},
			"kvs_config":                    {"container": "c", "segment": "s", "read_timeout": 1000, "write_timeout": 1000, "ttl": 86400000}
		}`
		config, err := New(makeReader(cfg, validClabeRules()))
		assert.NoError(t, err)
		assert.Equal(t, "onepx-mlm-schedules-omega", config.PXCheckoutInitializerAPI.OptionalConfigurations.Discriminators["mobile-scheduled"])
	})

	t.Run("mobile-scheduled key is absent when not in discriminators", func(t *testing.T) {
		config, err := New(makeReader(validFullConfig(), validClabeRules()))
		assert.NoError(t, err)
		assert.Empty(t, config.PXCheckoutInitializerAPI.OptionalConfigurations.Discriminators["mobile-scheduled"])
	})
}

func TestValidate_OptionalAPIsTimeoutOnlyWhenPresent(t *testing.T) {
	t.Run("rule_engine_api with base_url but no timeout fails", func(t *testing.T) {
		cfg := `{
			"cap_provider": "rule_engine",
			"transaction_intent_api":        {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"search_transaction_intent_api": {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"account_balance_api":           {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"accumulators_api":              {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5, "optional_configurations": {"cap_id": "CAP_ID"}},
			"rule_engine_api":               {"base_url": "http://example.com", "timeout": 0, "circuit_breaker_ratio": 0.5},
			"px_checkout_initializer_api":   {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"flow_control_engine_api":       {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"metadata_api":                  {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"public_key_api":                {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"check_user_api":                {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"expected_dates_api":            {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"audit_api":                     {"name": "test", "timeout": 1000, "retries": 1},
			"kvs_config":                    {"container": "c", "segment": "s", "read_timeout": 1000, "write_timeout": 1000, "ttl": 86400000}
		}`
		_, err := New(makeReader(cfg, validClabeRules()))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rule_engine_api")
	})

	t.Run("rule_engine_api absent passes", func(t *testing.T) {
		_, err := New(makeReader(validFullConfig(), validClabeRules()))
		assert.NoError(t, err)
	})

	t.Run("trusted_networks_api with base_url but no timeout fails", func(t *testing.T) {
		cfg := `{
			"cap_provider": "accumulators",
			"transaction_intent_api":        {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"search_transaction_intent_api": {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"account_balance_api":           {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"accumulators_api":              {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5, "optional_configurations": {"cap_id": "CAP_ID"}},
			"px_checkout_initializer_api":   {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"flow_control_engine_api":       {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"metadata_api":                  {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"public_key_api":                {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"check_user_api":                {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"expected_dates_api":            {"base_url": "http://example.com", "timeout": 1000, "circuit_breaker_ratio": 0.5},
			"audit_api":                     {"name": "test", "timeout": 1000, "retries": 1},
			"kvs_config":                    {"container": "c", "segment": "s", "read_timeout": 1000, "write_timeout": 1000, "ttl": 86400000},
			"trusted_networks_api":          {"base_url": "http://example.com", "timeout": 0, "circuit_breaker_ratio": 0.5}
		}`
		_, err := New(makeReader(cfg, validClabeRules()))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "trusted_networks_api")
	})
}
