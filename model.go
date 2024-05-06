package bpmnio

import (
	"encoding/xml"
)

type Definition struct {
	XMLName            xml.Name               `xml:"definitions"`
	Id                 string                 `xml:"id,attr"`
	Name               string                 `xml:"name,attr"`
	TargetNamespace    string                 `xml:"targetNamespace,attr"`
	ExpressionLanguage string                 `xml:"expressionLanguage,attr"`
	TypeLanguage       string                 `xml:"typeLanguage,attr"`
	Exporter           string                 `xml:"exporter,attr"`
	ExporterVersion    string                 `xml:"exporterVersion,attr"`
	Collaboration      *Collaboration         `xml:"collaboration"`
	Category           *Category              `xml:"category"`
	Processes          []*Process             `xml:"process"`
	Messages           []*TMessage            `xml:"message"`
	BpmnDiagram        []*Diagram             `xml:"BPMNDiagram"`
	_BaseElementMap    map[string]BaseElement // [Id]BaseElement
	_BaseElementGroup  map[string]*Group      // [Id]Group
	_BaseElementParent map[string]BaseElement // [Id]ParentElement (e.g. process, lane, sub-process)
	_bpmnXML           string
}

type Diagram struct {
	XMLName       xml.Name `xml:"BPMNDiagram"`
	Id            string   `xml:"id,attr"`
	Name          string   `xml:"name,attr"`
	Documentation string   `xml:"documentation,attr"`
	Resolution    string   `xml:"resolution,attr"`
	BpmnPlane     []*Plane `xml:"BPMNPlane"`
}

type Plane struct {
	XMLName     xml.Name `xml:"BPMNPlane"`
	BpmnElement string   `xml:"bpmnElement,attr"`
	BpmnShape   []*Shape `xml:"BPMNShape,omitempty"`
	BpmnEdge    []*Edge  `xml:"BPMNEdge,omitempty"`
}

type Shape struct {
	XMLName          xml.Name `xml:"BPMNShape"`
	Id               string   `xml:"id,attr"`
	BpmnElement      string   `xml:"bpmnElement,attr"`
	IsHorizontal     bool     `xml:"isHorizontal,attr"`
	IsExpanded       bool     `xml:"isExpanded,attr,omitempty"`
	IsMarkerVisible  bool     `xml:"isMarkerVisible,attr,omitempty"`
	IsMessageVisible bool     `xml:"isMessageVisible,attr,omitempty"`
	BpmnBounds       Bounds   `xml:"Bounds"`
	BpmnLabel        []*Label `xml:"BPMNLabel,omitempty"`
}

type Bounds struct {
	XMLName xml.Name `xml:"Bounds"`
	X       float32  `xml:"x,attr"`
	Y       float32  `xml:"y,attr"`
	Width   float32  `xml:"width,attr"`
	Height  float32  `xml:"height,attr"`
	area    float32
}

type Label struct {
	XMLName    xml.Name `xml:"BPMNLabel"`
	BpmnBounds Bounds   `xml:"Bounds"`
}

type Edge struct {
	XMLName      xml.Name    `xml:"BPMNEdge"`
	Id           string      `xml:"id,attr"`
	BpmnElement  string      `xml:"bpmnElement,attr"`
	BpmnWaypoint []*Waypoint `xml:"waypoint"`
}

type Waypoint struct {
	XMLName xml.Name `xml:"waypoint"`
	X       float32  `xml:"x,attr"`
	Y       float32  `xml:"y,attr"`
}

type Collaboration struct {
	XMLName           xml.Name           `xml:"collaboration"`
	Id                string             `xml:"id,attr"`
	Name              string             `xml:"name,attr"`
	Participants      []*Participant     `xml:"participant"`
	MessageFlows      []*MessageFlow     `xml:"messageFlow"`
	Groups            []*Group           `xml:"group"`
	Documentation     string             `xml:"documentation"`
	ExtensionElements *ExtensionElements `xml:"extensionElements"`
}

type Participant struct {
	XMLName           xml.Name           `xml:"participant"`
	Id                string             `xml:"id,attr"`
	Name              string             `xml:"name,attr"`
	ProcessRef        string             `xml:"processRef,attr"`
	Documentation     string             `xml:"documentation"`
	ExtensionElements *ExtensionElements `xml:"extensionElements"`
}

