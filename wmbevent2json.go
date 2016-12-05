package wmbevent2json

import "github.com/fausto/wmbevent2json/model"

func Transform(wmbEventXML string) (model.Event, error) {
	if (wmbEventXML != nil) {
		return model.Event{}, nil
	} else {
		return nil, nil
	}
}