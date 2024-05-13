package bpmnio

// ApplyFunctionToShapes applies a function to all shapes within Definition
func (d *Definition) ApplyFunctionToShapes(f func(shape *Shape)) {
	if d == nil {
		return
	}
	// Each diagram has a different plane
	for _, diagram := range d.BpmnDiagram {
		for _, plane := range diagram.BpmnPlane {
			for _, shape := range plane.BpmnShape {
				f(shape)
			}
		}
	}
}

// BpmnIdShapeMap stores and returns a mapping of Shape by BPMN_ID if it hasn't already been created
func (d *Definition) BpmnIdShapeMap() map[string]*Shape {
	if len(d._BpmnIdShapeMap) == 0 {
		shapes := d.FindShapes()
		d._BpmnIdShapeMap = make(map[string]*Shape, len(shapes))
		for _, s := range shapes {
			d._BpmnIdShapeMap[s.BpmnElement] = s
		}
	}
	return d._BpmnIdShapeMap
}

// FindShapes finds all shapes in the definition
func (d *Definition) FindShapes() (shapes []*Shape) {
	fAppendShape := func(shape *Shape) {
		shapes = append(shapes, shape)
	}
	d.ApplyFunctionToShapes(fAppendShape)
	return
}

// FindShapeByElementId finds the Shape that matches the id
func (d *Definition) FindShapeByElementId(bpmnId string) (shape *Shape) {
	d.BpmnIdShapeMap()
	return d._BpmnIdShapeMap[bpmnId]
}

// ShapeOfBaseElement finds the Shape that matches the base element
func (d *Definition) ShapeOfBaseElement(be BaseElement) (shape *Shape) {
	d.BpmnIdShapeMap()
	return d._BpmnIdShapeMap[be.GetId()]
}

// FindShapesByElementTypes finds all Shapes that have one of the specific Types
// Interfaces are just pointers and types so no need for *
func (d *Definition) FindShapesByElementTypes(types ...ElementType) (shapes []*Shape) {
	// Create a function that checks the type and appends if found
	var matchType = make(map[ElementType]bool, len(types))
	for _, et := range types {
		matchType[et] = true
	}
	d.BpmnIdBaseElementMap() // Create the map of elements
	appendType := func(shape *Shape) {
		if element := d._BaseElementMap[shape.BpmnElement]; element != nil {
			if matchType[element.GetType()] {
				shapes = append(shapes, shape)
			}
		}
	}
	d.ApplyFunctionToShapes(appendType)
	return
}

func shapeWithinShape(s Shape, w Shape) bool {
	return boundsWithinBounds(s.BpmnBounds, w.BpmnBounds)
}

func boundsWithinBounds(b Bounds, within Bounds) bool {
	//fmt.Printf("%v within %v --> %t, %t, %t, %t\n", b, within,
	//	b.X > within.X, b.Y > within.Y, b.X+b.Width < within.X+within.Width, b.Y+b.Height < within.Y+within.Height)
	return b.X > within.X && b.Y > within.Y && b.X+b.Width < within.X+within.Width && b.Y+b.Height < within.Y+within.Height
}