type MessageFlow struct {
	XMLName           xml.Name           `xml:"messageFlow"`
	Id                string             `xml:"id,attr"`
	Name              string             `xml:"name,attr"`
	SourceRef         string             `xml:"sourceRef"`
	TargetRef         string             `xml:"targetRef"`
	Documentation     string             `xml:"documentation"`
	ExtensionElements *ExtensionElements `xml:"extensionElements"`
}

type Group struct {
	XMLName           xml.Name           `xml:"group"`
	Id                string             `xml:"id,attr"`
	Name              string             `xml:"name,attr"`
	CategoryValueRef  string             `xml:"categoryValueRef,attr"`
	Documentation     string             `xml:"documentation"`
	ExtensionElements *ExtensionElements `xml:"extensionElements"`
}

type Category struct {
	XMLName           xml.Name           `xml:"category"`
	Id                string             `xml:"id,attr"`
	Name              string             `xml:"name,attr"`
	CategoryValues    []*CategoryValue   `xml:"categoryValue"`
	Documentation     string             `xml:"documentation"`
	ExtensionElements *ExtensionElements `xml:"extensionElements"`
}

type CategoryValue struct {
	XMLName           xml.Name           `xml:"categoryValue"`
	Id                string             `xml:"id,attr"`
	Value             string             `xml:"value,attr"`
	Documentation     string             `xml:"documentation"`
	ExtensionElements *ExtensionElements `xml:"extensionElements"`
}

type Process struct {
	XMLName                      xml.Name                  `xml:"process"`
	Id                           string                    `xml:"id,attr"`
	Name                         string                    `xml:"name,attr"`
	ProcessType                  string                    `xml:"processType,attr"`
	IsClosed                     bool                      `xml:"isClosed,attr"`
	IsExecutable                 bool                      `xml:"isExecutable,attr"`
	DefinitionalCollaborationRef string                    `xml:"definitionalCollaborationRef,attr"`
	Groups                       []*Group                  `xml:"group"`
	LaneSet                      []*LaneSet                `xml:"laneSet"`
	StartEvents                  []*StartEvent             `xml:"startEvent"`
	EndEvents                    []*EndEvent               `xml:"endEvent"`
	SubProcesses                 []*SubProcess             `xml:"subProcess"`
	SequenceFlows                []*SequenceFlow           `xml:"sequenceFlow"`
	Tasks                        []*Task                   `xml:"task"`
	ScriptTasks                  []*ScriptTask             `xml:"scriptTask"`
	ManualTasks                  []*ManualTask             `xml:"manualTask"`
	ServiceTasks                 []*ServiceTask            `xml:"serviceTask"`
	UserTasks                    []*UserTask               `xml:"userTask"`
	ReceiveTasks                 []*ReceiveTask            `xml:"receiveTask"`
	SendTasks                    []*SendTask               `xml:"sendTask"`
	BusinessRuleTasks            []*BusinessRuleTask       `xml:"businessRuleTask"`
	CallActivities               []*CallActivity           `xml:"callActivity"`
	ParallelGateways             []*ParallelGateway        `xml:"parallelGateway"`
	ExclusiveGateways            []*ExclusiveGateway       `xml:"exclusiveGateway"`
	DataObjectReferences         []*DataObjectReference    `xml:"dataObjectReference"`
	DataObjects                  []*DataObject             `xml:"dataObject"`
	DataStoreReferences          []*DataStoreReference     `xml:"dataStoreReference"`
	IntermediateThrowEvents      []*IntermediateThrowEvent `xml:"intermediateThrowEvent"`
	IntermediateCatchEvents      []*IntermediateCatchEvent `xml:"intermediateCatchEvent"`
	EventBasedGateways           []*EventBasedGateway      `xml:"eventBasedGateway"`
	ExtensionElements            *ExtensionElements        `xml:"extensionElements"`
	Documentation                string                    `xml:"documentation"`
	_processFlow                 []string
}

type LaneSet struct {
	XMLName           xml.Name           `xml:"laneSet"`
	Id                string             `xml:"id,attr"`
	Name              string             `xml:"name,attr"`
	Lanes             []*Lane            `xml:"lane"`
	Documentation     string             `xml:"documentation"`
	ExtensionElements *ExtensionElements `xml:"extensionElements"`
}

