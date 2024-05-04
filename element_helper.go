package bpmnio

import "strings"

func NodeElementTypes() []ElementType {
	return []ElementType{
		B2Task,
		B2ScriptTask,
		B2ServiceTask,
		B2ManualTask,
		B2UserTask,
		B2ReceiveTask,
		B2SendTask,
		B2BusinessRuleTask,
		B2CallActivity,
		B2StartEvent,
		B2EndEvent,
		B2EventBasedGateway,
		B2ParallelGateway,
		B2ExclusiveGateway,
		B2IntermediateCatchEvent,
		B2SubProcess,
	}
}

func LinkElementTypes() []ElementType {
	return []ElementType{
		B2MessageFlow,
		B2SequenceFlow,
	}
}

func ContainerElementTypes() []ElementType {
	return []ElementType{
		B2Lane,
		B2Collaboration,
		B2Process,
		B2SubProcess,
	}
}

func (et ElementType) ToCamelCase(lowerCaseFirst bool) string {
	n := strings.Builder{}
	var s = string(et)
	n.Grow(len(s))
	isUpperCaseNext := !lowerCaseFirst
	isLastUpperCase := false
	for i, v := range []byte(s) {
		isUpperCase := v >= 'A' && v <= 'Z'
		isLowerCase := v >= 'a' && v <= 'z'
		if isUpperCaseNext {
			if isLowerCase {
				v += 'A'
				v -= 'a'
			}
		} else if i == 0 {
			if isUpperCase {
				v += 'a'
				v -= 'A'
			}
		} else if isLastUpperCase && isUpperCase {
			v += 'a'
			v -= 'A'
		}
		isLastUpperCase = isUpperCase

		if isUpperCase || isLowerCase {
			n.WriteByte(v)
			isUpperCaseNext = false
		} else if vIsNum := v >= '0' && v <= '9'; vIsNum {
			n.WriteByte(v)
			isUpperCaseNext = true
		} else {
			isUpperCaseNext = v == '_' || v == ' ' || v == '-' || v == '.'
		}
	}
	return n.String()
}
