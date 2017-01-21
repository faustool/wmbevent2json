package wmbevent2json

import (
	"encoding/xml"
	"bytes"
	"github.com/fausto/jsonenc"
	"strings"
	"io"
	"regexp"
)

const WMB_XML_NS = "http://www.ibm.com/xmlns/prod/websphere/messagebroker/6.1.0/monitoring/event"
const SIMPLE_CONTENT = "simpleContent"
const COMPLEX_CONTENT = "complexContent"
const DATA_ELEMENT = "data"
const ATTRIBUTE_SEPARATOR = "_"
const TEXT_VALUE_PREFIX = "#text"
const ATTRIBUTE_NAME_PREFIX = "@"

var trimmer = NewAllStringTrimmer()

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

// A data structure that holds a Json Stream and the buffer used in it
type bufferedJsonStream struct {
	buffer *bytes.Buffer
	stream *jsonenc.Stream
}

// Initializes a bufferedJsonStream
func newBufferedJsonStream() *bufferedJsonStream {
	bufferedJsonStream := bufferedJsonStream{}
	bufferedJsonStream.buffer = newEmptyBuffer()
	bufferedJsonStream.stream = jsonenc.NewStream(bufferedJsonStream.buffer)
	return &bufferedJsonStream
}

// Creates a new empty buffer
func newEmptyBuffer() *bytes.Buffer {
	return bytes.NewBuffer(make([]byte, 0))
}

// A data structure that holds separate Json Streams for the whole event Json object and for simple and complex
// streams so they can be built separately
type jsonStreams struct{
	event *bufferedJsonStream
	simple *bufferedJsonStream
	complex *bufferedJsonStream
}

// Transforms a WMBEvent XML string into a Json object
func TransformWMBEventXMLToJson(wmbEventXML string) ([]byte, error) {
	d := xml.NewDecoder(strings.NewReader(wmbEventXML))

	json := jsonStreams{
		event:newBufferedJsonStream(),
		simple:newBufferedJsonStream(),
		complex:newBufferedJsonStream(),
	}

	currentWmbElementName := ""
	json.event.stream.WriteStartObject()
	for t, tokenErr := d.Token(); tokenErr != io.EOF; t, tokenErr = d.Token() {
		if tokenErr != nil {
			return nil, tokenErr
		}
		switch t := t.(type) {
		case xml.StartElement:
			currentWmbElementName = handleStartElement(t, currentWmbElementName, &json)
		case xml.EndElement:
			currentWmbElementName = handleEndElement(currentWmbElementName, &json, t)
		case xml.CharData:
			handleElementValue(t, currentWmbElementName, &json)
		}
	}

	addSimpleContent(&json)
	addComplexContent(&json)

	json.event.stream.WriteEndObject()

	return json.event.buffer.Bytes(), nil
}

// Add the value of complex Json Stream to the event Json stream as a Json array of objects
func addComplexContent(json *jsonStreams) {
	json.event.stream.WriteStartArrayWithName(COMPLEX_CONTENT)
	json.event.stream.WriteLiteralValue(string(json.complex.buffer.Bytes()))
	json.event.stream.WriteEndArray()
}

// Add the value of simple Json Stream to the event Json stream as a Json array of objects
func addSimpleContent(json *jsonStreams) {
	json.event.stream.WriteStartArrayWithName(SIMPLE_CONTENT)
	json.event.stream.WriteLiteralValue(string(json.simple.buffer.Bytes()))
	json.event.stream.WriteEndArray()
}

// Captures the XML element value and appends either to the complex or event stream
func handleElementValue(t xml.CharData, currentWmbElementName string, json *jsonStreams) {
	value := trimmer.Trim(string(t))
	if value != "" {
		if currentWmbElementName == COMPLEX_CONTENT {
			json.complex.stream.WriteNameValueString(TEXT_VALUE_PREFIX, value)
		} else {
			json.event.stream.WriteNameValueString(currentWmbElementName, value)
		}
	}
}

// Closes the current object on the complex stream and resets the currentWmbElement name (returned)
// when it is finished
func handleEndElement(currentWmbElementName string, json *jsonStreams, t xml.EndElement) string {
	if currentWmbElementName == COMPLEX_CONTENT {
		json.complex.stream.WriteEndObject()
		if t.Name.Space == WMB_XML_NS {
			// closing the "data" object
			json.complex.stream.WriteEndObject()
			// cleanup the variable to avoid multiple closes
			currentWmbElementName = ""
		}
	}
	return currentWmbElementName
}

// Append attributes to the correct stream and sets the currentWmbElementName (returned) to the appropriate value
func handleStartElement(t xml.StartElement, currentWmbElementName string, json *jsonStreams) string {
	if t.Name.Space == WMB_XML_NS {
		currentWmbElementName = t.Name.Local
		switch currentWmbElementName {
		case SIMPLE_CONTENT:
			json.simple.stream.WriteStartObject()
			writeAttributes(t, json.simple.stream, "")
			json.simple.stream.WriteEndObject()
		case COMPLEX_CONTENT:
			json.complex.stream.WriteStartObject()
			writeAttributes(t, json.complex.stream, "")
			json.complex.stream.WriteStartObjectWithName(DATA_ELEMENT)
		default:
			writeAttributes(t, json.event.stream, currentWmbElementName+ATTRIBUTE_SEPARATOR)
		}
	} else {
		name := createNamespaceQualifiedString(t)
		json.complex.stream.WriteStartObjectWithName(name)
		writeAttributes(t, json.complex.stream, ATTRIBUTE_NAME_PREFIX)
	}
	return currentWmbElementName
}

// Creates a string in the form of {namespace}:local
func createNamespaceQualifiedString(t xml.StartElement) (string) {
	space := "{" + t.Name.Space + "}"
	name := space + ":" + t.Name.Local
	return name
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
