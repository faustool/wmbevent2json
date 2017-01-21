package wmbevent2json

import (
	"encoding/xml"
	"bytes"
	"github.com/fausto/jsonenc"
	"strings"
	"io"
	"regexp"
)

type AllStringTrimmer interface {
	Trim(value string) string
}

type RegExAllStringTrimmer struct {
	trimmerRegEx *regexp.Regexp
}

func NewAllStringTrimmer() AllStringTrimmer {
	trimmerRegEx, _ := regexp.Compile("[\\n\\t ]")
	return RegExAllStringTrimmer{trimmerRegEx}
}

func (trimmer RegExAllStringTrimmer) Trim(value string) string {
	return trimmer.trimmerRegEx.ReplaceAllString(value, "")
}

func Transform(wmbEventXML string) ([]byte, error) {
	d := xml.NewDecoder(strings.NewReader(wmbEventXML))

	buffer := bytes.NewBuffer(make([]byte, 0))
	stream := jsonenc.NewStream(buffer)

	trimmer := NewAllStringTrimmer()

	stream.WriteStartObject()
	for t, tokenErr := d.Token(); tokenErr != io.EOF; t, tokenErr = d.Token() {
		if tokenErr != nil {
			return nil, tokenErr
		}
		switch t := t.(type) {
		case xml.StartElement:
			stream.WriteStartObjectWithName(t.Name.Local)
			for _, attr := range t.Attr {
				stream.WriteNameValueString("@" + attr.Name.Local, attr.Value)
			}
		case xml.EndElement:
			stream.WriteEndObject()
		case xml.CharData:
			value := trimmer.Trim(string(t))
			if value != "" {
				stream.WriteNameValueString("#value", value)
			}
		}
	}
	stream.WriteEndObject()

	return buffer.Bytes(), nil
}
