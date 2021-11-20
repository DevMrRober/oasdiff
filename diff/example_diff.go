package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ExampleDiff describes the changes between a pair of example objects: https://swagger.io/specification/#example-object
type ExampleDiff struct {
	ExtensionsDiff    *ExtensionsDiff `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	SummaryDiff       *ValueDiff      `json:"summary,omitempty" yaml:"summary,omitempty"`
	DescriptionDiff   *ValueDiff      `json:"description,omitempty" yaml:"description,omitempty"`
	ValueDiff         *ValueDiff      `json:"value,omitempty" yaml:"value,omitempty"`
	ExternalValueDiff *ValueDiff      `json:"externalValue,omitempty" yaml:"externalValue,omitempty"`
}

// Empty indicates whether a change was found in this element
func (diff *ExampleDiff) Empty() bool {
	return diff == nil || *diff == ExampleDiff{}
}

// Breaking indicates whether this element includes a breaking change
func (diff *ExampleDiff) Breaking() bool {
	return false
}

func getExampleDiff(config *Config, value1, value2 *openapi3.Example) *ExampleDiff {
	diff := getExampleDiffInternal(config, value1, value2)

	if diff.Empty() {
		return nil
	}

	if config.BreakingOnly && !diff.Breaking() {
		return nil
	}

	return diff
}

func getExampleDiffInternal(config *Config, value1, value2 *openapi3.Example) *ExampleDiff {
	result := ExampleDiff{}

	result.ExtensionsDiff = getExtensionsDiff(config, value1.ExtensionProps, value2.ExtensionProps)
	result.SummaryDiff = getValueDiff(config, false, value1.Summary, value2.Summary)
	result.DescriptionDiff = getValueDiffConditional(config, false, config.ExcludeDescription, value1.Description, value2.Description)
	result.ValueDiff = getValueDiff(config, false, value1.Value, value2.Value)
	result.ExternalValueDiff = getValueDiff(config, false, value1.ExternalValue, value2.ExternalValue)

	return &result
}
