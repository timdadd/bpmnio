package bpmnio

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestElementType(t *testing.T) {
	elementTypes := map[ElementType]string{
		B2Collaboration:          "collaboration",
		B2Participant:            "participant",
		B2MessageFlow:            "messageFlow",
		B2Group:                  "group",
		B2Category:               "category",
		B2CategoryValue:          "categoryValue",
		B2LaneSet:                "laneSet",
		B2Lane:                   "lane",
		B2StartEvent:             "startEvent",
		B2EndEvent:               "endEvent",
		B2Process:                "process",
		B2SubProcess:             "subProcess",
		B2Task:                   "task",
		B2ManualTask:             "manualTask",
		B2ScriptTask:             "scriptTask",
		B2UserTask:               "userTask",
		B2ServiceTask:            "serviceTask",
		B2ParallelGateway:        "parallelGateway",
		B2ExclusiveGateway:       "exclusiveGateway",
		B2IntermediateCatchEvent: "intermediateCatchEvent",
		B2EventBasedGateway:      "eventBasedGateway",
		B2CallActivity:           "callActivity",
		B2SequenceFlow:           "sequenceFlow",
	}
	for et, bpmnioTag := range elementTypes {
		assert.Equal(t, bpmnioTag, et.ToCamelCase(true))
	}
}
