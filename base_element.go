package bpmnio

import (
	"encoding/xml"
	"fmt"
)

// BaseElement does not follow the strict BPMN 2 definition of a base element as described by [omg]
//
// [omg]: https://www.omg.org/spec/BPMN/2.0.1/PDF
type BaseElement interface {
	GetId() string
	GetName() string
	GetDocumentation() string
	GetIncomingAssociations() []string
	GetOutgoingAssociations() []string
	GetExtensionElement() *ExtensionElements
	GetType() ElementType
	GetXMLName() xml.Name
	ToString() string
	IsGateway() bool
}

// *** process methods ***
func (p *Process) GetId() string {
	return p.Id
}
func (p *Process) GetName() string {
	return p.Name
}
func (p *Process) GetDocumentation() string { return p.Documentation }
func (p *Process) GetIncomingAssociations() []string {
	return []string{}
}
func (p *Process) GetOutgoingAssociations() []string {
	return []string{}
}
func (p *Process) GetExtensionElement() *ExtensionElements { return p.ExtensionElements }
func (p *Process) GetType() ElementType {
	return B2Process
}
func (p *Process) GetXMLName() xml.Name { return p.XMLName }
func (p *Process) ToString() string {
	return fmt.Sprintf("%s:%s (%s)",
		B2Process, p.Id, p.Name)
}
func (p *Process) IsGateway() bool { return false }

// *** LaneSet methods ***
func (ls *LaneSet) GetId() string {
	return ls.Id
}
func (ls *LaneSet) GetName() string {
	return ls.Name
}
func (ls *LaneSet) GetDocumentation() string { return ls.Documentation }
func (ls *LaneSet) GetIncomingAssociations() []string {
	return []string{}
}
func (ls *LaneSet) GetOutgoingAssociations() []string {
	return []string{}
}
func (ls *LaneSet) GetExtensionElement() *ExtensionElements { return ls.ExtensionElements }
func (ls *LaneSet) GetType() ElementType {
	return B2LaneSet
}
func (ls *LaneSet) GetXMLName() xml.Name { return ls.XMLName }
func (ls *LaneSet) ToString() string {
	return fmt.Sprintf("%s:%s (%s)",
		B2LaneSet, ls.Id, ls.Name)
}
func (ls *LaneSet) IsGateway() bool { return false }

// *** Lane methods ***
func (l *Lane) GetId() string {
	return l.Id
}
func (l *Lane) GetName() string {
	return l.Name
}
func (l *Lane) GetDocumentation() string { return l.Documentation }
func (l *Lane) GetIncomingAssociations() []string {
	return []string{}
}
func (l *Lane) GetOutgoingAssociations() []string {
	return []string{}
}
func (l *Lane) GetExtensionElement() *ExtensionElements { return l.ExtensionElements }
func (l *Lane) GetType() ElementType {
	return B2Lane
}
func (l *Lane) GetXMLName() xml.Name { return l.XMLName }
func (l *Lane) ToString() string {
	return fmt.Sprintf("%s:%s (%s) %v", B2Lane, l.Id, l.Name, l.FlowNodeRefs)
}
func (l *Lane) IsGateway() bool { return false }

// *** ChildLaneSet methods ***
func (cls *ChildLaneSet) GetId() string {
	return cls.Id
}
func (cls *ChildLaneSet) GetName() string {
	return cls.Name
}
func (cls *ChildLaneSet) GetDocumentation() string { return cls.Documentation }
func (cls *ChildLaneSet) GetIncomingAssociations() []string {
	return []string{}
}
func (cls *ChildLaneSet) GetOutgoingAssociations() []string {
	return []string{}
}
func (cls *ChildLaneSet) GetExtensionElement() *ExtensionElements { return cls.ExtensionElements }
func (cls *ChildLaneSet) GetType() ElementType {
	return B2ChildLaneSet
}
func (cls *ChildLaneSet) GetXMLName() xml.Name { return cls.XMLName }
func (cls *ChildLaneSet) ToString() string {
	return fmt.Sprintf("%s:%s (%s)",
		B2ChildLaneSet, cls.Id, cls.Name)
}
func (cls *ChildLaneSet) IsGateway() bool { return false }

