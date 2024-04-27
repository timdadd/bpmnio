package BPMNIO

import (
	"encoding/xml"
)

// ProcessInfoFromXML takes the xmlByte data and creates a ProcessInfo struct
func definitionsFromXML(xmlData []byte) (definitions Definition, err error) {
	err = xml.Unmarshal(xmlData, &definitions)
	return
}