type Lane struct {
	XMLName           xml.Name           `xml:"lane"`
	Id                string             `xml:"id,attr"`
	Name              string             `xml:"name,attr"`
	ProcessRef        string             `xml:"processRef"`
	FlowNodeRefs      []string           `xml:"flowNodeRef"`
	Documentation     string             `xml:"documentation"`
	ExtensionElements *ExtensionElements `xml:"extensionElements"`
}

// SubProcess has same attributes as Activity + triggered by event
type SubProcess struct {
	XMLName                 xml.Name                  `xml:"subProcess"`
	Id                      string                    `xml:"id,attr"`
	Name                    string                    `xml:"name,attr"`
	Default                 string                    `xml:"default,attr"`
	StartQuantity           int                       `xml:"startQuantity,attr"`
	CompletionQuantity      int                       `xml:"completionQuantity,attr"`
	IsForCompensation       bool                      `xml:"isForCompensation,attr"`
	TriggeredByEvent        bool                      `xml:"triggeredByEvent,attr"`
	OperationRef            string                    `xml:"operationRef,attr"`
	Implementation          string                    `xml:"implementation,attr"`
	IncomingAssociations    []string                  `xml:"incoming"`
	OutgoingAssociations    []string                  `xml:"outgoing"`
	SubProcesses            []*SubProcess             `xml:"subProcess"`
	Groups                  []*Group                  `xml:"group"`
	StartEvents             []*StartEvent             `xml:"startEvent"`
	EndEvents               []*EndEvent               `xml:"endEvent"`
	SequenceFlows           []*SequenceFlow           `xml:"sequenceFlow"`
	Tasks                   []*Task                   `xml:"task"`
	ScriptTasks             []*ScriptTask             `xml:"scriptTask"`
	ManualTasks             []*ManualTask             `xml:"manualTask"`
	ServiceTasks            []*ServiceTask            `xml:"serviceTask"`
	BusinessRuleTasks       []*BusinessRuleTask       `xml:"businessRuleTask"`
	ReceiveTasks            []*ReceiveTask            `xml:"receiveTasks"`
	SendTasks               []*SendTask               `xml:"sendTask"`
	UserTasks               []*UserTask               `xml:"userTask"`
	CallActivities          []*CallActivity           `xml:"callActivity"`
	ParallelGateways        []*ParallelGateway        `xml:"parallelGateway"`
	ExclusiveGateways       []*ExclusiveGateway       `xml:"exclusiveGateway"`
	DataObjectReferences    []*DataObjectReference    `xml:"dataObjectReference"`
	DataObjects             []*DataObject             `xml:"dataObject"`
	DataStoreReferences     []*DataStoreReference     `xml:"dataStoreReference"`
	IntermediateThrowEvents []*IntermediateThrowEvent `xml:"intermediateThrowEvent"`
	IntermediateCatchEvents []*IntermediateCatchEvent `xml:"intermediateCatchEvent"`
	EventBasedGateways      []*EventBasedGateway      `xml:"eventBasedGateway"`
	Documentation           string                    `xml:"documentation"`
	ExtensionElements       *ExtensionElements        `xml:"extensionElements"`
}

type SequenceFlow struct {
	XMLName             xml.Name           `xml:"sequenceFlow"`
	Id                  string             `xml:"id,attr"`
	Name                string             `xml:"name,attr"`
	SourceRef           string             `xml:"sourceRef,attr"`
	TargetRef           string             `xml:"targetRef,attr"`
	ConditionExpression []*Expression      `xml:"conditionExpression"`
	Documentation       string             `xml:"documentation"`
	ExtensionElements   *ExtensionElements `xml:"extensionElements"`
}

type Expression struct {
	XMLName xml.Name `xml:"conditionExpression"`
	Text    string   `xml:",innerxml"`
}

type StartEvent struct {
	XMLName              xml.Name           `xml:"startEvent"`
	Id                   string             `xml:"id,attr"`
	Name                 string             `xml:"name,attr"`
	IsInterrupting       bool               `xml:"isInterrupting,attr"`
	ParallelMultiple     bool               `xml:"parallelMultiple,attr"`
	IncomingAssociations []string           `xml:"incoming"`
	OutgoingAssociations []string           `xml:"outgoing"`
	Documentation        string             `xml:"documentation"`
	ExtensionElements    *ExtensionElements `xml:"extensionElements"`
}