// ***  Collaboration methods ***
func (cn *Collaboration) GetId() string {
	return cn.Id
}
func (cn *Collaboration) GetName() string {
	return cn.Name
}
func (cn *Collaboration) GetDocumentation() string { return cn.Documentation }
func (cn *Collaboration) GetIncomingAssociations() []string {
	return []string{}
}
func (cn *Collaboration) GetOutgoingAssociations() []string {
	return []string{}
}
func (cn *Collaboration) GetExtensionElement() *ExtensionElements {
	return cn.ExtensionElements
}
func (cn *Collaboration) GetType() ElementType {
	return B2Collaboration
}
func (cn *Collaboration) GetXMLName() xml.Name { return cn.XMLName }
func (cn *Collaboration) ToString() string {
	return fmt.Sprintf("%s:%s (%s)", B2Collaboration, cn.Id, cn.Name)
}
func (cn *Collaboration) IsGateway() bool { return false }

// *** Participant methods ***
func (p *Participant) GetId() string {
	return p.Id
}
func (p *Participant) GetName() string {
	return p.Name
}
func (p *Participant) GetDocumentation() string { return p.Documentation }
func (p *Participant) GetIncomingAssociations() []string {
	return []string{}
}
func (p *Participant) GetOutgoingAssociations() []string {
	return []string{}
}
func (p *Participant) GetExtensionElement() *ExtensionElements {
	return p.ExtensionElements
}
func (p *Participant) GetType() ElementType {
	return B2Participant
}
func (p *Participant) GetXMLName() xml.Name { return p.XMLName }
func (p *Participant) ToString() string {
	return fmt.Sprintf("%s:%s (%s)", B2Participant, p.Id, p.Name)
}
func (p *Participant) IsGateway() bool { return false }

// ***  MessageFlow methods ***
func (mf *MessageFlow) GetId() string {
	return mf.Id
}
func (mf *MessageFlow) GetName() string {
	return mf.Name
}
func (mf *MessageFlow) GetDocumentation() string { return mf.Documentation }
func (mf *MessageFlow) GetIncomingAssociations() []string {
	return []string{mf.SourceRef}
}
func (mf *MessageFlow) GetOutgoingAssociations() []string {
	return []string{mf.TargetRef}
}
func (mf *MessageFlow) GetExtensionElement() *ExtensionElements {
	return mf.ExtensionElements
}
func (mf *MessageFlow) GetType() ElementType {
	return B2MessageFlow
}
func (mf *MessageFlow) GetXMLName() xml.Name { return mf.XMLName }
func (mf *MessageFlow) ToString() string {
	return fmt.Sprintf("%s:%s (%s)", B2MessageFlow, mf.Id, mf.Name)
}
func (mf *MessageFlow) IsGateway() bool { return false }

// ***  SequenceFlow methods ***
func (sf *SequenceFlow) GetId() string {
	return sf.Id
}
func (sf *SequenceFlow) GetName() string {
	return sf.Name
}
func (sf *SequenceFlow) GetDocumentation() string { return sf.Documentation }
func (sf *SequenceFlow) GetIncomingAssociations() []string {
	return []string{sf.SourceRef}
}
func (sf *SequenceFlow) GetOutgoingAssociations() []string {
	return []string{sf.TargetRef}
}
func (sf *SequenceFlow) GetExtensionElement() *ExtensionElements {
	return sf.ExtensionElements
}
func (sf *SequenceFlow) GetType() ElementType {
	return B2SequenceFlow
}
func (sf *SequenceFlow) GetXMLName() xml.Name { return sf.XMLName }
func (sf *SequenceFlow) ToString() string {
	return fmt.Sprintf("%s:%s (%s) %s --> %s", B2SequenceFlow, sf.Id, sf.Name, sf.SourceRef, sf.TargetRef)
}
func (sf *SequenceFlow) IsGateway() bool { return false }

