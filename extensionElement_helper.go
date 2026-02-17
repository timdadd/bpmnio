package bpmnio

// GetRules avoids any null errors
func (ee *ExtensionElements) GetRules() []*Rule {
	if ee != nil && ee.Rules != nil {
		return ee.Rules.Rules
	}
	return []*Rule{} // 0 length
}

// GetDisplayOrder avoids any null errors and returns 0 if none defined
func (ee *ExtensionElements) GetDisplayOrder() int {
	if ee != nil {
		return ee.DisplayOrder
	}
	return 0 // 0 Default
}

func (ee *ExtensionElements) GetImplementation() *Implementation {
	if ee != nil && ee.Implementation != nil {
		return ee.Implementation
	}
	return nil // No definition
}