type EndEvent struct {
	XMLName              xml.Name           `xml:"endEvent"`
	Id                   string             `xml:"id,attr"`
	Name                 string             `xml:"name,attr"`
	IncomingAssociations []string           `xml:"incoming"`
	OutgoingAssociations []string           `xml:"outgoing"`
	Documentation        string             `xml:"documentation"`
	ExtensionElements    *ExtensionElements `xml:"extensionElements"`
}

// Task inherits the attributes and model associations of Activity
type Task struct {
	XMLName              xml.Name           `xml:"task"`
	Id                   string             `xml:"id,attr"`
	Name                 string             `xml:"name,attr"`
	Default              string             `xml:"default,attr"`
	StartQuantity        int                `xml:"startQuantity,attr"`
	CompletionQuantity   int                `xml:"completionQuantity,attr"`
	IsForCompensation    bool               `xml:"isForCompensation,attr"`
	OperationRef         string             `xml:"operationRef,attr"`
	Implementation       string             `xml:"implementation,attr"`
	IncomingAssociations []string           `xml:"incoming"`
	OutgoingAssociations []string           `xml:"outgoing"`
	Documentation        string             `xml:"documentation"`
	ExtensionElements    *ExtensionElements `xml:"extensionElements"`
}

// ManualTask has same attributes as Task
type ManualTask struct {
	XMLName              xml.Name           `xml:"manualTask"`
	Id                   string             `xml:"id,attr"`
	Name                 string             `xml:"name,attr"`
	Default              string             `xml:"default,attr"`
	StartQuantity        int                `xml:"startQuantity,attr"`
	CompletionQuantity   int                `xml:"completionQuantity,attr"`
	IsForCompensation    bool               `xml:"isForCompensation,attr"`
	OperationRef         string             `xml:"operationRef,attr"`
	IncomingAssociations []string           `xml:"incoming"`
	OutgoingAssociations []string           `xml:"outgoing"`
	Documentation        string             `xml:"documentation"`
	ExtensionElements    *ExtensionElements `xml:"extensionElements"`
}

// ScriptTask has same attributes as Task + scriptFormat & script
type ScriptTask struct {
	XMLName              xml.Name           `xml:"scriptTask"`
	Id                   string             `xml:"id,attr"`
	Name                 string             `xml:"name,attr"`
	Default              string             `xml:"default,attr"`
	StartQuantity        int                `xml:"startQuantity,attr"`
	CompletionQuantity   int                `xml:"completionQuantity,attr"`
	IsForCompensation    bool               `xml:"isForCompensation,attr"`
	OperationRef         string             `xml:"operationRef,attr"`
	ScriptFormat         string             `xml:"scriptFormat,attr"`
	Script               string             `xml:"script,attr"`
	IncomingAssociations []string           `xml:"incoming"`
	OutgoingAssociations []string           `xml:"outgoing"`
	Documentation        string             `xml:"documentation"`
	ExtensionElements    *ExtensionElements `xml:"extensionElements"`
}

// BusinessRuleTask has same attributes as Task + scriptFormat & script
type BusinessRuleTask struct {
	XMLName              xml.Name           `xml:"businessRuleTask"`
	Id                   string             `xml:"id,attr"`
	Name                 string             `xml:"name,attr"`
	Default              string             `xml:"default,attr"`
	StartQuantity        int                `xml:"startQuantity,attr"`
	CompletionQuantity   int                `xml:"completionQuantity,attr"`
	IsForCompensation    bool               `xml:"isForCompensation,attr"`
	OperationRef         string             `xml:"operationRef,attr"`
	IncomingAssociations []string           `xml:"incoming"`
	OutgoingAssociations []string           `xml:"outgoing"`
	Documentation        string             `xml:"documentation"`
	ExtensionElements    *ExtensionElements `xml:"extensionElements"`
}

// ReceiveTask has same attributes as Task + scriptFormat & script
type ReceiveTask struct {
	XMLName              xml.Name           `xml:"receiveTask"`
	Id                   string             `xml:"id,attr"`
	Name                 string             `xml:"name,attr"`
	Default              string             `xml:"default,attr"`
	StartQuantity        int                `xml:"startQuantity,attr"`
	CompletionQuantity   int                `xml:"completionQuantity,attr"`
	IsForCompensation    bool               `xml:"isForCompensation,attr"`
	OperationRef         string             `xml:"operationRef,attr"`
	IncomingAssociations []string           `xml:"incoming"`
	OutgoingAssociations []string           `xml:"outgoing"`
	Documentation        string             `xml:"documentation"`
	ExtensionElements    *ExtensionElements `xml:"extensionElements"`
}