// ***  Group methods ***
func (g *Group) GetId() string {
	return g.Id
}
func (g *Group) GetName() string {
	return g.Name
}
func (g *Group) GetDocumentation() string { return g.Documentation }
func (g *Group) GetIncomingAssociations() []string {
	return []string{}
}
func (g *Group) GetOutgoingAssociations() []string {
	return []string{}
}
func (g *Group) GetExtensionElement() *ExtensionElements {
	return g.ExtensionElements
}
func (g *Group) GetType() ElementType {
	return B2Group
}
func (g *Group) GetXMLName() xml.Name { return g.XMLName }
func (g *Group) ToString() string {
	return fmt.Sprintf("%s:%s (%s)", B2Group, g.Id, g.Name)
}
func (g *Group) IsGateway() bool { return false }

// ***  Category methods ***
func (c *Category) GetId() string {
	return c.Id
}
func (c *Category) GetName() string {
	return c.Name
}
func (c *Category) GetDocumentation() string { return c.Documentation }
func (c *Category) GetIncomingAssociations() []string {
	return []string{}
}
func (c *Category) GetOutgoingAssociations() []string {
	return []string{}
}
func (c *Category) GetExtensionElement() *ExtensionElements {
	return c.ExtensionElements
}
func (c *Category) GetType() ElementType {
	return B2Category
}
func (c *Category) GetXMLName() xml.Name { return c.XMLName }
func (c *Category) ToString() string {
	return fmt.Sprintf("%s:%s (%s)", B2Category, c.Id, c.Name)
}
func (c *Category) IsGateway() bool { return false }

// ***  CategoryValue methods ***
func (cv *CategoryValue) GetId() string {
	return cv.Id
}
func (cv *CategoryValue) GetName() string {
	return cv.Value
}
func (cv *CategoryValue) GetDocumentation() string { return cv.Documentation }
func (cv *CategoryValue) GetIncomingAssociations() []string {
	return []string{}
}
func (cv *CategoryValue) GetOutgoingAssociations() []string {
	return []string{}
}
func (cv *CategoryValue) GetExtensionElement() *ExtensionElements {
	return cv.ExtensionElements
}
func (cv *CategoryValue) GetType() ElementType {
	return B2CategoryValue
}
func (cv *CategoryValue) GetXMLName() xml.Name { return cv.XMLName }
func (cv *CategoryValue) ToString() string {
	return fmt.Sprintf("%s:%s (%s)", B2CategoryValue, cv.Id, cv.Value)
}
func (cv *CategoryValue) IsGateway() bool { return false }

// *** StartEvent methods ***
func (se *StartEvent) GetId() string {
	return se.Id
}
func (se *StartEvent) GetName() string {
	return se.Name
}
func (se *StartEvent) GetDocumentation() string { return se.Documentation }
func (se *StartEvent) GetIncomingAssociations() []string {
	return se.IncomingAssociations
}
func (se *StartEvent) GetOutgoingAssociations() []string {
	return se.OutgoingAssociations
}
func (se *StartEvent) GetExtensionElement() *ExtensionElements { return se.ExtensionElements }
func (se *StartEvent) GetType() ElementType {
	return B2StartEvent
}
func (se *StartEvent) GetXMLName() xml.Name { return se.XMLName }
func (se *StartEvent) ToString() string {
	return fmt.Sprintf("%s:%s (%s) ia=%d, oa=%d",
		B2StartEvent, se.Id, se.Name, len(se.IncomingAssociations), len(se.OutgoingAssociations))
}
func (se *StartEvent) IsGateway() bool { return false }

// *** EndEvent methods ***
func (ee *EndEvent) GetId() string {
	return ee.Id
}
func (ee *EndEvent) GetName() string {
	return ee.Name
}
func (ee *EndEvent) GetDocumentation() string { return ee.Documentation }
func (ee *EndEvent) GetIncomingAssociations() []string {
	return ee.IncomingAssociations
}
func (ee *EndEvent) GetOutgoingAssociations() []string {
	return ee.OutgoingAssociations
}
func (ee *EndEvent) GetExtensionElement() *ExtensionElements { return ee.ExtensionElements }
func (ee *EndEvent) GetType() ElementType {
	return B2EndEvent
}
func (ee *EndEvent) GetXMLName() xml.Name { return ee.XMLName }
func (ee *EndEvent) ToString() string {
	return fmt.Sprintf("%s:%s (%s) ia=%d, oa=%d",
		B2EndEvent, ee.Id, ee.Name, len(ee.IncomingAssociations), len(ee.OutgoingAssociations))
}
func (ee *EndEvent) IsGateway() bool { return false }

