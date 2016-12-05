package wmbevent2json

import (
	"github.com/fausto/wmbevent2json/model"
	"github.com/fausto/stack"
	//"encoding/xml"
	//"strings"
	"encoding/xml"
	"strings"
)

func Transform(wmbEventXML string) (model.EventJson, error) {
	mapperStack := stack.ArrayStack{make([]interface{}, 0)}
	d := xml.NewDecoder(strings.NewReader(wmbEventXML))
	event := model.EventJson{}
	for {
		t, tokenErr := d.Token()
		if tokenErr != nil {
			return nil, tokenErr
		}
		switch t := t.(type) {
		case xml.StartElement:
			v, err := mapperStack.Peak()
			if (err != nil && v.(Element).isComplexContent()) {

			} else {
				element := Element{Name: t.Name, Attr: t.Attr}
				mapperStack.Push(element)
			}
		case xml.EndElement:
			v, err := mapperStack.Pop()
			if (err != nil) {
				mapper := v.(Element).getMapper()
				mapper.doMap(event)
			}
		case xml.CharData:
			v, err := mapperStack.Peak()
			if (err != nil) {
				v.(Element).CharData = t
			}
		}
	}

	return model.EventJson{}, nil
}
