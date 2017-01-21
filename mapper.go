package wmbevent2json

import (
	"encoding/xml"
	"github.com/fausto/wmbevent2json/model"
	"strconv"
	"encoding/json"
)

const wmbns = "http://www.ibm.com/xmlns/prod/websphere/messagebroker/6.1.0/monitoring/event"

type Mapper interface {
	doMap(event *model.EventJson)
}

type Element struct {
	Name     xml.Name
	Attr     []xml.Attr
	CharData xml.CharData
	Encoder  json.Encoder
}

func (element Element) getAttribute(name string) string {
	for _, v := range element.Attr {
		if v.Name.Local == name {
			return v.Value
		}
	}
	return ""
}

func (element Element) isComplexContent() bool {
	return element.Name.Space == wmbns && element.Name.Local == "complexContent"
}

type Ignored struct {
}

func (mapper Ignored) doMap(event *model.EventJson) {

}

type Any struct {
	Element Element
}

func (mapper Any) doMap(event *model.EventJson) {

}

type EventData struct {
	Element Element
}

func (mapper EventData) doMap(event *model.EventJson) {
	event.EventSourceAddress = mapper.Element.getAttribute("eventSourceAddress")
	event.ProductVersion = mapper.Element.getAttribute("productVersion")
}

type EventIdentity struct {
	Element Element
}

func (mapper EventIdentity) doMap(event *model.EventJson) {
	event.EventName = mapper.Element.getAttribute("eventName")
}

type EventSequence struct {
	Element Element
}

func (mapper EventSequence) doMap(event *model.EventJson) {
	counter := mapper.Element.getAttribute("counter")
	i, err := strconv.Atoi(counter)
	if err == nil {
		event.Counter = 0
	} else {
		event.Counter = i
	}
	event.Timestamp = mapper.Element.getAttribute("creationTime")
}

type EventCorrelation struct {
	Element Element
}

func (mapper EventCorrelation) doMap(event *model.EventJson) {
	event.LocalTransactionId = mapper.Element.getAttribute("localTransactionId")
	event.ParentTransactionId = mapper.Element.getAttribute("parentTransactionId")
	event.GlobalTransactionId = mapper.Element.getAttribute("globalTransactionId")
}

type Broker struct {
	Element Element
}

func (mapper Broker) doMap(event *model.EventJson) {
	event.BrokerName = mapper.Element.getAttribute("name")
	event.HostName = mapper.Element.getAttribute("hostName")
}

type ExecutionGroup struct {
	Element Element
}

func (mapper ExecutionGroup) doMap(event *model.EventJson) {
	event.ExecutionGroupName = mapper.Element.getAttribute("name")
}

type MessageFlow struct {
	Element Element
}

func (mapper MessageFlow) doMap(event *model.EventJson) {
	event.MessageFlowName = mapper.Element.getAttribute("name")
}

type Node struct {
	Element Element
}

func (mapper Node) doMap(event *model.EventJson) {
	event.NodeDetail = mapper.Element.getAttribute("detail")
	event.NodeLabel = mapper.Element.getAttribute("nodeLabel")
	event.NodeType = mapper.Element.getAttribute("nodeType")
	event.NodeTerminal = mapper.Element.getAttribute("terminal")
}

type SimpleContent struct {
	Element Element
}

func (mapper SimpleContent) doMap(event *model.EventJson) {
	name := mapper.Element.getAttribute("name")
	value := mapper.Element.getAttribute("value")
	event.SimpleContents[name] = value
}

type Bitstream struct {
	Element Element
}

func (mapper Bitstream) doMap(event *model.EventJson) {
	event.BitstreamEncoding = mapper.Element.getAttribute("encoding")
	event.Bitstream = string(mapper.Element.CharData)
}

type ComplexContent struct {
	Element Element
}

func (mapper ComplexContent) doMap(event *model.EventJson) {
	event.ComplexContents = append(event.ComplexContents, string(mapper.Element.CharData))
}

func (element Element) getMapper() Mapper {
	if element.Name.Space == wmbns {
		switch element.Name.Local {
		case "eventData":
			return EventData{element}
		case "eventIdentity":
			return EventIdentity{element}
		case "eventSequence":
			return EventSequence{element}
		case "eventCorrelation":
			return EventCorrelation{element}
		case "broker":
			return Broker{element}
		case "executionGroup":
			return ExecutionGroup{element}
		case "messageFlow":
			return MessageFlow{element}
		case "node":
			return Node{element}
		case "simpleContent":
			return SimpleContent{element}
		case "bitstream":
			return Bitstream{element}
		case "complexContent":
			return ComplexContent{element}
		}
		return Ignored{}
	} else {
		return Any{element}
	}
}