// *** Task methods ***
func (t *Task) GetId() string {
	return t.Id
}
func (t *Task) GetName() string {
	return t.Name
}
func (t *Task) GetDocumentation() string { return t.Documentation }
func (t *Task) GetIncomingAssociations() []string {
	return t.IncomingAssociations
}
func (t *Task) GetOutgoingAssociations() []string {
	return t.OutgoingAssociations
}
func (t *Task) GetExtensionElement() *ExtensionElements { return t.ExtensionElements }
func (t *Task) GetType() ElementType {
	return B2Task
}
func (t *Task) GetXMLName() xml.Name { return t.XMLName }
func (t *Task) ToString() string {
	return fmt.Sprintf("%s:%s (%s) ia=%d, oa=%d",
		B2Task, t.Id, t.Name, len(t.IncomingAssociations), len(t.OutgoingAssociations))
}
func (t *Task) IsGateway() bool { return false }

// *** ManualTask methods ***
func (mt *ManualTask) GetId() string {
	return mt.Id
}
func (mt *ManualTask) GetName() string {
	return mt.Name
}
func (mt *ManualTask) GetDocumentation() string { return mt.Documentation }
func (mt *ManualTask) GetIncomingAssociations() []string {
	return mt.IncomingAssociations
}
func (mt *ManualTask) GetOutgoingAssociations() []string {
	return mt.OutgoingAssociations
}
func (mt *ManualTask) GetExtensionElement() *ExtensionElements { return mt.ExtensionElements }
func (mt *ManualTask) GetType() ElementType {
	return B2ManualTask
}
func (mt *ManualTask) GetXMLName() xml.Name { return mt.XMLName }
func (mt *ManualTask) ToString() string {
	return fmt.Sprintf("%s:%s (%s) ia=%d, oa=%d",
		B2ManualTask, mt.Id, mt.Name, len(mt.IncomingAssociations), len(mt.OutgoingAssociations))
}
func (mt *ManualTask) IsGateway() bool { return false }

// *** ScriptTask methods ***
func (sct *ScriptTask) GetId() string {
	return sct.Id
}
func (sct *ScriptTask) GetName() string {
	return sct.Name
}
func (sct *ScriptTask) GetDocumentation() string { return sct.Documentation }
func (sct *ScriptTask) GetIncomingAssociations() []string {
	return sct.IncomingAssociations
}
func (sct *ScriptTask) GetOutgoingAssociations() []string {
	return sct.OutgoingAssociations
}
func (sct *ScriptTask) GetExtensionElement() *ExtensionElements { return sct.ExtensionElements }
func (sct *ScriptTask) GetType() ElementType {
	return B2ScriptTask
}
func (sct *ScriptTask) GetXMLName() xml.Name { return sct.XMLName }
func (sct *ScriptTask) ToString() string {
	return fmt.Sprintf("%s:%s (%s) ia=%d, oa=%d",
		B2ScriptTask, sct.Id, sct.Name, len(sct.IncomingAssociations), len(sct.OutgoingAssociations))
}
func (sct *ScriptTask) IsGateway() bool { return false }

