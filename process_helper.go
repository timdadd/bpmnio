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

// TopologicalSort
// Based upon algorithm published by Kai Kretschmann
// Design and Implementation of a Framework for the Creation of BPMN 2.0 Process Models based on Textual Descriptions
// But enhanced to prefer sequence flows over number of connections
//
// includeLinks is optional
func (p *Process) TopologicalSort(includeLinks bool) (orderedList []BaseElement) {
	nodeMap := p.FindNodes()
	linkMap := p.FindLinks() // Links/Edges with predecessor & successor
	// incomingSeqCounts: elements with number of incoming sequence flows & number of processed incoming sequence flows
	// If processed = total then this is next best step
	type incomingSeqCount struct {
		id                                     string
		be                                     BaseElement
		numberOfIncomingSequenceFlows          int
		numberOfProcessedIncomingSequenceFlows int
	}
	incomingSeqCounts := make(map[string]*incomingSeqCount, len(nodeMap))
	for k, be := range nodeMap {
		incomingSeqCounts[k] = &incomingSeqCount{id: k, be: be}
	}
	for _, l := range linkMap {
		for _, incomingAssocBPMNId := range l.GetOutgoingAssociations() {
			if incomingSeqCounts[incomingAssocBPMNId] == nil {
				fmt.Println("id on link isn't a node", incomingAssocBPMNId)
			} else {
				incomingSeqCounts[incomingAssocBPMNId].numberOfIncomingSequenceFlows++
			}
		}
	}
	//for k := range nodeMap {
	//	fmt.Println(k, incomingSeqCounts[k].numberOfIncomingSequenceFlows, nodeMap[k].ToString())
	//}
	// incomingSeqCounts Initialised

	inOrderedList := map[string]bool{}                 // id added if in the ordered list
	var pathPreferredNodes = make(map[string]bool, 4)  // This is to prefer a node joined by a link over anything else following the current route
	var otherPreferredNodes = make(map[string]bool, 4) // This is to prefer a node we've seen that is a different path
	// while incomingSeqCounts not empty do
	for len(incomingSeqCounts) > 0 {
		var isc *incomingSeqCount
		// Go around this loop twice, the first time try and find a preferred node
		//fmt.Println("preferred nodes:", len(pathPreferredNodes))
		//for id := range pathPreferredNodes {
		//	fmt.Println("..", nodeMap[id].ToString())
		//}
		for i := range 3 {
			var preferredNodes map[string]bool
			if i == 0 {
				preferredNodes = pathPreferredNodes
			}
			if i == 1 {
				preferredNodes = otherPreferredNodes
			}
			if i < 2 && len(preferredNodes) == 0 {
				continue
			}
			// triple: an element in incomingSeqCounts with incoming sequence flows = number of processed;
			for _, iscNP := range incomingSeqCounts {
				if iscNP.numberOfIncomingSequenceFlows == iscNP.numberOfProcessedIncomingSequenceFlows &&
					(preferredNodes[iscNP.id] || i == 2) {
					//fmt.Println("==", i, iscNP)
					isc = iscNP
					break
				}
			}
			if isc != nil {
				break
			}
			// if no such element can be found then
			// triple:an element in incomingSeqCounts with the highest processed sequence flow number;
			var maxNumberProcessed = -1
			for _, iscMax := range incomingSeqCounts {
				if iscMax.numberOfProcessedIncomingSequenceFlows > maxNumberProcessed &&
					(preferredNodes[iscMax.id] || i == 2) {
					isc = iscMax
					maxNumberProcessed = isc.numberOfProcessedIncomingSequenceFlows
					//fmt.Println("max", i, isc)
				}
			}
			if isc != nil {
				break
			}
		}
		if isc == nil {
			fmt.Println("How did we not find something!")
			break
		}
		//fmt.Println("Selected node", isc.be.ToString())

		// Remove the link from the map and add to ordered list if required
		var fromId string
		if len(orderedList) > 0 {
			fromId = orderedList[len(orderedList)-1].GetId()
		}
		var usedLink BaseElement
		var nextNodes []string
		for _, l := range linkMap {
			if l.GetIncomingAssociations()[0] == isc.id {
				for _, id := range l.GetOutgoingAssociations() {
					// We can only prefer things we haven't already processed
					if id != isc.id && !inOrderedList[id] {
						nextNodes = append(nextNodes, id)
					}
				}
			}
			if fromId > "" && l.GetIncomingAssociations()[0] == fromId && l.GetOutgoingAssociations()[0] == isc.id {
				usedLink = l
			}
		}
		if usedLink != nil {
			if includeLinks {
				orderedList = append(orderedList, usedLink)
			}
			delete(linkMap, usedLink.GetId())
		}
		orderedList = append(orderedList, isc.be) //orderedList← orderedList ∪ {elementID from triple};
		delete(pathPreferredNodes, isc.id)        // Take the node out of path preferred list
		delete(otherPreferredNodes, isc.id)       // Take the node out of other preferred list
		for k, v := range pathPreferredNodes {    // Move any nodes in current path preferred list to otherPreferredNodes
			otherPreferredNodes[k] = v
		}
		pathPreferredNodes = make(map[string]bool, len(nextNodes))
		for _, id := range nextNodes {
			pathPreferredNodes[id] = true
		}

		inOrderedList[isc.id] = true
		delete(incomingSeqCounts, isc.id)                           //incomingSeqCounts← incomingSeqCounts\{triple};
		newIncomingSequenceCounts := map[string]*incomingSeqCount{} // newIncomingSeqCount← {};
		for k, v := range incomingSeqCounts {                       //for all x in incomingSeqCounts do
			var processedLinks []BaseElement
			// processedEdges ←{e| e ∈ sequenceFlowArray, triple.elementID = e.predecessorID, x.elementId = e.successorID };
			for _, l := range linkMap {
				if isc.id == l.GetIncomingAssociations()[0] && k == l.GetOutgoingAssociations()[0] {
					processedLinks = append(processedLinks, l)
				}
			}
			newIncomingSequenceCounts[k] = &incomingSeqCount{
				id:                                     k,
				be:                                     v.be,
				numberOfIncomingSequenceFlows:          v.numberOfIncomingSequenceFlows,
				numberOfProcessedIncomingSequenceFlows: v.numberOfProcessedIncomingSequenceFlows + len(processedLinks),
			}
		}
		incomingSeqCounts = newIncomingSequenceCounts
	}
	return orderedList
}
