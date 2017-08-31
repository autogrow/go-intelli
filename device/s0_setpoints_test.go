package device

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func Test_myData(t *testing.T) {

	myBytes := []byte("1234567890")
	myString := string(myBytes[0])
	fmt.Print(myString)

	// should be equal
	if myString != "1" {
		t.Errorf("Error: myString = %s", myString)
	}

	r := json.Unmarshal([]byte("{}"), SetPointIClimate{})

	fmt.Printf("%s", r)
}

func Test_convert(t *testing.T) {
	var resultDataJSON []byte
	var err error

	d1Response := make([]byte, 64)
	testData := SetPointIClimate{
		LightBank:     "1",
		LightOn:       1201,
		LightDuration: 184,
		DayTemp:       21.5,
		NightDropDeg:  3.1,
		RhDay:         1,
		RhMax:         3,
		RhNight:       256,
		CO2:           34,
	}

	fillSetPointData(&d1Response, testData)

	testResultData := extractSetPointArray(d1Response)

	resultDataJSON, err = json.Marshal(testResultData)

	if err == nil {
		fmt.Printf(" %s  ", err)
	}

	if !reflect.DeepEqual(testData, testResultData[0]) {

		t.Log("Content is not equal")
		t.Log("Given :")
		t.Log(string(resultDataJSON))

		t.Log("Result:")
		t.Log(d1Response)

		t.Log("Result:")
		t.Log(string(resultDataJSON))

		t.FailNow()

	}

}