// *** UserTask methods ***
func (ut *UserTask) GetId() string {
	return ut.Id
}
func (ut *UserTask) GetName() string {
	return ut.Name
}
func (ut *UserTask) GetDocumentation() string { return ut.Documentation }
func (ut *UserTask) GetIncomingAssociations() []string {
	return ut.IncomingAssociations
}
func (ut *UserTask) GetOutgoingAssociations() []string {
	return ut.OutgoingAssociations
}
func (ut *UserTask) GetExtensionElement() *ExtensionElements { return ut.ExtensionElements }
func (ut *UserTask) GetType() ElementType {
	return B2UserTask
}
func (ut *UserTask) GetXMLName() xml.Name { return ut.XMLName }
func (ut *UserTask) ToString() string {
	return fmt.Sprintf("%s:%s (%s) ia=%d, oa=%d",
		B2UserTask, ut.Id, ut.Name, len(ut.IncomingAssociations), len(ut.OutgoingAssociations))
}
func (ut *UserTask) IsGateway() bool { return false }

// *** ServiceTask methods ***
func (svt *ServiceTask) GetId() string {
	return svt.Id
}
func (svt *ServiceTask) GetName() string {
	return svt.Name
}
func (svt *ServiceTask) GetDocumentation() string { return svt.Documentation }
func (svt *ServiceTask) GetIncomingAssociations() []string {
	return svt.IncomingAssociations
}
func (svt *ServiceTask) GetOutgoingAssociations() []string {
	return svt.OutgoingAssociations
}
func (svt *ServiceTask) GetExtensionElement() *ExtensionElements { return svt.ExtensionElements }
func (svt *ServiceTask) GetType() ElementType {
	return B2ServiceTask
}
func (svt *ServiceTask) GetXMLName() xml.Name { return svt.XMLName }
func (svt *ServiceTask) ToString() string {
	return fmt.Sprintf("%s:%s (%s) ia=%d, oa=%d",
		B2ServiceTask, svt.Id, svt.Name, len(svt.IncomingAssociations), len(svt.OutgoingAssociations))
}
func (svt *ServiceTask) IsGateway() bool { return false }

// *** BusinessRuleTask methods ***
func (brt *BusinessRuleTask) GetId() string {
	return brt.Id
}
func (brt *BusinessRuleTask) GetName() string {
	return brt.Name
}
func (brt *BusinessRuleTask) GetDocumentation() string { return brt.Documentation }
func (brt *BusinessRuleTask) GetIncomingAssociations() []string {
	return brt.IncomingAssociations
}
func (brt *BusinessRuleTask) GetOutgoingAssociations() []string {
	return brt.OutgoingAssociations
}
func (brt *BusinessRuleTask) GetExtensionElement() *ExtensionElements { return brt.ExtensionElements }
func (brt *BusinessRuleTask) GetType() ElementType {
	return B2BusinessRuleTask
}
func (brt *BusinessRuleTask) GetXMLName() xml.Name { return brt.XMLName }
func (brt *BusinessRuleTask) ToString() string {
	return fmt.Sprintf("%s:%s (%s) ia=%d, oa=%d",
		B2BusinessRuleTask, brt.Id, brt.Name, len(brt.IncomingAssociations), len(brt.OutgoingAssociations))
}
func (brt *BusinessRuleTask) IsGateway() bool { return false }

// *** ReceiveTask methods ***
func (rt *ReceiveTask) GetId() string {
	return rt.Id
}
func (rt *ReceiveTask) GetName() string {
	return rt.Name
}
func (rt *ReceiveTask) GetDocumentation() string { return rt.Documentation }
func (rt *ReceiveTask) GetIncomingAssociations() []string {
	return rt.IncomingAssociations
}
func (rt *ReceiveTask) GetOutgoingAssociations() []string {
	return rt.OutgoingAssociations
}
func (rt *ReceiveTask) GetExtensionElement() *ExtensionElements { return rt.ExtensionElements }
func (rt *ReceiveTask) GetType() ElementType {
	return B2ReceiveTask
}
func (rt *ReceiveTask) GetXMLName() xml.Name { return rt.XMLName }
func (rt *ReceiveTask) ToString() string {
	return fmt.Sprintf("%s:%s (%s) ia=%d, oa=%d",
		B2ReceiveTask, rt.Id, rt.Name, len(rt.IncomingAssociations), len(rt.OutgoingAssociations))
}
func (rt *ReceiveTask) IsGateway() bool { return false }

