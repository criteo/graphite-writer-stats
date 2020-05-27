package stats

import (
	"go.uber.org/zap"
	"strings"
)

type MetricMetadata struct {
	Rules        Rules
	ComponentsNb uint
}

// ExtractedMetric will be filled with with AplicationName, Type & rebuilt MetricPath from matching rule.
type ExtractedMetric struct {
	ExtractedMetric string
	ApplicationName string
	ApplicationType string
}

func (stats *Stats) getMetric(metricPath string, metricTags map[string]string) ExtractedMetric {
	statsMetric := ExtractedMetric{ExtractedMetric: "None", ApplicationName: "None", ApplicationType: "None"}
	components := getComponents(metricPath, stats.MetricMetadata.ComponentsNb)
	rule := getRule(components, metricTags, stats.MetricMetadata.Rules)
	if rule.Name == "" {
		stats.Logger.Warn("Metric Path did not match any rules", zap.String("metricPath", metricPath))
	} else if int(rule.ApplicationNamePosition) < len(components) {
		statsMetric.ApplicationType = rule.Name // rule.Name is check in rules.go
		if tag, hasTag := getMatchingTag(metricTags, rule); hasTag {
			statsMetric.ApplicationName = metricTags[tag]
		} else {
			statsMetric.ApplicationName = components[rule.ApplicationNamePosition] // the ApplicationNamePosition is check in rules.go ( must be > 0 )
		}
		statsMetric.ExtractedMetric = strings.Join(components, ".")
	} else {
		stats.Logger.Error("bad metric ", zap.String("metricPath", metricPath), zap.String("rule", rule.Name))
	}
	return statsMetric
}

func getComponents(metricPath string, componentsLen uint) []string {

	currentIndex := 0
	var componentIndex uint = 0
	nextDotIndex := strings.IndexByte(metricPath[currentIndex:], '.')
	components := make([]string, componentsLen)
	for ; componentIndex < componentsLen && nextDotIndex != -1; componentIndex, nextDotIndex = componentIndex+1, strings.IndexByte(metricPath[currentIndex:], '.') {
		components[componentIndex] = metricPath[currentIndex : currentIndex+nextDotIndex]
		currentIndex += nextDotIndex + 1
	}
	if nextDotIndex == -1 && componentIndex < componentsLen {
		components[componentIndex] = metricPath[currentIndex:]
		components = components[:componentIndex+1]
	}

	return components
}

func getMatchingTag(tags map[string]string, rule Rule) (string, bool) {
	if 0 == len(rule.UseTags) {
		return "", false
	}

	for k := range tags {
		for _, label := range rule.UseTags {
			if label == k {
				return k, true
			}
		}
	}

	return "", false
}

func isMatchingRule(components []string, tags map[string]string, rule Rule) bool {
	_, match := getMatchingTag(tags, rule)

	// If rule is a tag rule, skip path component test.
	if 0 != len(rule.UseTags) {
		return match
	}

	patternLen := len(rule.Pattern)
	if patternLen == 0 {
		match = true
	} else if len(components) >= patternLen && patternLen > 0 {
		extractedComponent := components[0:patternLen]
		match = cheapEqual(rule.Pattern, extractedComponent)
	}
	return match
}

func cheapEqual(array1 []string, array2 []string) bool {
	equals := false
	if len(array2) == len(array1) {
		i := 0
		for ; i < len(array1) && array1[i] == array2[i]; i++ {

		}
		if i == len(array1) {
			equals = true
		}
	}
	return equals
}

func getRule(components []string, metricTags map[string]string, allRules Rules) Rule {
	i := 0
	var rule Rule
	for ; i < len(allRules.Rules) && !isMatchingRule(components, metricTags, allRules.Rules[i]); i++ {
	}
	if i < len(allRules.Rules) {
		rule = allRules.Rules[i]
	}
	return rule
}
