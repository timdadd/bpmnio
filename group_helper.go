package bpmnio

import "fmt"

// BpmnIdGroupMap builds & stores a map of BPMN IDs and the group they belong to
// This currently assumes an item only belongs to one group
func (d *Definition) BpmnIdGroupMap() map[string]*Group {
	if len(d._BaseElementGroup) != 0 {
		return d._BaseElementGroup
	}
	d.BpmnIdBaseElementMap()
	d._BaseElementGroup = make(map[string]*Group) // Element belongs to Group
	type shapeBE struct {
		shape *Shape
		group *Group
	}
	// Each diagram has a different plane
	for _, diagram := range d.BpmnDiagram {
		for _, plane := range diagram.BpmnPlane {
			// Find all the groups within the diagram / Plane
			var groupShapes []shapeBE
			for _, shape := range plane.BpmnShape {
				if e, inMap := d._BaseElementMap[shape.BpmnElement]; inMap {
					if e.GetType() == B2Group {
						group := d._BaseElementMap[shape.BpmnElement].(*Group)
						groupShapes = append(groupShapes, shapeBE{shape: shape, group: group})
						// Fix the group name
						if group.Name == "" {
							if be, ok := d._BaseElementMap[group.CategoryValueRef]; ok {
								group.Name = be.GetName()
							}
						}
					}
				} else {
					fmt.Printf("Why BPMN element not in the map %s\n", shape.BpmnElement)
				}
			}
			// Now find all the shapes that are lying within the bounds of each group
			// This only handles an item sitting within one group
			for _, shape := range plane.BpmnShape {
				// Does this element exist within a group
				for _, groupShape := range groupShapes {
					if groupShape.group.Id == shape.Id {
						continue // Ignore itself
					}
					if shapeWithinShape(*shape, *groupShape.shape) {
						d._BaseElementGroup[shape.BpmnElement] = groupShape.group
					}
				}
			}
		}
	}
	return d._BaseElementGroup
}

func (d *Definition) GroupId(be BaseElement) string {
	if group, inMap := d._BaseElementGroup[be.GetId()]; inMap {
		return group.Id
	}
	return ""
}

func (d *Definition) GroupName(be BaseElement) string {
	if group, inMap := d._BaseElementGroup[be.GetId()]; inMap {
		return group.Name
	}
	return ""
}

func (d *Definition) Group(be BaseElement) *Group {
	return d._BaseElementGroup[be.GetId()]
}
