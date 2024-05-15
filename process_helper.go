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

type TopologyBaseElement struct {
	BaseElement                    BaseElement
	Step                           string
	SortStep                       string
	incomingBaseElements           []BaseElement
	processedIncomingSequenceFlows int
	outgoingBaseElements           []BaseElement
}
type topology struct {
	topologyBaseElements map[string]*TopologyBaseElement
	processed            map[BaseElement]bool // have we already processed this element
	nodeMap              map[string]BaseElement
	linkMap              map[string]BaseElement
	iscLock              sync.RWMutex
	pathLock             sync.RWMutex
	dataLock             sync.RWMutex
}

// TopologicalSort
// Loosely based upon algorithm published by Kai Kretschmann
// Design and Implementation of a Framework for the Creation of BPMN 2.0 Process Models based on Textual Descriptions
// But enhanced to prefer sequence flows over number of connections
//
// includeLinks is optional
func (p *Process) TopologicalSort(includeLinks bool) []BaseElement {
	return SortTopologicalMap(p.TopologicalSortMap(includeLinks))
}

// SortTopologicalMap sorts returns an ordered list of BaseElements from the map generated by TopologicalSortMap
func SortTopologicalMap(tbeMap map[string]*TopologyBaseElement) (orderedList []BaseElement) {
	orderedList = make([]BaseElement, 0, len(tbeMap))
	for _, isc := range tbeMap {
		orderedList = append(orderedList, isc.BaseElement)
		//fmt.Println(isc.BaseElement.GetType(), isc.BaseElement.GetName(), isc.SortStep)
	}
	sort.Slice(orderedList, func(i, j int) bool {
		return tbeMap[orderedList[i].GetId()].SortStep < tbeMap[orderedList[j].GetId()].SortStep
	})
	//for _, be := range orderedList {
	//	fmt.Printf("%s: [%s]\n", be.ToString(), tbeMap[be.GetId()].SortStep)
	//}
	return
}

// TopologicalSortMap returns the TopologyBaseElement map is by bpmnid
func (p *Process) TopologicalSortMap(includeLinks bool) map[string]*TopologyBaseElement {
	t := &topology{
		topologyBaseElements: nil,
		nodeMap:              p.FindNodes(),
		linkMap:              p.FindLinks(B2Process),
		iscLock:              sync.RWMutex{}, // lock for topologyBaseElements
		pathLock:             sync.RWMutex{},
		dataLock:             sync.RWMutex{},
	}
	t.topologyBaseElements = make(map[string]*TopologyBaseElement, len(t.nodeMap)) // map[bpmnID]
	t.processed = make(map[BaseElement]bool, len(t.nodeMap))                       // map[baseEle]

	// topologyBaseElements: elements with number of incoming sequence flows & number of processed incoming sequence flows
	// If processed = total then this is next best step
	var nextElement BaseElement
	for k, be := range t.nodeMap {
		t.topologyBaseElements[k] = &TopologyBaseElement{BaseElement: be}
		if be.GetType() == B2StartEvent {
			nextElement = be
		}
		//t.orderedList = append(t.orderedList, be)
	}
	for _, l := range t.linkMap {
		switch l.(type) {
		case *SequenceFlow:
			link := l.(*SequenceFlow)
			var tbeTo, tbeFrom *TopologyBaseElement
			if tbeTo = t.topologyBaseElements[link.TargetRef]; tbeTo == nil {
				//fmt.Println("sourceRef on link isn't a node", link.SourceRef)
				continue
			}
			if tbeFrom = t.topologyBaseElements[link.SourceRef]; tbeFrom == nil {
				//fmt.Println("targetRef on link isn't a node", link.TargetRef)
				continue
			}
			tbeTo.incomingBaseElements = append(tbeTo.incomingBaseElements, tbeFrom.BaseElement)
			tbeFrom.outgoingBaseElements = append(tbeFrom.outgoingBaseElements, tbeTo.BaseElement)
		}
	}
	//link := l()
	//for _, toBpmnID := range l.GetOutgoingAssociations() {
	//	if tbeTo := t.topologyBaseElements[toBpmnID]; tbeTo == nil {
	//		fmt.Println("outgoing id on link isn't a node", toBpmnID)
	//	} else {
	//		//tbe.incomingSequenceFlows++
	//		// Now make a record the outgoing BE on the incoming ISC
	//		for _, fromBpmnID := range l.GetIncomingAssociations() {
	//			if tbeFrom := t.topologyBaseElements[fromBpmnID]; tbeFrom == nil {
	//				fmt.Println("incoming id on link isn't a node", fromBpmnID)
	//			} else {
	//				tbeFrom.outgoingBaseElements = append(tbeFrom.outgoingBaseElements, t.nodeMap[toBpmnID])
	//				tbeTo.incomingBaseElements = append(tbeTo.incomingBaseElements,t.nodeMap[fromBpmnID])
	//			}
	//		}
	//	}
	//for k, v := range t.topologyBaseElements {
	//	fmt.Println(k, len(v.incomingBaseElements), v.BaseElement.ToString())
	//}

	t.followPath("", "", 1, nextElement)
	return t.topologyBaseElements

}

