package stats

import (
	"encoding/json"
	"fmt"
)

type Rules struct {
	Rules []Rule `json:"rules"`
}
type Rule struct {
	Name                    string   `json:"name"`
	Pattern                 []string `json:"pattern"`
	ApplicationNamePosition uint     `json:"applicationNamePosition"`
}

func GetRulesFromBytes(jsonBytes []byte) (Rules, error) {
	var rules Rules
	err := json.Unmarshal(jsonBytes, &rules)
	err = checkRules(rules)
	return rules, err
}

func checkRules(rules Rules) error {
	var err error
	if len(rules.Rules) <= 0 {
		err = fmt.Errorf("No rules defined ")
	}
	isDefaultRule := false
	for i, rule := range rules.Rules {
		if len(rule.Name) <= 0 {
			err = fmt.Errorf("bad rule name `%v` at indice `%v`", rule.Name, i)
		}
		if len(rule.Pattern) == 0 {
			isDefaultRule = true
		}
	}
	if !isDefaultRule {
		err = fmt.Errorf("No default rule ( without Pattern )")
	}
	return err
}
