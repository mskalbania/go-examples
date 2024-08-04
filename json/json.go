package json

import (
	"encoding/json"
	"fmt"
)

const JSON = `{
  "name": "John",
  "surname": "Doe",
  "addresses": [
    {
      "street": "Street"
    }
  ]
}`

type Person struct {
	Name       string //if name match tag not required
	MiddleName string //defaults to zero value when not present in json - "" here
	//omitempty - exclude >zero-values< from output json
	Age       *int      `json:"age,omitempty"` //so to sometimes it is better to have pointers, zero value might be ambiguous
	Surname   string    `json:"surname"`
	Addresses []Address `json:"addresses"` //nil ptr when no in json, empty slice when empty json array
	Job       *Job      //for ptr type nil, for non-ptr inits with zero values
}
type Address struct {
	Street string `json:"street"`
}

type Job struct {
	Title string `json:"title"`
}

func UnmarshallExample() {
	var person Person
	json.Unmarshal([]byte(JSON), &person)
	fmt.Printf("%v", person)

	var asMap map[string]any
	json.Unmarshal([]byte(JSON), &asMap)
	//access this way is ugly but possible, type casting required on every step
	street := asMap["addresses"].([]interface{})[0].(map[string]interface{})["street"].(string)
	fmt.Printf("\n\n%v | %s", asMap, street)
}

func MarshallExample() {
	asMap := make(map[string]interface{})
	asMap["name"] = "John"
	asMap["addresses"] = []map[string]interface{}{
		{"street": "street1"},
		{"street": "street2"},
	}
	out, _ := json.MarshalIndent(asMap, "", "   ")
	fmt.Printf("%v", string(out))

	var person Person //takes either tag or struct parameter name to construct json fields
	out, _ = json.MarshalIndent(&person, "", "   ")
	fmt.Printf("%v", string(out))
}