// SendTask has same attributes as Task + scriptFormat & script
type SendTask struct {
	XMLName              xml.Name           `xml:"sendTask"`
	Id                   string             `xml:"id,attr"`
	Name                 string             `xml:"name,attr"`
	Default              string             `xml:"default,attr"`
	StartQuantity        int                `xml:"startQuantity,attr"`
	CompletionQuantity   int                `xml:"completionQuantity,attr"`
	IsForCompensation    bool               `xml:"isForCompensation,attr"`
	OperationRef         string             `xml:"operationRef,attr"`
	IncomingAssociations []string           `xml:"incoming"`
	OutgoingAssociations []string           `xml:"outgoing"`
	Documentation        string             `xml:"documentation"`
	ExtensionElements    *ExtensionElements `xml:"extensionElements"`
}

// UserTask has same attributes as Task + implementation
type UserTask struct {
	XMLName              xml.Name           `xml:"userTask"`
	Id                   string             `xml:"id,attr"`
	Name                 string             `xml:"name,attr"`
	Default              string             `xml:"default,attr"`
	StartQuantity        int                `xml:"startQuantity,attr"`
	CompletionQuantity   int                `xml:"completionQuantity,attr"`
	IsForCompensation    bool               `xml:"isForCompensation,attr"`
	OperationRef         string             `xml:"operationRef,attr"`
	Implementation       string             `xml:"implementation,attr"`
	IncomingAssociations []string           `xml:"incoming"`
	OutgoingAssociations []string           `xml:"outgoing"`
	Documentation        string             `xml:"documentation"`
	ExtensionElements    *ExtensionElements `xml:"extensionElements"`
}

// ServiceTask has same attributes as Task + implementation
type ServiceTask struct {
	XMLName              xml.Name           `xml:"serviceTask"`
	Id                   string             `xml:"id,attr"`
	Name                 string             `xml:"name,attr"`
	Default              string             `xml:"default,attr"`
	StartQuantity        int                `xml:"startQuantity,attr"`
	CompletionQuantity   int                `xml:"completionQuantity,attr"`
	IsForCompensation    bool               `xml:"isForCompensation,attr"`
	OperationRef         string             `xml:"operationRef,attr"`
	Implementation       string             `xml:"implementation,attr"`
	IncomingAssociations []string           `xml:"incoming"`
	OutgoingAssociations []string           `xml:"outgoing"`
	Documentation        string             `xml:"documentation"`
	ExtensionElements    *ExtensionElements `xml:"extensionElements"`
}

// CallActivity inherits the attributes and model associations of CallActivity
type CallActivity struct {
	XMLName              xml.Name           `xml:"callActivity"`
	Id                   string             `xml:"id,attr"`
	Name                 string             `xml:"name,attr"`
	IncomingAssociations []string           `xml:"incoming"`
	OutgoingAssociations []string           `xml:"outgoing"`
	Documentation        string             `xml:"documentation"`
	ExtensionElements    *ExtensionElements `xml:"extensionElements"`
}

type ParallelGateway struct {
	XMLName              xml.Name           `xml:"parallelGateway"`
	Id                   string             `xml:"id,attr"`
	Name                 string             `xml:"name,attr"`
	IncomingAssociations []string           `xml:"incoming"`
	OutgoingAssociations []string           `xml:"outgoing"`
	Documentation        string             `xml:"documentation"`
	ExtensionElements    *ExtensionElements `xml:"extensionElements"`
}

type ExclusiveGateway struct {
	XMLName              xml.Name           `xml:"exclusiveGateway"`
	Id                   string             `xml:"id,attr"`
	Name                 string             `xml:"name,attr"`
	IncomingAssociations []string           `xml:"incoming"`
	OutgoingAssociations []string           `xml:"outgoing"`
	Documentation        string             `xml:"documentation"`
	ExtensionElements    *ExtensionElements `xml:"extensionElements"`
}

