package stats

import (
	"reflect"
	"testing"
)

func TestGetRules(t *testing.T) {
	var jsonRules = []byte(`{
  "rules": [
    {
      "name": "aggregated",
      "pattern": [
        "criteo",
        "aggregated"
      ],
      "applicationNamePosition": 2
    },
    {
      "name": "agglo",
      "pattern": [
        "criteo",
        "agglo"
      ],
      "applicationNamePosition": 2
    },
    {
      "name": "aggregated-all",
      "pattern": [
        "criteo",
        "aggregated-all"
      ],
      "applicationNamePosition": 2
    },
    {
      "name": "legacy-hosting",
      "pattern": [
        "prometheus",
        "hosting"
      ],
      "applicationNamePosition": 1
    },
    {
      "name": "start-by-criteo",
      "pattern": [
        "criteo"
      ],
      "applicationNamePosition": 1
    },
    {
      "name": "start-by-app",
      "applicationNamePosition": 0
    }
  ]
}`)
	aggregatedRule := Rule{"aggregated", []string{"criteo", "aggregated"}, 2}
	aggloRule := Rule{"agglo", []string{"criteo", "agglo"}, 2}
	aggregatedAllRule := Rule{"aggregated-all", []string{"criteo", "aggregated-all"}, 2}
	legacyHostingRule := Rule{"legacy-hosting", []string{"prometheus", "hosting"}, 1}
	startWithCriteoRule := Rule{"start-by-criteo", []string{"criteo"}, 1}
	startbyAppRule := Rule{"start-by-app", nil, 0}
	rulesExpected := []Rule{aggregatedRule, aggloRule, aggregatedAllRule, legacyHostingRule, startWithCriteoRule, startbyAppRule}
	rules, err := GetRulesFromBytes(jsonRules)
	if (!reflect.DeepEqual(rules.Rules, rulesExpected)) || err != nil {
		t.Errorf("fail to parse rules : expected: '%v' actual: '%v', err: '%v'", rulesExpected, rules.Rules, err)
	}
}
func TestCheckRules(t *testing.T) {
	aggregatedRule := Rule{"aggregated", []string{"criteo", "aggregated"}, 2}
	aggloRule := Rule{"agglo", []string{"criteo", "agglo"}, 2}
	aggregatedAllRule := Rule{"aggregated-all", []string{"criteo", "aggregated-all"}, 2}
	legacyHostingRule := Rule{"legacy-hosting", []string{"prometheus", "hosting"}, 1}
	startWithCriteoRule := Rule{"start-by-criteo", []string{"criteo"}, 1}
	startbyAppRule := Rule{"start-by-app", nil, 0}
	rules := Rules{Rules: []Rule{aggregatedRule, aggloRule, aggregatedAllRule, legacyHostingRule, startWithCriteoRule, startbyAppRule}}
	err := checkRules(rules)
	if err != nil {
		t.Errorf("should not get the error: `%v`", err)
	}
	startbyAppRule = Rule{"", nil, 0}
	rules = Rules{Rules: []Rule{aggregatedRule, aggloRule, aggregatedAllRule, legacyHostingRule, startWithCriteoRule, startbyAppRule}}
	err = checkRules(rules)
	if err == nil {
		t.Errorf("the rule should have a name: `%v`", err)
	}
	startbyAppRule = Rule{"aa", []string{"za"}, 0}
	rules = Rules{Rules: []Rule{aggregatedRule, aggloRule, aggregatedAllRule, legacyHostingRule, startWithCriteoRule, startbyAppRule}}
	err = checkRules(rules)
	if err == nil {
		t.Errorf("the rules should have a default rule (without pattern) : `%v`", err)
	}
	err = checkRules(Rules{nil})
	if err == nil {
		t.Errorf("the rules should have a default rule : `%v`", err)
	}
}
