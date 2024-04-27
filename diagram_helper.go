package bpmnio

//func (definition *Definition) BuildShapeBoundsMap() {
//	definition._shapeBoundsMap = map[string]Bounds{}
//	for _, diagram := range definition.BpmnDiagram {
//		for _, plane := range diagram.BpmnPlane {
//			for _, shape := range plane.BpmnShape {
//				definition._shapeBoundsMap[shape.BpmnElement] = shape.BpmnBounds
//			}
//		}
//	}
//}

//// BuildShapeHierarchyMap builds a hierarchy of shapes
//func (d *Definition) BuildShapeHierarchyMap() {
//	//if len(d._shapeBoundsMap) == 0 {
//	//	d.BuildShapeBoundsMap()
//	//}
//	d.BuildBaseElementsMap()
//	//pc := make(map[string][]Shape) // Parent has children
//	cp := make(map[*Shape]*Shape)   // Children has parent
//	gc := make(map[*Shape][]*Shape) // Group has Children
//	for _, diagram := range d.BpmnDiagram {
//		var maxArea int
//		for _, plane := range diagram.BpmnPlane {
//			for _, shape := range plane.BpmnShape {
//				shape.BpmnBounds.area = shape.BpmnBounds.Width * shape.BpmnBounds.Height
//				//shape.BpmnBounds.area = max(shape.BpmnBounds.area, -shape.BpmnBounds.area) // ABS
//				//fmt.Printf("shape: %s:%d\n", shape.Id, shape.BpmnBounds.area)
//				if shape.BpmnBounds.area > maxArea {
//					maxArea = shape.BpmnBounds.area
//				}
//				//if shape.BpmnBounds.area == 0 {
//				//	fmt.Printf("WHY ZERO?: %s (%v):%d\n", shape.Id, &shape, shape.BpmnBounds.area)
//				//}
//			}
//			//fmt.Printf("max area: %d\n", maxArea)
//			for _, shape := range plane.BpmnShape {
//				//if shape.BpmnBounds.area == 0 {
//				//	fmt.Printf("WHY ZERO?: %s (%v):%d\n", shape.Id, &shape, shape.BpmnBounds.area)
//				//}
//				// Does this element exist within another element?
//				for _, parentShape := range plane.BpmnShape {
//					if parentShape.Id == shape.Id {
//						continue
//					}
//					if shapeWithinShape(*shape, *parentShape) {
//						//fmt.Println(shape.Id, shape.BpmnElement, "within", parentShape.Id, parentShape.BpmnElement)
//						//fmt.Printf("%s within %s\n",
//						//	d._BaseElementMap[shape.BpmnElement].ToString(),
//						//	d._BaseElementMap[parentShape.BpmnElement].ToString())
//						//if parentShape.BpmnBounds.area == 0 {
//						//	fmt.Printf("Parent Area=0: %s:%v\n", d._BaseElementMap[shape.BpmnElement].ToString(), shape.BpmnBounds)
//						//}
//						//if shape.BpmnBounds.area == 0 {
//						//	fmt.Printf("Area=0: %s:%v\n", d._BaseElementMap[shape.BpmnElement].ToString(), shape.BpmnBounds)
//						//}
//
//						// If parent is a group then add anyway everything is added to the group
//						if d._BaseElementMap[parentShape.BpmnElement].GetType() == B2Group {
//							gc[parentShape] = append(gc[parentShape], shape)
//							continue // Don't make group part of any other set
//						}
//
//						if pShape, ok := cp[shape]; ok {
//							//if pShape.BpmnBounds.area == 0 {
//							//	fmt.Printf("AREA=0: %s:%d\n", d._BaseElementMap[shape.BpmnElement].ToString(), shape.BpmnBounds.area)
//							//}
//							// Associate to the parent with the smallest area
//							if pShape.BpmnBounds.area > parentShape.BpmnBounds.area {
//								cp[shape] = parentShape
//							}
//						} else {
//							cp[shape] = parentShape
//						}
//
//					}
//				}
//			}
//		}
//	}
//	d._ShapeHierarchyMap = make(map[string][]BaseElement, len(cp))
//	for childShape, parentShape := range cp {
//		d._ShapeHierarchyMap[parentShape.BpmnElement] =
//			append(d._ShapeHierarchyMap[parentShape.BpmnElement], d._BaseElementMap[childShape.BpmnElement])
//	}
//	for groupShape, shapes := range gc {
//		d._ShapeHierarchyMap[groupShape.BpmnElement] = make([]BaseElement, len(shapes))
//		for i, shape := range shapes {
//			d._ShapeHierarchyMap[groupShape.BpmnElement][i] = d._BaseElementMap[shape.BpmnElement]
//		}
//	}
//}

func shapeWithinShape(s Shape, w Shape) bool {
	return boundsWithinBounds(s.BpmnBounds, w.BpmnBounds)
}

func boundsWithinBounds(b Bounds, within Bounds) bool {
	//fmt.Printf("%v within %v --> %t, %t, %t, %t\n", b, within,
	//	b.X > within.X, b.Y > within.Y, b.X+b.Width < within.X+within.Width, b.Y+b.Height < within.Y+within.Height)
	return b.X > within.X && b.Y > within.Y && b.X+b.Width < within.X+within.Width && b.Y+b.Height < within.Y+within.Height
}
