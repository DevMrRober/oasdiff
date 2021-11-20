package diff

// InterfaceMap is a map of string to interface
type InterfaceMap map[string]interface{}

// InterfaceMapDiff describes the changes between a pair of InterfaceMap
type InterfaceMapDiff struct {
	Added    StringList         `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted  StringList         `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	Modified ModifiedInterfaces `json:"modified,omitempty" yaml:"modified,omitempty"`

	breaking bool // whether this diff is considered breaking within its specific context
}

// ModifiedInterfaces is map of interface names to their respective diffs
type ModifiedInterfaces map[string]*ValueDiff

// Empty indicates whether a change was found in this element
func (diff *InterfaceMapDiff) Empty() bool {
	if diff == nil {
		return true
	}

	return len(diff.Added) == 0 &&
		len(diff.Deleted) == 0 &&
		len(diff.Modified) == 0
}

// Breaking indicates whether this element includes a breaking change
func (diff *InterfaceMapDiff) Breaking() bool {
	return diff.breaking
}

func newInterfaceMapDiff() *InterfaceMapDiff {
	return &InterfaceMapDiff{
		Added:    StringList{},
		Deleted:  StringList{},
		Modified: ModifiedInterfaces{},
	}
}

func getInterfaceMapDiff(config *Config, breaking bool, map1, map2 InterfaceMap, filter StringSet) *InterfaceMapDiff {
	diff := getInterfaceMapDiffInternal(config, map1, map2, filter)

	if diff.Empty() {
		return nil
	}

	diff.breaking = breaking
	if config.BreakingOnly && !diff.Breaking() {
		return nil
	}

	return diff
}

func getInterfaceMapDiffInternal(config *Config, map1, map2 InterfaceMap, filter StringSet) *InterfaceMapDiff {

	result := newInterfaceMapDiff()

	for name1, interface1 := range map1 {
		if _, ok := filter[name1]; ok {
			if interface2, ok := map2[name1]; ok {
				if diff := getValueDiff(config, false, interface1, interface2); !diff.Empty() {
					result.Modified[name1] = diff
				}
			} else {
				result.Deleted = append(result.Deleted, name1)
			}
		}
	}

	for name2 := range map2 {
		if _, ok := filter[name2]; ok {
			if _, ok := map1[name2]; !ok {
				result.Added = append(result.Deleted, name2)
			}
		}
	}

	return result
}