// *** SendTask methods ***
func (st *SendTask) GetId() string {
	return st.Id
}
func (st *SendTask) GetName() string {
	return st.Name
}
func (st *SendTask) GetDocumentation() string { return st.Documentation }
func (st *SendTask) GetIncomingAssociations() []string {
	return st.IncomingAssociations
}
func (st *SendTask) GetOutgoingAssociations() []string {
	return st.OutgoingAssociations
}
func (st *SendTask) GetExtensionElement() *ExtensionElements { return st.ExtensionElements }
func (st *SendTask) GetType() ElementType {
	return B2SendTask
}
func (st *SendTask) GetXMLName() xml.Name { return st.XMLName }
func (st *SendTask) ToString() string {
	return fmt.Sprintf("%s:%s (%s) ia=%d, oa=%d",
		B2SendTask, st.Id, st.Name, len(st.IncomingAssociations), len(st.OutgoingAssociations))
}
func (st *SendTask) IsGateway() bool { return false }

// *** CallActivity methods ***
func (ca *CallActivity) GetId() string {
	return ca.Id
}
func (ca *CallActivity) GetName() string {
	return ca.Name
}
func (ca *CallActivity) GetDocumentation() string { return ca.Documentation }
func (ca *CallActivity) GetIncomingAssociations() []string {
	return ca.IncomingAssociations
}
func (ca *CallActivity) GetOutgoingAssociations() []string {
	return ca.OutgoingAssociations
}
func (ca *CallActivity) GetExtensionElement() *ExtensionElements {
	return ca.ExtensionElements
}
func (ca *CallActivity) GetType() ElementType {
	return B2CallActivity
}
func (ca *CallActivity) GetXMLName() xml.Name { return ca.XMLName }
func (ca *CallActivity) ToString() string {
	return fmt.Sprintf("%s:%s (%s) ia=%d, oa=%d",
		B2CallActivity, ca.Id, ca.Name, len(ca.IncomingAssociations), len(ca.OutgoingAssociations))
}
func (ca *CallActivity) IsGateway() bool { return false }

// *** SubProcess methods ***
func (sp *SubProcess) GetId() string {
	return sp.Id
}
func (sp *SubProcess) GetName() string {
	return sp.Name
}
func (sp *SubProcess) GetDocumentation() string { return sp.Documentation }
func (sp *SubProcess) GetIncomingAssociations() []string {
	return sp.IncomingAssociations
}
func (sp *SubProcess) GetOutgoingAssociations() []string {
	return sp.OutgoingAssociations
}
func (sp *SubProcess) GetExtensionElement() *ExtensionElements { return sp.ExtensionElements }
func (sp *SubProcess) GetType() ElementType {
	return B2SubProcess
}
func (sp *SubProcess) GetXMLName() xml.Name { return sp.XMLName }
func (sp *SubProcess) ToString() string {
	return fmt.Sprintf("%s:%s (%s) ia=%d, oa=%d",
		B2SubProcess, sp.Id, sp.Name, len(sp.IncomingAssociations), len(sp.OutgoingAssociations))
}
func (sp *SubProcess) IsGateway() bool { return false }

func (pg *ParallelGateway) GetId() string {
	return pg.Id
}
func (pg *ParallelGateway) GetName() string {
	return pg.Name
}
func (pg *ParallelGateway) GetDocumentation() string { return pg.Documentation }
func (pg *ParallelGateway) GetIncomingAssociations() []string {
	return pg.IncomingAssociations
}
func (pg *ParallelGateway) GetOutgoingAssociations() []string {
	return pg.OutgoingAssociations
}
func (pg *ParallelGateway) GetExtensionElement() *ExtensionElements {
	return pg.ExtensionElements
}
func (pg *ParallelGateway) GetType() ElementType {
	return B2ParallelGateway
}
func (pg *ParallelGateway) GetXMLName() xml.Name { return pg.XMLName }
func (pg *ParallelGateway) ToString() string {
	return fmt.Sprintf("%s:%s (%s) ia=%d, oa=%d",
		B2ParallelGateway, pg.Id, pg.Name, len(pg.IncomingAssociations), len(pg.OutgoingAssociations))
}
func (pg *ParallelGateway) IsGateway() bool { return true }