// The path starts at base element and finishes at a gateway or end of the line
func (t *topology) followPath(prefix, sortPrefix string, step int, nextElement BaseElement) {
	for {
		//SortTopologicalMap(t.topologyBaseElements)
		t.iscLock.Lock()
		var tbe *TopologyBaseElement
		for {
			for _, tbeNP := range t.topologyBaseElements {
				if tbeNP.SortStep == "" && len(tbeNP.incomingBaseElements) == tbeNP.processedIncomingSequenceFlows &&
					(nextElement == nil || tbeNP.BaseElement == nextElement) {
					//fmt.Println("==", i, tbeNP)
					tbe = tbeNP
					break
				}
			}
			if tbe != nil {
				break
			}
			// if no such element can be found then
			// triple:an element in topologyBaseElements with the highest processed sequence flow number;
			var maxNumberProcessed = -1
			for _, tbeMax := range t.topologyBaseElements {
				if tbeMax.SortStep == "" && tbeMax.processedIncomingSequenceFlows > maxNumberProcessed &&
					(nextElement == nil || tbeMax.BaseElement == nextElement) {
					tbe = tbeMax
					maxNumberProcessed = tbe.processedIncomingSequenceFlows
					break
					//fmt.Println("max", i, tbe)
				}
			}
			if tbe != nil || nextElement == nil {
				break
			}
			nextElement = nil
		}
		// Found something suitable
		var outgoingBEs []BaseElement
		if tbe != nil {
			tbe.processedIncomingSequenceFlows++
			// If this is the last thing to coming into this tbe then continue
			if tbe.processedIncomingSequenceFlows == 1 {
				outgoingBEs = tbe.outgoingBaseElements
			}
			tbe.SortStep = fmt.Sprintf("%s%04d", sortPrefix, step)
			tbe.Step = fmt.Sprintf("%s%04d", prefix, step)
			t.processed[tbe.BaseElement] = true
			//if tbe.BaseElement.GetId() == "Event_1fdhkgw" {
			//	fmt.Println("Watch carefully")
			//}
		}
		t.iscLock.Unlock()
		if tbe == nil {
			break
		}

		// Now start a path for the next items
		var mainBEs []BaseElement
		mainBeMap := map[BaseElement]int{}
		for i, be := range outgoingBEs {
			if t.processed[be] {
				continue // Skip processed stuff
			}
			mainBEs = append(mainBEs, be)
			if t.isMainPath(be, nil) {
				mainBeMap[be] = i
			} else {
				mainBeMap[be] = i + len(outgoingBEs)
			}
		}
		sort.Slice(mainBEs, func(i, j int) bool {
			return mainBeMap[mainBEs[i]] < mainBeMap[mainBEs[j]]
		})
		step++
		for j, be := range mainBEs {
			if j == 0 {
				t.followPath(prefix, sortPrefix, step+j, be)
			} else if tbe != nil {
				//if tbe.BaseElement.GetName() == "Approval Needed?" {
				//	fmt.Println("Watch carefully")
				//}
				t.followPath(fmt.Sprintf("%s.", tbe.Step), fmt.Sprintf("%s.", tbe.Step), j, be)
			}
		}

		//i := 0
		//for _, be := range outgoingBEs {
		//	if t.processed[be] {
		//		continue // Skip processed stuff
		//	}
		//	if t.isMainPath(be, nil) {
		//		mainBEs = append(mainBEs, be)
		//	} else if tbe != nil {
		//		if len(outgoingBEs) == 1 {
		//			step++
		//			t.followPath(prefix, sortPrefix, step, be)
		//		} else {
		//			// Follow this path first to get it out of the way
		//			i++
		//			t.followPath(fmt.Sprintf("%s.", tbe.Step), fmt.Sprintf("%s.", tbe.SortStep), i, be)
		//		}
		//	}
		//}
		//for j, be := range mainBEs {
		//	if j == 0 {
		//		step++
		//		t.followPath(prefix, sortPrefix, step+j, be)
		//	} else if tbe != nil {
		//		t.followPath(fmt.Sprintf("%s.", tbe.Step), fmt.Sprintf("%s.", tbe.Step), j+i, be)
		//
		//	}
		//}
	}
	return
}

