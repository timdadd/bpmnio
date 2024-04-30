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

// Flow converts the BPMN elements into Nodes & Links and then returns a sequential flow
// using a deepest first sort
func (p *Process) Flow(includeLinks bool) (flow []BaseElement) {
	nodeMap := p.FindNodes()
	linkMap := p.FindLinks()
	sourceNodeLinkMap := make(map[string][]BaseElement)

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
				if includeLinks {
					sourceNodeLinkMap[fromNode.GetId()] = append(sourceNodeLinkMap[fromNode.GetId()], l)
				}
			}
		}
	}
	result := g.topologicalSort()

	//flow = make([]BaseElement, len(result))
	for _, r := range result {
		flow = append(flow, nodeMap[r])
		for _, f := range sourceNodeLinkMap[r] {
			flow = append(flow, f)
		}
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
		B2ReceiveTask:            true,
		B2SendTask:               true,
		B2BusinessRuleTask:       true,
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

// G Set of nodes to sort
// L Empty list for the sorted elem.
// S Empty Set for nodes with no incoming edges
// B Set of backwards edges
// while G is non-empty do
// // search free nodes
// foreach node n
// 8 if incoming link counter of n = 0 then
// 9 insert m into S
// 10 if S is non-empty then
// 11 // ordinary top-sort
// 12 remove a node n from S and G
// 13 insert n into L
// 14 foreach node m with an edge e
// 15 from n to m do
// 16 remove e
// 17 decrement incoming link counter from m
// 18 else // cycle found
// 19 // find loop entry
// 20 foreach join j in G do
// 21 if incoming link counter of j
// 22 < initial incoming link counter of j then
// 23 J j and break
// 24 // process loop entry
// 25 foreach remaining incoming edge e of J do
// 26 replace l with backwards edge b
// 27 add b into B
// 28
// 29 output: L,B
// Flow converts the BPMN elements into Nodes & Links and then returns a sequential flow
// using a deepest first sort

// sequenceFlowArray← array of triples with predecessorID, successorID, edgeID;
//incomingSeqCounts← array of triples with elementID, number of incoming sequence
//ows, number of processed incoming sequence ows (initialized as 0);
//orderedList← {}; array to contain the topological sorted element IDs
//loopEdgeArray← {}; array to contain all sequence ow IDs considered as loop edges
//while incomingSeqCounts not empty do
//triple← an element in incomingSeqCounts with incoming sequence ows = number
//of processed;
//if no such element can be found then
//triple← an element in incomingSeqCounts with highest processed sequence ow
//number;
//loopEdges← all edgeIDs from sequenceFlowArray where successorID = triple elementID and predecessorID ∈/ orderedList;
//for all loopEdges do
//loopEdgeArray← loopEdgeArray ∪ loopEdges;
//end for
//end if
//orderedList← orderedList ∪ {elementID from triple};
//incomingSeqCounts← incomingSeqCounts\{triple};
//newIncomingSeqCount← {};
//for all x in incomingSeqCounts do
//processedEdges ←{e| e ∈ sequenceFlowArray, triple.elementID = e.predecessorID,
//x.elementId = e.successorID };
//newIncomingSeqCount ← newIncomingSeqCount ∪ {(x.elementID, x.incoming,
//x.processed + |processedEdges|)}
//end for
//incomingSeqCounts← newIncomingSeqCount;
//end while
//return loopEdgeArray, orderedList;
//func (p *Process) TopologicalSort(includeLinks bool) (flow []BaseElement) {
//	nodeMap := p.FindNodes()
//	linkMap := p.FindLinks()
//	type incomingSequenceCount struct {
//		numberOfIncomingSequenceFlows          int
//		numberOfProcessedIncomingSequenceFlows int
//	}
//	incomingSequenceCounts := make(map[string]*incomingSequenceCount, len(nodeMap))
//	for k := range nodeMap {
//		incomingSequenceCounts[k] = &incomingSequenceCount{}
//	}
//	for _, l := range linkMap {
//		for _, incomingAssocBPMNId := range l.GetOutgoingAssociations() {
//			incomingSequenceCounts[incomingAssocBPMNId].numberOfIncomingSequenceFlows++
//		}
//	}
//	for k := range nodeMap {
//		fmt.Println(k, incomingSequenceCounts[k].numberOfIncomingSequenceFlows, nodeMap[k])
//	}
//	var loopEdges []BaseElement
//	numberOfProcessedIncomingSequenceFlows := 0
//	// while incomingSeqCounts not empty do
//	for len(incomingSequenceCounts) > 0 {
//		var isc *incomingSequenceCount
//		// triple← an element in incomingSeqCounts with incoming sequence ows = number of processed;
//		for _, isc = range incomingSequenceCounts {
//			if isc.numberOfIncomingSequenceFlows == numberOfProcessedIncomingSequenceFlows {
//				break
//			}
//		}
//		// if no such element can be found then
//		// triple← an element in incomingSeqCounts with highest processed sequence flow number;
//			var maxNumberProcessed = -1
//			for _, iscMax := range incomingSequenceCounts {
//				if iscMax.numberOfProcessedIncomingSequenceFlows > maxNumberProcessed {
//isc=iscMax
//				}
//			}
//
//		}
//	}
//
//for ; len(nodeMap) > 0; {
//	for _,n:=
//}
//sourceNodeLinkMap := make(map[string][]BaseElement)
//
//g := &stringGraph{
//	edges:    make(map[string][]string, len(nodeMap)), // Node n connects to nodes m,o,p...
//	vertices: make([]string, len(nodeMap)),
//}
//v := 0
//for _, node := range nodeMap {
//	g.vertices[v] = node.GetId()
//	v++
//}
//// Use Links to find edges of nodes
//for _, l := range linkMap {
//	var sourceRef, targetRef string
//	if msgFlow, mfOK := l.(*MessageFlow); mfOK {
//		sourceRef = msgFlow.SourceRef
//		targetRef = msgFlow.TargetRef
//	} else if seqFlow, sfOK := l.(*SequenceFlow); sfOK {
//		sourceRef = seqFlow.SourceRef
//		targetRef = seqFlow.TargetRef
//	}
//	//fmt.Println(sourceRef, targetRef)
//	if fromNode, sourceInMap := nodeMap[sourceRef]; sourceInMap {
//		if toNode, targetInMap := nodeMap[targetRef]; targetInMap {
//			g.addEdge(fromNode.GetId(), toNode.GetId())
//			if includeLinks {
//				sourceNodeLinkMap[fromNode.GetId()] = append(sourceNodeLinkMap[fromNode.GetId()], l)
//			}
//		}
//	}
//}
//result := g.topologicalSort()
//
////flow = make([]BaseElement, len(result))
//for _, r := range result {
//	flow = append(flow, nodeMap[r])
//	for _, f := range sourceNodeLinkMap[r] {
//		flow = append(flow, f)
//	}
//}
//return
//}
