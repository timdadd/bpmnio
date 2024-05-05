package bpmnio

import (
	"encoding/xml"
	"fmt"
	"html"
	"strings"
)

// FindSequenceFlows finds all the sequence glows in the process and sub processes of the process that match the ids given
// If no ids given then all sequence flows returned
func (d *Definition) FindSequenceFlows(ids []string) (sequenceFlows []*SequenceFlow) {
	for _, p := range d.Processes {
		sequenceFlows = append(sequenceFlows, p.FindSequenceFlows(ids)...)
	}
	return sequenceFlows
}

// FindSourceRefs returns all the sourceRefs of the sequence flows
func FindSourceRefs(sequenceFlows []*SequenceFlow, id string) (ret []string) {
	for _, flow := range sequenceFlows {
		if id == flow.Id {
			ret = append(ret, flow.SourceRef)
		}
	}
	return
}

// FindBaseElementsByType finds all BaseElements that have a specific Type
// Interfaces are just pointers and types so no need for *
func (d *Definition) FindBaseElementsByType(t ElementType) (elements []BaseElement) {
	// Create a function that checks the type and appends if found
	appendType := func(element BaseElement) bool {
		if element.GetType() == t {
			elements = append(elements, element)
		}
		return true
	}
	d.ApplyFunctionToBaseElements(appendType)
	return
}

// FindBaseElementById finds all BaseElements that have a specific ID
func (d *Definition) FindBaseElementById(id string) (element BaseElement) {
	d.BpmnIdBaseElementMap()
	if be, inMap := d._BaseElementMap[id]; inMap {
		return be
	}
	return nil
}

// FindBaseElementsByTypeId checks the type as well
func (d *Definition) FindBaseElementsByTypeId(t ElementType, id string) (elements BaseElement) {
	if be, inMap := d._BaseElementMap[id]; inMap && be.GetType() == t {
		return be
	}
	return nil
}

// FindElements finds all base elements in the definition
func (d *Definition) FindElements() (elements []BaseElement) {
	fAppendElement := func(element BaseElement) bool {
		elements = append(elements, element)
		return true
	}
	d.ApplyFunctionToBaseElements(fAppendElement)
	return
}

// BpmnIdBaseElementMap stores and returns a mapping of baseElements by BPMN_ID if it hasn't already been created
func (d *Definition) BpmnIdBaseElementMap() map[string]BaseElement {
	if len(d._BaseElementMap) == 0 {
		elements := d.FindElements()
		d._BaseElementMap = make(map[string]BaseElement, len(elements))
		for _, be := range elements {
			d._BaseElementMap[be.GetId()] = be
		}
	}
	return d._BaseElementMap
}

// GetRules gets all rules in the definition
func (d *Definition) GetRules() (rules []*Rule) {
	fAppendRules := func(element BaseElement) bool {
		rules = append(rules, element.GetRules()...)
		return true
	}
	d.ApplyFunctionToBaseElements(fAppendRules)
	return
}

// HasConditionExpression returns true, if there's exactly 1 expression present (as by the spec)
// and there's some non-whitespace-characters available
func (sf *SequenceFlow) HasConditionExpression() bool {
	return len(sf.ConditionExpression) == 1 && len(strings.TrimSpace(sf.GetConditionExpression())) > 0
}

// GetConditionExpression returns the embedded expression. There will be a panic thrown, in case none exists!
func (sf *SequenceFlow) GetConditionExpression() string {
	return html.UnescapeString(sf.ConditionExpression[0].Text)
}

// ApplyFunctionToBaseElements applies a function to Definition and all children as long as function returns true
func (d *Definition) ApplyFunctionToBaseElements(f func(element BaseElement) bool) {
	d.Collaboration.ApplyFunctionToBaseElements(f) // Handle collaboration first
	d.Category.ApplyFunctionToBaseElements(f)      // Then categories
	for _, process := range d.Processes {
		process.ApplyFunctionToBaseElements(f)
	}
}

