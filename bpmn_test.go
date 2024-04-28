package bpmnio

import (
	"encoding/xml"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

// TestBPMNXML reads the BPMN files and unmarshalls them, then marshalls them and finallt ch
func TestBPMNXML(t *testing.T) {
	t.Log("TestUnmarshalXML")
	dir := "./bpmn"
	items, _ := os.ReadDir(dir)
	for _, item := range items {
		if item.IsDir() {
			continue
		}
		// Only test with the test BPMN
		if !strings.HasPrefix(item.Name(), "test") {
			continue
		}
		t.Logf("Found file %s", item.Name())
		// Read the file
		var bpmnXML []byte
		var err error
		if bpmnXML, err = os.ReadFile(dir + "/" + item.Name()); err != nil {
			t.Fatalf("could not read the XML file, got %v", err)
		}

		// Now unmarshall into structs
		var origD Definition
		if err = xml.Unmarshal(bpmnXML, &origD); err != nil {
			t.Fatalf("could not unmarshal XML into definitions, got %v", err)
		}
		t.Logf("%v", origD)
		// Now marshall back again into XML
		var newD Definition
		if bpmnXML, err = xml.Marshal(origD); err != nil {
			t.Fatalf("could not marshal XML from standard map, got %v", err)
		}

		// Now unmarshall the newly marshalled XML
		if err = xml.Unmarshal(bpmnXML, &newD); err != nil {
			t.Fatalf("could not unmarshal XML into new definition, got %v", err)
		}

		// Now see if they are the same
		compareDefinitions(t, &origD, &newD)

		t.Logf("\nRules:")
		for _, r := range origD.GetRules() {
			t.Logf("  %s:%s %s %s", r.Type, r.Id, r.Name, r.Description)
		}

		origD.BpmnIdBaseElementMap()
		// BaseElementsMap
		t.Logf("\nBaseElementsMap:%d", len(origD._BaseElementMap))
		for _, be := range origD._BaseElementMap {
			t.Logf("%s:%s (%s) oa=%d, ia=%d",
				be.GetXMLName().Local, be.GetId(), be.GetName(), len(be.GetOutgoingAssociations()), len(be.GetIncomingAssociations()))
			//t.Logf(".. %s:%v", bpmnID, be)
		}

		origD.BpmnIdGroupMap()
		t.Logf("\nGroupMap:%d", len(origD._BaseElementGroup))
		for bpmnID, group := range origD._BaseElementGroup {
			t.Logf(".. %s belongs to group %s\n", origD._BaseElementMap[bpmnID].ToString(), group.ToString())
		}

		origD.BpmnIdParentMap()
		t.Logf("\nParentMap:%d", len(origD._BaseElementParent))
		for bpmnID, parent := range origD._BaseElementParent {
			t.Logf(".. '%s', %s has parent %s\n", bpmnID, origD._BaseElementMap[bpmnID].ToString(), parent.ToString())
		}
		t.Log(origD._BaseElementParent["_6-61"].ToString())

		//testStringGraph(t)

		// Process Flow Hierarchy
		for _, p := range origD.Processes {
			t.Logf("\nProcess FLow:%s", p.ToString())
			//t.Logf("Sequence Flows:%v", p.SequenceFlows)
			for _, be := range p.Flow(false) {
				var rules []string
				for _, r := range be.GetRules() {
					rules = append(rules, fmt.Sprintf("(%s) %s %s", r.Type, r.Name, r.Description))
				}

				if be.GetType() == B2SequenceFlow {
					t.Logf("  %s: %s: %s : %s : %s",
						origD.GroupName(be),
						be.GetType(),
						be.GetName(),
						strings.Join(rules, ", "),
						origD.ParentName(be))
				} else {
					t.Logf("  %s: %s: %s : %s : %s",
						origD.GroupName(be),
						be.GetType(),
						be.GetName(),
						strings.Join(rules, ", "),
						origD.ParentName(be))
				}
			}
		}

		// Check we can't find a baseElement
		be := origD.FindBaseElementById("tim")
		assert.Nil(t, be, "could not find base element should be NIL")

		break
	}
}

func compareDefinitions(t *testing.T, d1 *Definition, d2 *Definition) {
	var d1Elements, d2Elements []BaseElement
	appendD1E := func(element BaseElement) bool {
		d1Elements = append(d1Elements, element)
		return true
	}
	d1.ApplyFunctionToBaseElements(appendD1E)

	appendD2E := func(element BaseElement) bool {
		d2Elements = append(d2Elements, element)
		return true
	}
	d2.ApplyFunctionToBaseElements(appendD2E)

	for i, b1 := range d1Elements {
		b2 := d2Elements[i]
		compareBaseElements(t, b1, b2)
	}
}

func compareBaseElements(t *testing.T, b1 BaseElement, b2 BaseElement) {
	assert.Equal(t, b1.GetId(), b2.GetId())
	assert.Equal(t, b1.GetName(), b2.GetName())
	assert.Equal(t, b1.GetOutgoingAssociations(), b2.GetOutgoingAssociations())
	assert.Equal(t, b1.GetIncomingAssociations(), b2.GetIncomingAssociations())
	//var ok = "OK"
	//if b1.GetId() != b2.GetId() || b1.GetName() != b2.GetName() {
	//	ok = "NOT OK"
	//}
	//for i, oa := range b1.GetOutgoingAssociations() {
	//	if oa != b2.GetOutgoingAssociations()[i] {
	//		ok = "NOT OK"
	//	}
	//}
	//for i, ia := range b1.GetIncomingAssociations() {
	//	if ia != b2.GetIncomingAssociations()[i] {
	//		ok = "NOT OK"
	//	}
	//}
	//t.Logf("%s --> %s", b1.ToString(), ok)

	for i, r1 := range b1.GetRules() {
		r2 := b2.GetRules()[i]
		assert.Equal(t, r1.Id, r2.Id)
		assert.Equal(t, r1.Name, r2.Name)
		assert.Equal(t, r1.Type, r2.Type)
		assert.Equal(t, r1.Description, r2.Description)
		//ok = "OK"
		//if r1.Id != r2.Id || r1.Name != r2.Name || r1.Type != r2.Type || r1.Description != r2.Description {
		//	ok = "NOT OK"
		//}
		//t.Logf(".. %s rule %s (%s) --> %s", r1.Type, r1.Id, r1.Name, ok)
	}

}

func displayDefinitions(t *testing.T, d *Definition) {
	for _, p := range d.Processes {
		t.Logf("Process definition %s", p.Id)
		for _, task := range p.Tasks {
			t.Logf("..Task definition %s (%s)", task.Id, task.Name)
		}
		for _, sp := range p.SubProcesses {
			t.Logf(".. Sub Process definition %s (%s)", sp.Id, sp.Name)
			for _, task := range sp.Tasks {
				t.Logf("   ..Task definition %s (%s)", task.Id, task.Name)
			}
		}
	}
}

func testStringGraph(t *testing.T) {
	g := &stringGraph{
		edges:    make(map[string][]string),
		vertices: []string{"5", "4", "2", "3", "1", "0"},
	}

	g.addEdge("5", "2")
	g.addEdge("5", "0")
	g.addEdge("4", "0")
	g.addEdge("4", "1")
	g.addEdge("2", "3")
	g.addEdge("3", "1")
	expected := []string{"5", "4", "0", "2", "3", "1"}

	result := g.topologicalSort()
	//result, _ := g.kahn()

	t.Logf("Got %v, wanted %v", result, expected)

	//for _, v := range result {
	//	t.Logf("%s\n", v)
	//}
}
