package bpmnio

import (
	"encoding/xml"
)

type ExtensionElements struct {
	XMLName        xml.Name        `xml:"extensionElements"`
	DisplayOrder   int             `xml:"displayOrder,attr"`
	Form           *Form           `xml:"form"`
	Rules          *Rules          `xml:"rules,omitempty"`
	Implementation *Implementation `xml:"implementation,omitempty"`
}

type Form struct {
	XMLName  xml.Name `xml:"form"`
	Name     string   `xml:"name,attr"`
	Scenario string   `xml:"scenario,attr"`
}

type Implementation struct {
	XMLName       xml.Name `xml:"implementation"`
	Type          string   `xml:"type,attr"`
	Configuration string   `xml:"configuration,attr"`
}

type Rules struct {
	XMLName xml.Name `xml:"rules"`
	Rules   []*Rule  `xml:"rule,omitempty"`
}

type Rule struct {
	XMLName      xml.Name `xml:"rule"`
	DisplayOrder int      `xml:"displayOrder,attr"`
	Id           string   `xml:"id,attr"`
	Type         string   `xml:"type,attr"`
	Code         string   `xml:"code,attr"`
	Name         string   `xml:"name,attr"`
	Description  string   `xml:"description,attr"`
}