// ApplyFunctionToBaseElements applies a function to Collaboration and all children as long as function returns true
func (cn *Collaboration) ApplyFunctionToBaseElements(f func(element BaseElement) bool) {
	if cn == nil {
		return
	}
	if f(cn) {
		for _, participant := range cn.Participants {
			f(participant)
		}
		for _, messageFlow := range cn.MessageFlows {
			f(messageFlow)
		}
		for _, group := range cn.Groups {
			f(group)
		}
	}
}

// ApplyFunctionToBaseElements applies a function to Process and all children as long as function returns true
func (p *Process) ApplyFunctionToBaseElements(f func(element BaseElement) bool) {
	if p == nil {
		return //Nothing to see here
	}
	if f(p) { // Apply to process and then items within process
		for _, group := range p.Groups {
			f(group)
		}
		for _, laneSet := range p.LaneSet {
			laneSet.ApplyFunctionToBaseElements(f)
		}
		for _, startEvent := range p.StartEvents {
			f(startEvent)
		}
		for _, endEvent := range p.EndEvents {
			f(endEvent)
		}
		for _, subProcess := range p.SubProcesses {
			subProcess.ApplyFunctionToBaseElements(f)
		}
		for _, sequenceFlow := range p.SequenceFlows {
			f(sequenceFlow)
		}
		for _, task := range p.Tasks {
			f(task)
		}
		for _, scriptTask := range p.ScriptTasks {
			f(scriptTask)
		}
		for _, manualTask := range p.ManualTasks {
			f(manualTask)
		}
		for _, userTask := range p.UserTasks {
			f(userTask)
		}
		for _, businessRuleTask := range p.BusinessRuleTasks {
			f(businessRuleTask)
		}
		for _, receiveTask := range p.ReceiveTasks {
			f(receiveTask)
		}
		for _, sendTask := range p.SendTasks {
			f(sendTask)
		}
		for _, serviceTask := range p.ServiceTasks {
			f(serviceTask)
		}
		for _, callActivity := range p.CallActivities {
			f(callActivity)
		}
		for _, parallelGateway := range p.ParallelGateways {
			f(parallelGateway)
		}
		for _, exclusiveGateway := range p.ExclusiveGateways {
			f(exclusiveGateway)
		}
		for _, intermediateCatchEvent := range p.IntermediateCatchEvents {
			f(intermediateCatchEvent)
		}
		for _, eventBasedGateway := range p.EventBasedGateways {
			f(eventBasedGateway)
		}
	}
}

// ApplyFunctionToBaseElements applies a function to the SubProcess and all children if function returns true
func (sp *SubProcess) ApplyFunctionToBaseElements(f func(element BaseElement) bool) {
	if sp == nil {
		return // Nothing to see here
	}
	if f(sp) {
		for _, subProcess := range sp.SubProcesses {
			subProcess.ApplyFunctionToBaseElements(f)
		}
		for _, group := range sp.Groups {
			f(group)
		}
		for _, startEvent := range sp.StartEvents {
			f(startEvent)
		}
		for _, endEvent := range sp.EndEvents {
			f(endEvent)
		}
		for _, sequenceFlow := range sp.SequenceFlows {
			f(sequenceFlow)
		}
		for _, task := range sp.Tasks {
			f(task)
		}
		for _, scriptTask := range sp.ScriptTasks {
			f(scriptTask)
		}
		for _, manualTask := range sp.ManualTasks {
			f(manualTask)
		}
		for _, userTask := range sp.UserTasks {
			f(userTask)
		}
		for _, serviceTask := range sp.ServiceTasks {
			f(serviceTask)
		}
		for _, businessRuleTask := range sp.BusinessRuleTasks {
			f(businessRuleTask)
		}
		for _, receiveTask := range sp.ReceiveTasks {
			f(receiveTask)
		}
		for _, sendTask := range sp.SendTasks {
			f(sendTask)
		}
		for _, callActivity := range sp.CallActivities {
			f(callActivity)
		}
		for _, parallelGateway := range sp.ParallelGateways {
			f(parallelGateway)
		}
		for _, exclusiveGateway := range sp.ExclusiveGateways {
			f(exclusiveGateway)
		}
		for _, intermediateCatchEvent := range sp.IntermediateCatchEvents {
			f(intermediateCatchEvent)
		}
		for _, eventBasedGateway := range sp.EventBasedGateways {
			f(eventBasedGateway)
		}
	}
}

