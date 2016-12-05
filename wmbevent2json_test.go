package wmbevent2json

import (
	"fmt"
	"testing"
	"path/filepath"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"github.com/fausto/wmbevent2json/model"
)

func TestTransform(t *testing.T) {

	wmbEventXMLFile, err := GetFile("test/WMBEvent.xml")
	if (err != nil) {
		t.Fatal(err)
	}

	expectedFile, err := GetFile("test/event.json")
	if (err != nil) {
		t.Fatal(err)
	}

	wmbEventXML, err := ioutil.ReadFile(wmbEventXMLFile)
	if (err != nil) {
		t.Fatal(err)
	}

	expectedBytes, err := ioutil.ReadFile(expectedFile)
	if (err != nil) {
		t.Fatal(err)
	}

	expected := model.Event{}
	err = json.Unmarshal(expectedBytes, &expected)
	if (err != nil) {
		t.Fatal(err)
	}

	actual, err := Transform(string(wmbEventXML))

	assert.Equal(t, expected, actual, "Unexpected events")
}

func GetFile(filePath string) (string, error) {
	files, err := filepath.Glob(filePath);

	if (err != nil) {
		return "", err
	} else {
		numberOfFiles := len(files)
		if (numberOfFiles == 0) {
			return "", fmt.Errorf("%s not found", filePath)
		} else if (numberOfFiles > 1) {
			return "", fmt.Errorf("More than one %s was found, which is awkward", filePath)
		}
	}

	return files[0], nil
}