type IntermediateThrowEvent struct {
	XMLName                xml.Name               `xml:"intermediateThrowEvent"`
	Id                     string                 `xml:"id,attr"`
	Name                   string                 `xml:"name,attr"`
	IncomingAssociations   []string               `xml:"incoming"`
	OutgoingAssociations   []string               `xml:"outgoing"`
	MessageEventDefinition MessageEventDefinition `xml:"messageEventDefinition"`
	TimerEventDefinition   TimerEventDefinition   `xml:"timerEventDefinition"`
	ParallelMultiple       bool                   `xml:"parallelMultiple"`
	Documentation          string                 `xml:"documentation"`
	ExtensionElements      *ExtensionElements     `xml:"extensionElements"`
}

type IntermediateCatchEvent struct {
	XMLName                xml.Name               `xml:"intermediateCatchEvent"`
	Id                     string                 `xml:"id,attr"`
	Name                   string                 `xml:"name,attr"`
	IncomingAssociations   []string               `xml:"incoming"`
	OutgoingAssociations   []string               `xml:"outgoing"`
	MessageEventDefinition MessageEventDefinition `xml:"messageEventDefinition"`
	TimerEventDefinition   TimerEventDefinition   `xml:"timerEventDefinition"`
	ParallelMultiple       bool                   `xml:"parallelMultiple"`
	Documentation          string                 `xml:"documentation"`
	ExtensionElements      *ExtensionElements     `xml:"extensionElements"`
}

type DataObjectReference struct {
	XMLName              xml.Name           `xml:"dataObjectReference"`
	Id                   string             `xml:"id,attr"`
	Name                 string             `xml:"name,attr"`
	IncomingAssociations []string           `xml:"incoming"`
	OutgoingAssociations []string           `xml:"outgoing"`
	DataObjectRef        string             `xml:"dataObjectRef,attr"`
	Documentation        string             `xml:"documentation"`
	ExtensionElements    *ExtensionElements `xml:"extensionElements"`
}

type DataObject struct {
	XMLName              xml.Name           `xml:"dataObject"`
	Id                   string             `xml:"id,attr"`
	Name                 string             `xml:"name,attr"`
	IncomingAssociations []string           `xml:"incoming"`
	OutgoingAssociations []string           `xml:"outgoing"`
	DataObjectRef        string             `xml:"dataObjectRef,attr"`
	Documentation        string             `xml:"documentation"`
	ExtensionElements    *ExtensionElements `xml:"extensionElements"`
}

type DataStoreReference struct {
	XMLName              xml.Name           `xml:"dataStoreReference"`
	Id                   string             `xml:"id,attr"`
	Name                 string             `xml:"name,attr"`
	IncomingAssociations []string           `xml:"incoming"`
	OutgoingAssociations []string           `xml:"outgoing"`
	Documentation        string             `xml:"documentation"`
	ExtensionElements    *ExtensionElements `xml:"extensionElements"`
}

type EventBasedGateway struct {
	XMLName              xml.Name           `xml:"eventBasedGateway"`
	Id                   string             `xml:"id,attr"`
	Name                 string             `xml:"name,attr"`
	IncomingAssociations []string           `xml:"incoming"`
	OutgoingAssociations []string           `xml:"outgoing"`
	Documentation        string             `xml:"documentation"`
	ExtensionElements    *ExtensionElements `xml:"extensionElements"`
}

type MessageEventDefinition struct {
	XMLName    xml.Name `xml:"messageEventDefinition"`
	Id         string   `xml:"id,attr"`
	MessageRef string   `xml:"messageRef,attr"`
}

type TimerEventDefinition struct {
	XMLName      xml.Name     `xml:"timerEventDefinition"`
	Id           string       `xml:"id,attr"`
	TimeDuration TimeDuration `xml:"timeDuration"`
}

type TMessage struct {
	XMLName xml.Name `xml:"message"`
	Id      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
}

type TimeDuration struct {
	XMLName xml.Name `xml:"timeDuration"`
	XMLText string   `xml:",innerxml"`
}

type ExtensionElements struct {
	XMLName xml.Name `xml:"extensionElements"`
	Rules   *Rules   `xml:"rules"`
}

type Rules struct {
	XMLName xml.Name `xml:"rules"`
	Rules   []*Rule  `xml:"rule"`
}

type Rule struct {
	XMLName     xml.Name `xml:"rule"`
	Id          string   `xml:"id,attr"`
	Type        string   `xml:"type,attr"`
	Code        string   `xml:"code,attr"`
	Name        string   `xml:"name,attr"`
	Description string   `xml:"description,attr"`
}