// Does this BE lead to the main path
// If it does then there will be a stop event or intermediate throw event at the end
func (t *topology) isMainPath(be BaseElement, visited map[BaseElement]bool) bool {
	if be == nil {
		return false
	}
	if visited == nil {
		visited = make(map[BaseElement]bool)
	}
	switch be.GetType() {
	case B2IntermediateThrowEvent, B2EndEvent:
		return true
	}
	if len(be.GetOutgoingAssociations()) == 0 {
		return false
	}
	tBEs := t.topologyBaseElements[be.GetId()]
	if tBEs != nil {
		for _, nextBE := range tBEs.outgoingBaseElements {
			if visited[nextBE] {
				continue
			}
			visited[nextBE] = true
			if t.isMainPath(nextBE, visited) {
				return true
			}
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
//	// topologyBaseElements: elements with number of incoming sequence flows & number of processed incoming sequence flows
//	// If processed = total then this is next best step
//	type TopologyBaseElement struct {
//		id                                     string
//		be                                     BaseElement
//		incomingSequenceFlows          int
//		processedIncomingSequenceFlows int
//	}
//	topologyBaseElements := make(map[string]*TopologyBaseElement, len(nodeMap))
//	for k, be := range nodeMap {
//		topologyBaseElements[k] = &TopologyBaseElement{id: k, be: be}
//	}
//	for _, l := range linkMap {
//		for _, incomingAssocBPMNId := range l.GetOutgoingAssociations() {
//			if topologyBaseElements[incomingAssocBPMNId] == nil {
//				fmt.Println("id on link isn't a node", incomingAssocBPMNId)
//			} else {
//				topologyBaseElements[incomingAssocBPMNId].incomingSequenceFlows++
//			}
//		}
//	}
//	for k := range nodeMap {
//		fmt.Println(k, topologyBaseElements[k].incomingSequenceFlows, nodeMap[k].ToString())
//	}
//	// topologyBaseElements Initialised
//
//	inOrderedList := map[string]bool{}                 // id added if in the ordered list
//	var pathPreferredNodes = make(map[string]bool, 4)  // This is to prefer a node joined by a link over anything else following the current route
//	var otherPreferredNodes = make(map[string]bool, 4) // This is to prefer a node we've seen that is a different path
//	// while topologyBaseElements not empty do
//	for len(topologyBaseElements) > 0 {
//		var isc *TopologyBaseElement
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
//			// triple: an element in topologyBaseElements with incoming sequence flows = number of processed;
//			for _, iscNP := range topologyBaseElements {
//				if iscNP.incomingSequenceFlows == iscNP.processedIncomingSequenceFlows &&
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
//			// triple:an element in topologyBaseElements with the highest processed sequence flow number;
//			var maxNumberProcessed = -1
//			for _, iscMax := range topologyBaseElements {
//				if iscMax.processedIncomingSequenceFlows > maxNumberProcessed &&
//					(preferredNodes[iscMax.id] || i == 2) {
//					isc = iscMax
//					maxNumberProcessed = isc.processedIncomingSequenceFlows
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
//		delete(topologyBaseElements, isc.id)                           //topologyBaseElements← topologyBaseElements\{triple};
//		newIncomingSequenceCounts := map[string]*TopologyBaseElement{} // newIncomingSeqCount← {};
//		for k, v := range topologyBaseElements {                       //for all x in topologyBaseElements do
//			var processedLinks []BaseElement
//			// processedEdges ←{e| e ∈ sequenceFlowArray, triple.elementID = e.predecessorID, x.elementId = e.successorID };
//			for _, l := range linkMap {
//				if len(l.GetIncomingAssociations()) > 0 && len(l.GetOutgoingAssociations()) > 0 {
//					if isc.id == l.GetIncomingAssociations()[0] && k == l.GetOutgoingAssociations()[0] {
//						processedLinks = append(processedLinks, l)
//					}
//				}
//			}
//			newIncomingSequenceCounts[k] = &TopologyBaseElement{
//				id:                                     k,
//				be:                                     v.be,
//				incomingSequenceFlows:          v.incomingSequenceFlows,
//				processedIncomingSequenceFlows: v.processedIncomingSequenceFlows + len(processedLinks),
//			}
//		}
//		topologyBaseElements = newIncomingSequenceCounts
//	}
//	return orderedList
//}
