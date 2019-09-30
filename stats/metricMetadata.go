package stats

import (
	"go.uber.org/zap"
	"strings"
)

type MetricMetadata struct {
	Rules        Rules
	ComponentsNb uint
}
type metric struct {
	ExtractedMetric string
	ApplicationName string
	ApplicationType string
}

func (stats *Stats) getMetric(metricPath string) metric {
	statsMetric := metric{ExtractedMetric: "None", ApplicationName: "None", ApplicationType: "None"}
	components := getComponents(metricPath, stats.MetricMetadata.ComponentsNb)
	rule := getRule(stats.Logger,components, stats.MetricMetadata.Rules)
	if rule.Name != "" {
		statsMetric.ApplicationType = rule.Name                                // rule.Name is check in rules.go
		statsMetric.ApplicationName = components[rule.ApplicationNamePosition] // the ApplicationNamePosition is check in rules.go ( must be > 0 )
		statsMetric.ExtractedMetric = strings.Join(components, ".")
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

func isMatchingRule(components []string, rule Rule) bool {
	match := true
	patternLen := len(rule.Pattern)
	if len(components) >= patternLen && patternLen > 0 {
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

func getRule(logger *zap.Logger,components []string, allRules Rules) Rule {
	i := 0
	var rule Rule
	for ; i < len(allRules.Rules) && !isMatchingRule(components, allRules.Rules[i]); i++ {
	}
	if  i < len(allRules.Rules) {
		rule = allRules.Rules[i]
	} else {
		logger.Warn("Metric Path did not match any rules",zap.Strings("metricPath",components))
	}
	return rule
}
