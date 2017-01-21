package wmbevent2json

import (
	"encoding/xml"
	"bytes"
	"github.com/fausto/jsonenc"
	"strings"
	"io"
	"regexp"
)

// String trimmer
type AllStringTrimmer interface {
	// Trims all spaces of a String, left and right, including \t and \t characters
	Trim(value string) string
}

// Implementation of AllStringTrimmer using regular expressions
type RegExAllStringTrimmer struct {
	trimmerRegEx *regexp.Regexp
}

// Creates a new Trimmer using regular expressions as underlying implementation
// The current regex in use to match unwanted spaces is [\\n\\t\\r ]
func NewAllStringTrimmer() AllStringTrimmer {
	trimmerRegEx, _ := regexp.Compile("[\\n\\t\\r ]")
	return RegExAllStringTrimmer{trimmerRegEx}
}

// Uses regular expression to trim spaces
func (trimmer RegExAllStringTrimmer) Trim(value string) string {
	return trimmer.trimmerRegEx.ReplaceAllString(value, "")
}

// Transforms a WMBEvent XML string into a Json object
func Transform(wmbEventXML string) ([]byte, error) {
	d := xml.NewDecoder(strings.NewReader(wmbEventXML))

	eventBuffer := bytes.NewBuffer(make([]byte, 0))
	eventStream := jsonenc.NewStream(eventBuffer)

	simpleBuffer := bytes.NewBuffer(make([]byte, 0))
	simpleStream := jsonenc.NewStream(simpleBuffer)

	complexBuffer := bytes.NewBuffer(make([]byte, 0))
	complexStream := jsonenc.NewStream(complexBuffer)

	trimmer := NewAllStringTrimmer()

	currentWmbElementName := ""
	eventStream.WriteStartObject()
	for t, tokenErr := d.Token(); tokenErr != io.EOF; t, tokenErr = d.Token() {
		if tokenErr != nil {
			return nil, tokenErr
		}
		switch t := t.(type) {
		case xml.StartElement:
			if t.Name.Space == "http://www.ibm.com/xmlns/prod/websphere/messagebroker/6.1.0/monitoring/event" {
				currentWmbElementName = t.Name.Local
				if currentWmbElementName == "simpleContent" {
					simpleStream.WriteStartObject()
					writeAttributes(t, simpleStream, "")
					simpleStream.WriteEndObject()
				} else if currentWmbElementName == "complexContent" {
					complexStream.WriteStartObject()
					writeAttributes(t, complexStream, "")
					complexStream.WriteStartObjectWithName("data")
				} else {
					writeAttributes(t, eventStream, currentWmbElementName+ "_")
				}
			} else {
				space := "{" + t.Name.Space + "}"
				name := space + "#" + t.Name.Local
				complexStream.WriteStartObjectWithName(name)
				writeAttributes(t, complexStream, "@")
			}
		case xml.EndElement:
			if currentWmbElementName == "complexContent" {
				complexStream.WriteEndObject()
				if t.Name.Space == "http://www.ibm.com/xmlns/prod/websphere/messagebroker/6.1.0/monitoring/event" {
					// closing the "data" object
					complexStream.WriteEndObject()
					// cleanup the variable to avoid multiple closes
					currentWmbElementName = ""
				}
			}
		case xml.CharData:
			value := trimmer.Trim(string(t))
			if value != "" {
				if currentWmbElementName == "complexContent" {
					complexStream.WriteNameValueString("#text", value)
				} else {
					eventStream.WriteNameValueString(currentWmbElementName, value)
				}
			}
		}
	}

	eventStream.WriteStartArrayWithName("simpleContents")
	eventStream.WriteLiteralValue(string(simpleBuffer.Bytes()))
	eventStream.WriteEndArray()

	eventStream.WriteStartArrayWithName("complexContents")
	eventStream.WriteLiteralValue(string(complexBuffer.Bytes()))
	eventStream.WriteEndArray()

	eventStream.WriteEndObject()

	return eventBuffer.Bytes(), nil
}

// Write all XML attributes of a given element to a jsonenc.Stream with an optional prefix
//
// The first argument is the XML object that contains the element name. Only the local name will be used.
// The second argument is the jsonenc.Stream where the attributes will be written to
// The last argument is a prefix that can be added to the Json attribute name, e.g. an @ sign.
func writeAttributes(t xml.StartElement, currentStream *jsonenc.Stream, prefix string) {
	for _, attr := range t.Attr {
		if attr.Value != "" {
			currentStream.WriteNameValueString(prefix+attr.Name.Local, attr.Value)
		}
	}
}
