package bpmnio

type ElementCategory string

const (
	Activity   ElementCategory = "activity"
	Event      ElementCategory = "event"
	Gateway    ElementCategory = "gateway"
	Flow       ElementCategory = "flow"
	Swimlane   ElementCategory = "swimlane"
	Artifact   ElementCategory = "artifact"
	Collection ElementCategory = "collection"
)

type ElementType string

const (
	B2Participant            ElementType = "PARTICIPANT"
	B2MessageFlow            ElementType = "MESSAGE_FLOW"
	B2Group                  ElementType = "GROUP"
	B2Category               ElementType = "CATEGORY"
	B2CategoryValue          ElementType = "CATEGORY_VALUE"
	B2LaneSet                ElementType = "LANE_SET"
	B2Lane                   ElementType = "LANE"
	B2ChildLaneSet           ElementType = "CHILD_LANE_SET"
	B2StartEvent             ElementType = "START_EVENT"
	B2EndEvent               ElementType = "END_EVENT"
	B2Process                ElementType = "PROCESS"
	B2SubProcess             ElementType = "SUB_PROCESS"
	B2Task                   ElementType = "TASK"
	B2ManualTask             ElementType = "MANUAL_TASK"
	B2ScriptTask             ElementType = "SCRIPT_TASK"
	B2UserTask               ElementType = "USER_TASK"
	B2ReceiveTask            ElementType = "RECEIVE_TASK"
	B2SendTask               ElementType = "SEND_TASK"
	B2BusinessRuleTask       ElementType = "BUSINESS_RULE_TASK"
	B2Collaboration          ElementType = "COLLABORATION"
	B2ServiceTask            ElementType = "SERVICE_TASK"
	B2ParallelGateway        ElementType = "PARALLEL_GATEWAY"
	B2ExclusiveGateway       ElementType = "EXCLUSIVE_GATEWAY"
	B2DataObjectReference    ElementType = "DATA_OBJECT_REFERENCE"
	B2DataObject             ElementType = "DATA_OBJECT"
	B2DataStoreReference     ElementType = "DATA_STORE_REFERENCE"
	B2IntermediateCatchEvent ElementType = "INTERMEDIATE_CATCH_EVENT"
	B2IntermediateThrowEvent ElementType = "INTERMEDIATE_THROW_EVENT"
	B2EventBasedGateway      ElementType = "EVENT_BASED_GATEWAY"
	B2CallActivity           ElementType = "CALL_ACTIVITY"
	B2SequenceFlow           ElementType = "SEQUENCE_FLOW"
)

var elmCatMap = map[ElementType]ElementCategory{
	B2Participant:            Swimlane,
	B2MessageFlow:            Flow,
	B2Group:                  Artifact,
	B2Category:               Artifact,
	B2CategoryValue:          Artifact,
	B2LaneSet:                Collection,
	B2Lane:                   Swimlane,
	B2ChildLaneSet:           Collection,
	B2StartEvent:             Event,
	B2EndEvent:               Event,
	B2Process:                Swimlane,
	B2SubProcess:             Activity,
	B2Task:                   Activity,
	B2ManualTask:             Activity,
	B2ScriptTask:             Activity,
	B2UserTask:               Activity,
	B2ReceiveTask:            Activity,
	B2SendTask:               Activity,
	B2BusinessRuleTask:       Activity,
	B2Collaboration:          Collection,
	B2ServiceTask:            Activity,
	B2ParallelGateway:        Gateway,
	B2ExclusiveGateway:       Gateway,
	B2DataObjectReference:    Artifact,
	B2DataObject:             Artifact,
	B2DataStoreReference:     Artifact,
	B2IntermediateCatchEvent: Event,
	B2IntermediateThrowEvent: Event,
	B2EventBasedGateway:      Gateway,
	B2CallActivity:           Activity,
	B2SequenceFlow:           Flow,
}

func (et ElementType) Category() ElementCategory {
	return elmCatMap[et]
}