// ApplyFunctionToBaseElements applies a function to Category and all children as long as function returns true
func (c *Category) ApplyFunctionToBaseElements(f func(element BaseElement) bool) {
	if c == nil {
		return // Nothing to see here
	}
	if f(c) {
		for _, categoryValue := range c.CategoryValues {
			f(categoryValue)
		}
	}
}

// ApplyFunctionToBaseElements applies a function to LaneSet and all children as long as function returns true
func (ls *LaneSet) ApplyFunctionToBaseElements(f func(element BaseElement) bool) {
	if ls == nil {
		return // Nothing to see here
	}

	if f(ls) {
		for _, lane := range ls.Lanes {
			f(lane)
		}
	}
}

// BpmnIdParentMap assigns the most appropriate parent to each element
func (d *Definition) BpmnIdParentMap() map[string]BaseElement {
	if len(d._BaseElementParent) != 0 {
		return d._BaseElementParent
	}
	d.BpmnIdBaseElementMap()
	d._BaseElementParent = make(map[string]BaseElement, 20)
	var currentParent BaseElement
	fMap := func(be BaseElement) bool {
		switch be.(type) {
		case *Process: // Process doesn't have a parent but is a parent
			// Process can be a participant, if it is then the participant is the parent not the process
			currentParent = d._BaseElementParent[be.GetId()]
			if currentParent == nil {
				currentParent = be
			}
		case *Participant: // Participant and parent can be the same thing
			// Participant isn't naturally a parent of anything
			//if currentParent != nil {
			//	d._BaseElementParent[be.GetId()] = currentParent
			//}
			currentParent = nil
			if processRef := be.(*Participant).ProcessRef; processRef != "" {
				d._BaseElementParent[processRef] = be
			}
		case *Lane: // Lanes have a parent but also tell us who their children are
			if currentParent != nil {
				d._BaseElementParent[be.GetId()] = currentParent
			}
			for _, fnr := range be.(*Lane).FlowNodeRefs {
				d._BaseElementParent[fnr] = be // Belongs to this lane
			}
		case *SequenceFlow, *MessageFlow: // Flows don't have parents, only nodes
		default: // Everything else just inherits from the parent if one is known
			if currentParent != nil {
				if id := be.GetId(); id > "" {
					if _, inMap := d._BaseElementParent[id]; !inMap {
						d._BaseElementParent[id] = currentParent
					}
				}
			}
		}
		return true
	}
	d.ApplyFunctionToBaseElements(fMap)
	return d._BaseElementParent
}

func (d *Definition) Parent(be BaseElement) BaseElement {
	//fmt.Println(id, len(d._BaseElementParent), d._BaseElementParent[id], d._BaseElementParent[id].GetName())
	return d._BaseElementParent[be.GetId()]
}

func (d *Definition) ParentName(be BaseElement) string {
	//fmt.Println(id, len(d._BaseElementParent), d._BaseElementParent[id], d._BaseElementParent[id].GetName())
	if parent, inMap := d._BaseElementParent[be.GetId()]; inMap {
		if name := parent.GetName(); name > "" {
			return name
		}
		return parent.GetId()
	}
	return "No Parent"
}

func NewDefinition(bpmnXML []byte) (d *Definition, err error) {
	if err = xml.Unmarshal(bpmnXML, &d); err != nil {
		return nil, fmt.Errorf("could not unmarshal XML into definitions, got %v", err)
	}
	d._bpmnXML = string(bpmnXML)
	// Build the cache maps
	d.BpmnIdBaseElementMap()
	d.BpmnIdGroupMap()
	d.BpmnIdParentMap()
	return
}