func (eg *ExclusiveGateway) GetId() string {
	return eg.Id
}
func (eg *ExclusiveGateway) GetName() string {
	return eg.Name
}
func (eg *ExclusiveGateway) GetDocumentation() string { return eg.Documentation }
func (eg *ExclusiveGateway) GetIncomingAssociations() []string {
	return eg.IncomingAssociations
}
func (eg *ExclusiveGateway) GetOutgoingAssociations() []string {
	return eg.OutgoingAssociations
}
func (eg *ExclusiveGateway) GetExtensionElement() *ExtensionElements {
	return eg.ExtensionElements
}
func (eg *ExclusiveGateway) GetType() ElementType {
	return B2ExclusiveGateway
}
func (eg *ExclusiveGateway) GetXMLName() xml.Name { return eg.XMLName }
func (eg *ExclusiveGateway) ToString() string {
	return fmt.Sprintf("%s:%s (%s) ia=%d, oa=%d",
		B2ExclusiveGateway, eg.Id, eg.Name, len(eg.IncomingAssociations), len(eg.OutgoingAssociations))
}
func (eg *ExclusiveGateway) IsGateway() bool { return true }

func (dor *DataObjectReference) GetId() string {
	return dor.Id
}
func (dor *DataObjectReference) GetName() string {
	return dor.Name
}
func (dor *DataObjectReference) GetDocumentation() string { return dor.Documentation }
func (dor *DataObjectReference) GetIncomingAssociations() []string {
	return dor.IncomingAssociations
}
func (dor *DataObjectReference) GetOutgoingAssociations() []string {
	return dor.OutgoingAssociations
}
func (dor *DataObjectReference) GetExtensionElement() *ExtensionElements {
	return dor.ExtensionElements
}
func (dor *DataObjectReference) GetType() ElementType {
	return B2DataObjectReference
}
func (dor *DataObjectReference) GetXMLName() xml.Name {
	return dor.XMLName
}
func (dor *DataObjectReference) ToString() string {
	return fmt.Sprintf("%s:%s (%s) ia=%d, oa=%d",
		B2DataObjectReference, dor.Id, dor.Name, len(dor.IncomingAssociations), len(dor.OutgoingAssociations))
}
func (dor *DataObjectReference) IsGateway() bool { return false }

func (do *DataObject) GetId() string {
	return do.Id
}
func (do *DataObject) GetName() string {
	return do.Name
}
func (do *DataObject) GetDocumentation() string { return do.Documentation }
func (do *DataObject) GetIncomingAssociations() []string {
	return do.IncomingAssociations
}
func (do *DataObject) GetOutgoingAssociations() []string {
	return do.OutgoingAssociations
}
func (do *DataObject) GetExtensionElement() *ExtensionElements {
	return do.ExtensionElements
}
func (do *DataObject) GetType() ElementType {
	return B2DataObject
}
func (do *DataObject) GetXMLName() xml.Name {
	return do.XMLName
}
func (do *DataObject) ToString() string {
	return fmt.Sprintf("%s:%s (%s) ia=%d, oa=%d",
		B2DataObject, do.Id, do.Name, len(do.IncomingAssociations), len(do.OutgoingAssociations))
}
func (do *DataObject) IsGateway() bool { return false }

func (dsr *DataStoreReference) GetId() string {
	return dsr.Id
}
func (dsr *DataStoreReference) GetName() string {
	return dsr.Name
}
func (dsr *DataStoreReference) GetDocumentation() string { return dsr.Documentation }
func (dsr *DataStoreReference) GetIncomingAssociations() []string {
	return dsr.IncomingAssociations
}
func (dsr *DataStoreReference) GetOutgoingAssociations() []string {
	return dsr.OutgoingAssociations
}
func (dsr *DataStoreReference) GetExtensionElement() *ExtensionElements {
	return dsr.ExtensionElements
}
func (dsr *DataStoreReference) GetType() ElementType {
	return B2DataStoreReference
}
func (dsr *DataStoreReference) GetXMLName() xml.Name {
	return dsr.XMLName
}
func (dsr *DataStoreReference) ToString() string {
	return fmt.Sprintf("%s:%s (%s) ia=%d, oa=%d",
		B2DataStoreReference, dsr.Id, dsr.Name, len(dsr.IncomingAssociations), len(dsr.OutgoingAssociations))
}
func (dsr *DataStoreReference) IsGateway() bool { return false }

