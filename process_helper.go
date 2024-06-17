package bpmnio

import (
	"depgraph"
	"fmt"
)

// FindSequenceFlows finds all the sequence flows in the process and sub processes of the process
// that have match one of the ids passed in
func (p *Process) FindSequenceFlows(ids []string) (sequenceFlows []*SequenceFlow) {
	if len(ids) == 0 {
		sequenceFlows = p.SequenceFlows
	} else {
		for _, flow := range p.SequenceFlows {
			for _, id := range ids {
				if id == flow.Id {
					sequenceFlows = append(sequenceFlows, flow)
				}
			}
		}
	}

	for _, subProcess := range p.SubProcesses {
		if len(ids) == 0 {
			sequenceFlows = append(sequenceFlows, subProcess.SequenceFlows...)
		} else {
			for _, flow := range subProcess.SequenceFlows {
				for _, id := range ids {
					if id == flow.Id {
						sequenceFlows = append(sequenceFlows, flow)
					}
				}
			}
		}
	}
	return sequenceFlows
}

// FindSubProcessById will find the sub process with the Id given
func (p *Process) FindSubProcessById(id string) (subProcess *SubProcess) {
	for _, sp := range p.SubProcesses {
		if sp.Id == id {
			subProcess = sp
			break
		}
	}
	return
}

// FindNodes finds all the nodes and returns with a process map[id]BaseElement
func (p *Process) FindNodes() map[string]BaseElement {
	nodeTypes := NodeElementTypes()
	nodeTypeMap := make(map[ElementType]bool, len(nodeTypes))
	for _, nodeType := range nodeTypes {
		nodeTypeMap[nodeType] = true
	}
	nodeTypeMap[B2SubProcess] = false
	return p.MapBEs(nodeTypeMap)
}

// FindLinks finds all the Links and returns a map[id]BaseElement
func (p *Process) FindLinks(et ElementType) map[string]BaseElement {
	linkTypes := LinkElementTypes(et)
	linkTypeMap := make(map[ElementType]bool, len(linkTypes))
	for _, linkType := range linkTypes {
		linkTypeMap[linkType] = true
	}
	return p.MapBEs(linkTypeMap)
}

// MapBEs takes a map of ElementTypes and builds a list of BaseElement that match the map
func (p *Process) MapBEs(etMap map[ElementType]bool) (m map[string]BaseElement) {
	m = make(map[string]BaseElement, 20)
	fMap := func(be BaseElement) bool {
		var r = true
		if downHierarchy, inMap := etMap[be.GetType()]; inMap {
			m[be.GetId()] = be
			r = downHierarchy
		}
		//fmt.Printf("checking BE %s %t (%d)\n", be.ToString(), r, len(m))
		return r
	}
	p.ApplyFunctionToBaseElements(fMap)
	return
}

// FindIntermediateCatchEventsForContinuation finds all intermediate catch events of a process
func (p *Process) FindIntermediateCatchEventsForContinuation() (ICEs []*IntermediateCatchEvent) {
	// find potential event definitions
	ICEs = p.IntermediateCatchEvents
	return p.EliminateEventsWhichComeFromTheSameGateway(ICEs)
}

func (p *Process) EliminateEventsWhichComeFromTheSameGateway(events []*IntermediateCatchEvent) (ICEs []*IntermediateCatchEvent) {
	// a bubble-sort-like approach to find elements, which have the same incoming association
	for len(events) > 0 {
		event := events[0]
		events = events[1:]
		if event == nil {
			continue
		}
		ICEs = append(ICEs, event)
		for i := 0; i < len(events); i++ {
			if p.haveEqualInboundBaseElement(event, events[i]) {
				events[i] = nil
			}
		}
	}
	return
}

func (p *Process) haveEqualInboundBaseElement(event1 *IntermediateCatchEvent, event2 *IntermediateCatchEvent) bool {
	if event1 == nil || event2 == nil {
		return false
	}
	p.checkOnlyOneAssociationOrPanic(event1)
	p.checkOnlyOneAssociationOrPanic(event2)
	ref1 := FindSourceRefs(p.SequenceFlows, event1.IncomingAssociations[0])[0]
	ref2 := FindSourceRefs(p.SequenceFlows, event2.IncomingAssociations[0])[0]
	baseElement1 := p.FindBaseElementById(ref1)
	baseElement2 := p.FindBaseElementById(ref2)
	return baseElement1.GetId() == baseElement2.GetId()
}

func (p *Process) checkOnlyOneAssociationOrPanic(event *IntermediateCatchEvent) {
	if len(event.IncomingAssociations) != 1 {
		panic(fmt.Sprintf("Element with id=%s has %d incoming associations, but only 1 is supported by this engine.",
			event.Id, len(event.IncomingAssociations)))
	}
}

// FindBaseElementById finds all base elements in the definition
func (p *Process) FindBaseElementById(id string) (baseElement BaseElement) {
	fFindElementById := func(be BaseElement) bool {
		if be.GetId() == id {
			baseElement = be
		}
		return baseElement != nil
	}
	p.ApplyFunctionToBaseElements(fFindElementById)
	return
}

type TopologyBaseElement struct {
	BaseElement BaseElement `json:"base_element"`
	Step        string      `json:"step"`
	SortStep    string      `json:"sort_step"`
	Level       int         `json:"level"`
}

func (p *Process) TopologicalSort(includeLinks bool) (tbe []*TopologyBaseElement) {
	g := depgraph.New()
	nodeMap := p.FindNodes()
	for _, l := range p.FindLinks(B2Process) {
		switch l.(type) {
		case *SequenceFlow:
			fromNode := nodeMap[l.GetIncomingAssociations()[0]]
			toNode := nodeMap[l.GetOutgoingAssociations()[0]]
			// Ignore sequence flows that don't have known nodes
			if fromNode != nil && toNode != nil {
				_ = g.AddLink(l.GetName(), fromNode.GetId(), toNode.GetId())
			}
		}
	}
	sortedBPMN := g.SortedWithOrder()
	tbe = make([]*TopologyBaseElement, len(sortedBPMN))
	for i, s := range sortedBPMN {
		tbe[i] = &TopologyBaseElement{
			BaseElement: nodeMap[s.Node.(string)],
			Step:        s.Step,
			SortStep:    s.SortedStep,
			Level:       s.Level,
		}
	}
	return
}
