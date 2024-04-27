package bpmnio

import "fmt"

// FindSequenceFlows finds all the sequence glows in the process and sub processes of the process
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

func (p *Process) FindSubProcessById(id string) (subProcess *SubProcess) {
	for _, sp := range p.SubProcesses {
		if sp.Id == id {
			subProcess = sp
			break
		}
	}
	return
}

func (p *Process) sortFlow(node string, visited map[string]bool, edges map[string][]string, sorted *[]string) {
	visited[node] = true

	for _, u := range edges[node] {
		if !visited[u] {
			p.sortFlow(u, visited, edges, sorted)
		}
	}

	*sorted = append([]string{node}, *sorted...)
}

func (p *Process) Flow() (flow []BaseElement) {
	nodeMap := p.FindNodes()
	linkMap := p.FindLinks()

	g := &stringGraph{
		edges:    make(map[string][]string, len(nodeMap)), // Node n connects to nodes m,o,p...
		vertices: make([]string, len(nodeMap)),
	}
	v := 0
	for _, node := range nodeMap {
		g.vertices[v] = node.GetId()
		v++
	}
	// Use Links to find edges of nodes
	for _, l := range linkMap {
		var sourceRef, targetRef string
		if msgFlow, mfOK := l.(*MessageFlow); mfOK {
			sourceRef = msgFlow.SourceRef
			targetRef = msgFlow.TargetRef
		} else if seqFlow, sfOK := l.(*SequenceFlow); sfOK {
			sourceRef = seqFlow.SourceRef
			targetRef = seqFlow.TargetRef
		}
		//fmt.Println(sourceRef, targetRef)
		if fromNode, sourceInMap := nodeMap[sourceRef]; sourceInMap {
			if toNode, targetInMap := nodeMap[targetRef]; targetInMap {
				g.addEdge(fromNode.GetId(), toNode.GetId())
			}
		}
	}
	result := g.topologicalSort()

	flow = make([]BaseElement, len(result))
	for i, r := range result {
		flow[i] = nodeMap[r]
	}
	return
}

// FindNodes finds all the nodes within a process map[id]BaseElement
func (p *Process) FindNodes() map[string]BaseElement {
	nodeTypeMap := map[ElementType]bool{
		B2Task:                   true,
		B2ScriptTask:             true,
		B2ServiceTask:            true,
		B2ManualTask:             true,
		B2UserTask:               true,
		B2CallActivity:           true,
		B2StartEvent:             true,
		B2EndEvent:               true,
		B2EventBasedGateway:      true,
		B2ParallelGateway:        true,
		B2ExclusiveGateway:       true,
		B2IntermediateCatchEvent: true,
		B2SubProcess:             false,
	}
	return p.MapBEs(nodeTypeMap)
}

// FindLinks finds all the Links within a process map[id]BaseElement
func (p *Process) FindLinks() map[string]BaseElement {
	linkTypeMap := map[ElementType]bool{
		B2MessageFlow:  true,
		B2SequenceFlow: true,
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
