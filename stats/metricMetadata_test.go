package stats

import (
	"reflect"
	"testing"
)

func TestCheapEquals(t *testing.T) {
	isEquals := cheapEqual([]string{"a"}, []string{"a"})
	if !isEquals {
		t.Error("array should be equals ")
	}
	isEquals = cheapEqual([]string{"a"}, []string{"b"})
	if isEquals {
		t.Error("array should not be equals ")
	}
	isEquals = cheapEqual([]string{"a"}, nil)
	if isEquals {
		t.Error("array should not be equals ")
	}
	isEquals = cheapEqual(nil, []string{"a"})
	if isEquals {
		t.Error("array should not be equals ")
	}
	isEquals = cheapEqual(nil, nil)
	if !isEquals {
		t.Error("array should  be equals ")
	}
}
func TestIsMatchingRule(t *testing.T) {
	rule := Rule{"aggregated", []string{"criteo", "aggregated"}, 2}
	isMatchedRule := isMatchingRule([]string{"criteo", "aggregated", "d"}, rule)
	if !isMatchedRule {
		t.Error("should match the rule but for now it doesn't ")
	}
	rule = Rule{"aggregated", []string{"criteo", "agglo"}, 2}
	isMatchedRule = isMatchingRule([]string{"criteo", "aggregated", "d"}, rule)
	if isMatchedRule {
		t.Error("should not match the rule but for now it doesn't ")
	}
	var rule2 Rule
	isMatchedRule = isMatchingRule([]string{"criteo", "aggregated", "d"}, rule2)
	if !isMatchedRule {
		t.Error("should match because an empty rule is a default rule ! ")
	}
}

func TestGetComponents(t *testing.T) {
	metricPath := "a.b.c.d"
	components := getComponents("a.b.c.d", 3)
	componentsExpected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(components, componentsExpected) {
		t.Errorf("components not good for `%v` actual: `%v` expected:`%v`", metricPath, components, componentsExpected)
	}
	components = getComponents("a.b.c.d", 4)
	componentsExpected = []string{"a", "b", "c", "d"}
	if !reflect.DeepEqual(components, componentsExpected) {
		t.Errorf("components not good for `%v` actual: `%v` expected:`%v`", metricPath, components, componentsExpected)
	}
	//More than possible return the full components
	components = getComponents("a.b.c.d", 5)
	componentsExpected = []string{"a", "b", "c", "d"}
	if !reflect.DeepEqual(components, componentsExpected) {
		t.Errorf("components not good for `%v` actual: `%v` expected:`%v`", metricPath, components, componentsExpected)
	}
	//Less than possible  return empty array
	components = getComponents("a.b.c.d", 0)
	if len(components) != 0 {
		t.Errorf("components not good for `%v` actual: `%v` expected len of 0", metricPath, components)
	}
	//More than possible return the full components
	components = getComponents("a", 1)
	componentsExpected = []string{"a"}
	if !reflect.DeepEqual(components, componentsExpected) {
		t.Errorf("components not good for `%v` actual: `%v` expected:`%v`", metricPath, components, componentsExpected)
	}
}

func TestIsPrometheus(t *testing.T) {
	isPrometheus := isPrometheusMetric([]string{"prometheus", "myapp"})
	if !isPrometheus {
		t.Error("should be a prometheus metric")
	}
	isPrometheus = isPrometheusMetric([]string{"myapp"})
	if isPrometheus {
		t.Error("should not be a prometheus metric")
	}
	isPrometheus = isPrometheusMetric([]string{})
	if isPrometheus {
		t.Error("should not be a prometheus metric")
	}
}
func TestGetRule(t *testing.T) {
	aggregatedRule := Rule{"aggregated", []string{"criteo", "aggregated"}, 2}
	aggloRule := Rule{"agglo", []string{"criteo", "agglo"}, 2}
	aggregatedAllRule := Rule{"aggregated-all", []string{"criteo", "aggregated-all"}, 2}
	legacyHostingRule := Rule{"legacy-hosting", []string{"prometheus", "hosting"}, 1}
	startWithCriteoRule := Rule{"start-by-criteo", []string{"criteo"}, 1}
	startbyAppRule := Rule{"start-by-app", nil, 0}
	rules := Rules{Rules: []Rule{aggregatedRule, aggloRule, aggregatedAllRule, legacyHostingRule, startWithCriteoRule, startbyAppRule}}
	rule := getRule([]string{"criteo", "aggregated", "myapp"}, rules)
	if rule.Name != aggregatedRule.Name {
		t.Error("should be an aggregated metric")
	}
}