func (ite *IntermediateThrowEvent) GetId() string {
	return ite.Id
}
func (ite *IntermediateThrowEvent) GetName() string {
	return ite.Name
}
func (ite *IntermediateThrowEvent) GetDocumentation() string { return ite.Documentation }
func (ite *IntermediateThrowEvent) GetIncomingAssociations() []string {
	return ite.IncomingAssociations
}
func (ite *IntermediateThrowEvent) GetOutgoingAssociations() []string {
	return ite.OutgoingAssociations
}
func (ite *IntermediateThrowEvent) GetExtensionElement() *ExtensionElements {
	return ite.ExtensionElements
}
func (ite *IntermediateThrowEvent) GetType() ElementType {
	return B2IntermediateThrowEvent
}
func (ite *IntermediateThrowEvent) GetXMLName() xml.Name {
	return ite.XMLName
}
func (ite *IntermediateThrowEvent) ToString() string {
	return fmt.Sprintf("%s:%s (%s) ia=%d, oa=%d",
		B2IntermediateThrowEvent, ite.Id, ite.Name, len(ite.IncomingAssociations), len(ite.OutgoingAssociations))
}
func (ite *IntermediateThrowEvent) IsGateway() bool { return false }

func (ice *IntermediateCatchEvent) GetId() string {
	return ice.Id
}
func (ice *IntermediateCatchEvent) GetName() string {
	return ice.Name
}
func (ice *IntermediateCatchEvent) GetDocumentation() string { return ice.Documentation }
func (ice *IntermediateCatchEvent) GetIncomingAssociations() []string {
	return ice.IncomingAssociations
}
func (ice *IntermediateCatchEvent) GetOutgoingAssociations() []string {
	return ice.OutgoingAssociations
}
func (ice *IntermediateCatchEvent) GetExtensionElement() *ExtensionElements {
	return ice.ExtensionElements
}
func (ice *IntermediateCatchEvent) GetType() ElementType {
	return B2IntermediateCatchEvent
}
func (ice *IntermediateCatchEvent) GetXMLName() xml.Name {
	return ice.XMLName
}
func (ice *IntermediateCatchEvent) ToString() string {
	return fmt.Sprintf("%s:%s (%s) ia=%d, oa=%d",
		B2IntermediateCatchEvent, ice.Id, ice.Name, len(ice.IncomingAssociations), len(ice.OutgoingAssociations))
}
func (ice *IntermediateCatchEvent) IsGateway() bool { return false }

func (ebg *EventBasedGateway) GetId() string {
	return ebg.Id
}
func (ebg *EventBasedGateway) GetName() string {
	return ebg.Name
}
func (ebg *EventBasedGateway) GetDocumentation() string { return ebg.Documentation }
func (ebg *EventBasedGateway) GetIncomingAssociations() []string {
	return ebg.IncomingAssociations
}
func (ebg *EventBasedGateway) GetOutgoingAssociations() []string {
	return ebg.OutgoingAssociations
}
func (ebg *EventBasedGateway) GetExtensionElement() *ExtensionElements {
	return ebg.ExtensionElements
}
func (ebg *EventBasedGateway) GetType() ElementType {
	return B2EventBasedGateway
}
func (ebg *EventBasedGateway) GetXMLName() xml.Name { return ebg.XMLName }
func (ebg *EventBasedGateway) ToString() string {
	return fmt.Sprintf("%s:%s (%s) ia=%d, oa=%d",
		B2EventBasedGateway, ebg.Id, ebg.Name, len(ebg.IncomingAssociations), len(ebg.OutgoingAssociations))
}
func (ebg *EventBasedGateway) IsGateway() bool { return true }
