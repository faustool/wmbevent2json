package wmbevent2json

import
(
	"fmt"
	"testing"
	"path/filepath"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

func TestTransform(t *testing.T) {
	wmbEventXMLFile, err := GetFile("test/WMBEvent.xml")
	if err != nil {
		t.Fatal(err)
	}
	wmbEventXML, err := ioutil.ReadFile(wmbEventXMLFile)
	if err != nil {
		t.Fatal(err)
	}
	expectedFile, err := GetFile("test/event.json")
	if err != nil {
		t.Fatal(err)
	}
	expectedBytes, err := ioutil.ReadFile(expectedFile)
	if err != nil {
		t.Fatal(err)
	}
	var expected map[string]interface{}
	err = json.Unmarshal(expectedBytes, &expected)
	if err != nil {
		t.Fatal(err)
	}

	xmlString := string(wmbEventXML)
	actualBytes, err := TransformWMBEventXMLToJson(xmlString)
	if err != nil {
		t.Fatal(err)
	}

	println(string(actualBytes))

	var actual map[string]interface{}
	err = json.Unmarshal(actualBytes, &actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, actual, ":(")
}

func GetFile(filePath string) (string, error) {
	files, err := filepath.Glob(filePath)
	if err != nil {
		return "", err
	} else {
		numberOfFiles := len(files)
		if numberOfFiles == 0 {
			return "", fmt.Errorf("%s not found", filePath)
		} else if numberOfFiles > 1 {
			return "", fmt.Errorf("More than one %s was found, which is awkward", filePath)
		}
	}
	return files[0], nil
}
