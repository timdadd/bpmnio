package bpmnio

import (
	"fmt"
	"sort"
	"sync"
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

// FindLinks finds all the Links and returns a process map[id]BaseElement
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

type incomingSeqCount struct {
	id                                     string
	be                                     BaseElement
	numberOfIncomingSequenceFlows          int
	numberOfProcessedIncomingSequenceFlows int
	outgoingBaseElements                   []BaseElement
}
type topology struct {
	incomingSeqCounts map[string]*incomingSeqCount
	orderedList       []BaseElement          // [step]baseElement
	paths             map[BaseElement]string // [step]baseElement
	nodeMap           map[string]BaseElement
	linkMap           map[string]BaseElement
	iscLock           sync.RWMutex
	pathLock          sync.RWMutex
	dataLock          sync.RWMutex
}

// TopologicalSort
// Based upon algorithm published by Kai Kretschmann
// Design and Implementation of a Framework for the Creation of BPMN 2.0 Process Models based on Textual Descriptions
// But enhanced to prefer sequence flows over number of connections
//
// includeLinks is optional
type item struct {
	step        string
	baseElement BaseElement
}

func (p *Process) TopologicalSort(includeLinks bool) (orderedList []BaseElement) {
	t := &topology{
		incomingSeqCounts: nil,
		orderedList:       nil,
		paths:             nil,
		nodeMap:           p.FindNodes(),
		linkMap:           p.FindLinks(B2Process),
		iscLock:           sync.RWMutex{}, // lock for incomingSeqCounts
		pathLock:          sync.RWMutex{},
		dataLock:          sync.RWMutex{},
	}
	t.incomingSeqCounts = make(map[string]*incomingSeqCount, len(t.nodeMap))
	t.orderedList = make([]BaseElement, 0, len(t.nodeMap))
	t.paths = make(map[BaseElement]string, len(t.nodeMap))

	// incomingSeqCounts: elements with number of incoming sequence flows & number of processed incoming sequence flows
	// If processed = total then this is next best step
	for k, be := range t.nodeMap {
		t.incomingSeqCounts[k] = &incomingSeqCount{id: k, be: be}
		t.orderedList = append(t.orderedList, be)
	}
	for _, l := range t.linkMap {
		for _, toBpmnID := range l.GetOutgoingAssociations() {
			if isc := t.incomingSeqCounts[toBpmnID]; isc == nil {
				fmt.Println("outgoing id on link isn't a node", toBpmnID)
			} else {
				isc.numberOfIncomingSequenceFlows++
				// Now make a record the outgoing BE on the incoming ISC
				for _, fromBpmnID := range l.GetIncomingAssociations() {
					if isc = t.incomingSeqCounts[fromBpmnID]; isc == nil {
						fmt.Println("incoming id on link isn't a node", fromBpmnID)
					} else {
						isc.outgoingBaseElements = append(isc.outgoingBaseElements, t.nodeMap[toBpmnID])
					}
				}
			}
		}
	}
	for k := range t.nodeMap {
		fmt.Println(k, t.incomingSeqCounts[k].numberOfIncomingSequenceFlows, t.nodeMap[k].ToString())
	}

	t.followPath("", 1, nil)

	sort.Slice(t.orderedList, func(i, j int) bool {
		return t.paths[t.orderedList[i]] < t.paths[t.orderedList[j]]
	})
	//for i, be := range t.orderedList {
	//	fmt.Println(i, t.paths[be], be.GetType(), be.GetName())
	//}
	return t.orderedList

}

// The path starts at base element and finishes at a gateway or end of the line
func (t *topology) followPath(prefix string, step int, nextElement BaseElement) {
	t.iscLock.Lock()
	var isc *incomingSeqCount
	for {
		for _, iscNP := range t.incomingSeqCounts {
			if iscNP.numberOfIncomingSequenceFlows == iscNP.numberOfProcessedIncomingSequenceFlows &&
				(nextElement == nil || iscNP.be == nextElement) {
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
		for _, iscMax := range t.incomingSeqCounts {
			if iscMax.numberOfProcessedIncomingSequenceFlows > maxNumberProcessed &&
				(nextElement == nil || iscMax.be == nextElement) {
				isc = iscMax
				maxNumberProcessed = isc.numberOfProcessedIncomingSequenceFlows
				//fmt.Println("max", i, isc)
			}
		}
		if isc != nil || nextElement == nil {
			break
		}
		nextElement = nil
	}
	// Found something suitable
	var outgoingBEs []BaseElement
	if isc != nil {
		isc.numberOfProcessedIncomingSequenceFlows++
		// If this is the last thing to coming into this isc then continue
		if isc.numberOfProcessedIncomingSequenceFlows == 1 {
			outgoingBEs = isc.outgoingBaseElements
		}
	}
	t.iscLock.Unlock()

	t.paths[isc.be] = fmt.Sprintf("%s%04d", prefix, step)

	// Now start a path for the next items
	var mainBEs []BaseElement
	i := 0
	for _, be := range outgoingBEs {
		if t.paths[be] > "" {
			continue
		}
		if t.isMainPath(be) {
			mainBEs = append(mainBEs, be)
		} else {
			// Follow this path first to get it out of the way
			i++
			t.followPath(fmt.Sprintf("%s.", t.paths[isc.be]), i, be)
		}
	}
	for j, be := range mainBEs {
		if j == 0 {
			t.followPath(prefix, step+1, be)
		} else {
			t.followPath(fmt.Sprintf("%s.", t.paths[isc.be]), j+i, be)

		}
	}
	return
}

// Does this BE lead to the main path
// If it does then there will be a stop event or intermediate throw event at the end
func (t *topology) isMainPath(be BaseElement) bool {
	if be == nil {
		return false
	}
	switch be.GetType() {
	case B2IntermediateThrowEvent, B2EndEvent:
		return true
	}
	if len(be.GetOutgoingAssociations()) == 0 {
		return false
	}
	for _, nextBE := range t.incomingSeqCounts[be.GetId()].outgoingBaseElements {
		if t.isMainPath(nextBE) {
			return true
		}
	}
	return false
}

//// TopologicalSort
//// Based upon algorithm published by Kai Kretschmann
//// Design and Implementation of a Framework for the Creation of BPMN 2.0 Process Models based on Textual Descriptions
//// But enhanced to prefer sequence flows over number of connections
////
//// includeLinks is optional
//func (p *Process) TopologicalSort(includeLinks bool) (orderedList []BaseElement) {
//	nodeMap := p.FindNodes()
//	linkMap := p.FindLinks(B2Process) // Links/Edges with predecessor & successor
//	// incomingSeqCounts: elements with number of incoming sequence flows & number of processed incoming sequence flows
//	// If processed = total then this is next best step
//	type incomingSeqCount struct {
//		id                                     string
//		be                                     BaseElement
//		numberOfIncomingSequenceFlows          int
//		numberOfProcessedIncomingSequenceFlows int
//	}
//	incomingSeqCounts := make(map[string]*incomingSeqCount, len(nodeMap))
//	for k, be := range nodeMap {
//		incomingSeqCounts[k] = &incomingSeqCount{id: k, be: be}
//	}
//	for _, l := range linkMap {
//		for _, incomingAssocBPMNId := range l.GetOutgoingAssociations() {
//			if incomingSeqCounts[incomingAssocBPMNId] == nil {
//				fmt.Println("id on link isn't a node", incomingAssocBPMNId)
//			} else {
//				incomingSeqCounts[incomingAssocBPMNId].numberOfIncomingSequenceFlows++
//			}
//		}
//	}
//	for k := range nodeMap {
//		fmt.Println(k, incomingSeqCounts[k].numberOfIncomingSequenceFlows, nodeMap[k].ToString())
//	}
//	// incomingSeqCounts Initialised
//
//	inOrderedList := map[string]bool{}                 // id added if in the ordered list
//	var pathPreferredNodes = make(map[string]bool, 4)  // This is to prefer a node joined by a link over anything else following the current route
//	var otherPreferredNodes = make(map[string]bool, 4) // This is to prefer a node we've seen that is a different path
//	// while incomingSeqCounts not empty do
//	for len(incomingSeqCounts) > 0 {
//		var isc *incomingSeqCount
//		// Go around this loop twice, the first time try and find a preferred node
//		//fmt.Println("preferred nodes:", len(pathPreferredNodes))
//		//for id := range pathPreferredNodes {
//		//	fmt.Println("..", nodeMap[id].ToString())
//		//}
//		for i := range 3 {
//			var preferredNodes map[string]bool
//			if i == 0 {
//				preferredNodes = pathPreferredNodes
//			}
//			if i == 1 {
//				preferredNodes = otherPreferredNodes
//			}
//			if i < 2 && len(preferredNodes) == 0 {
//				continue
//			}
//			// triple: an element in incomingSeqCounts with incoming sequence flows = number of processed;
//			for _, iscNP := range incomingSeqCounts {
//				if iscNP.numberOfIncomingSequenceFlows == iscNP.numberOfProcessedIncomingSequenceFlows &&
//					(preferredNodes[iscNP.id] || i == 2) {
//					//fmt.Println("==", i, iscNP)
//					isc = iscNP
//					break
//				}
//			}
//			if isc != nil {
//				break
//			}
//			// if no such element can be found then
//			// triple:an element in incomingSeqCounts with the highest processed sequence flow number;
//			var maxNumberProcessed = -1
//			for _, iscMax := range incomingSeqCounts {
//				if iscMax.numberOfProcessedIncomingSequenceFlows > maxNumberProcessed &&
//					(preferredNodes[iscMax.id] || i == 2) {
//					isc = iscMax
//					maxNumberProcessed = isc.numberOfProcessedIncomingSequenceFlows
//					//fmt.Println("max", i, isc)
//				}
//			}
//			if isc != nil {
//				break
//			}
//		}
//		if isc == nil {
//			fmt.Println("How did we not find something!")
//			break
//		}
//		//fmt.Println("Selected node", isc.be.ToString())
//
//		// Remove the link from the map and add to ordered list if required
//		var fromId string
//		if len(orderedList) > 0 {
//			fromId = orderedList[len(orderedList)-1].GetId()
//		}
//		var usedLink BaseElement
//		var nextNodes []string
//		for _, l := range linkMap {
//			for _, ia := range l.GetIncomingAssociations() {
//				if ia == isc.id {
//					for _, id := range l.GetOutgoingAssociations() {
//						// We can only prefer things we haven't already processed
//						// If this is a gateway we should not prefer to continue but wait
//						// for other paths to catch up.
//						// that is isc.processed still less than isc.incoming
//						if id != isc.id && !inOrderedList[id] {
//							nextNodes = append(nextNodes, id)
//						}
//					}
//				}
//				if fromId > "" && ia == fromId && ia == isc.id {
//					usedLink = l
//				}
//				break // Only look at first incoming association
//			}
//		}
//		if usedLink != nil {
//			if includeLinks {
//				orderedList = append(orderedList, usedLink)
//			}
//			delete(linkMap, usedLink.GetId())
//		}
//		orderedList = append(orderedList, isc.be) //orderedList← orderedList ∪ {elementID from triple};
//		delete(pathPreferredNodes, isc.id)        // Take the node out of path preferred list
//		delete(otherPreferredNodes, isc.id)       // Take the node out of other preferred list
//		for k, v := range pathPreferredNodes {    // Move any nodes in current path preferred list to otherPreferredNodes
//			otherPreferredNodes[k] = v
//		}
//		pathPreferredNodes = make(map[string]bool, len(nextNodes))
//		for _, id := range nextNodes {
//			pathPreferredNodes[id] = true
//		}
//
//		inOrderedList[isc.id] = true
//		delete(incomingSeqCounts, isc.id)                           //incomingSeqCounts← incomingSeqCounts\{triple};
//		newIncomingSequenceCounts := map[string]*incomingSeqCount{} // newIncomingSeqCount← {};
//		for k, v := range incomingSeqCounts {                       //for all x in incomingSeqCounts do
//			var processedLinks []BaseElement
//			// processedEdges ←{e| e ∈ sequenceFlowArray, triple.elementID = e.predecessorID, x.elementId = e.successorID };
//			for _, l := range linkMap {
//				if len(l.GetIncomingAssociations()) > 0 && len(l.GetOutgoingAssociations()) > 0 {
//					if isc.id == l.GetIncomingAssociations()[0] && k == l.GetOutgoingAssociations()[0] {
//						processedLinks = append(processedLinks, l)
//					}
//				}
//			}
//			newIncomingSequenceCounts[k] = &incomingSeqCount{
//				id:                                     k,
//				be:                                     v.be,
//				numberOfIncomingSequenceFlows:          v.numberOfIncomingSequenceFlows,
//				numberOfProcessedIncomingSequenceFlows: v.numberOfProcessedIncomingSequenceFlows + len(processedLinks),
//			}
//		}
//		incomingSeqCounts = newIncomingSequenceCounts
//	}
//	return orderedList
//}
