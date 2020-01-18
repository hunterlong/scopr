package scoper

import (
	"encoding/json"
	"io"
	"reflect"
	"strings"
)

var EmptyShows = true

// A Decoder reads and decodes JSON values from an input stream.
type Decoder struct {
	r       io.Reader
	buf     []byte
	err     error
}

// An Encoder writes JSON values to an output stream.
type Encoder struct {
	w          io.Writer
	err        error
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (enc *Encoder) Encode(v interface{}, scope string) error {
	jsonD, err := Marshal(v, scope)
	if err != nil {
		return err
	}
	_, err = enc.w.Write(jsonD)
	return err
}


func Marshal(data interface{}, scope string) ([]byte, error) {
	if reflect.ValueOf(data).Kind() == reflect.Slice {
		alldata := make([]map[string]interface{}, 0, 1)
		s := reflect.ValueOf(data)
		for i := 0; i < s.Len(); i++ {
			alldata = append(alldata, toSafeJson(s.Index(i).Interface(), scope))
		}
		return json.Marshal(alldata)
	}
	return json.Marshal(toSafeJson(data, scope))
}

func toSafeJson(input interface{}, scope string) map[string]interface{} {
	thisData := make(map[string]interface{})
	t := reflect.TypeOf(input)
	elem := reflect.ValueOf(input)
	d, _ := json.Marshal(input)

	var raw map[string]*json.RawMessage
	json.Unmarshal(d, &raw)

	if t.Kind() == reflect.Ptr {
		input = &input
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		tag := field.Tag.Get("scope")
		tags := strings.Split(tag, ",")

		jTags := field.Tag.Get("json")
		jsonTag := strings.Split(jTags, ",")

		if len(jsonTag) == 0 {
			continue
		}

		if jsonTag[0] == "" || jsonTag[0] == "-" {
			continue
		}

		trueValue := elem.Field(i).Interface()

		if len(jsonTag) == 2 {
			if jsonTag[1] == "omitempty" && trueValue == "" {
				continue
			}
		}

		if tag == "" && EmptyShows {
			thisData[jsonTag[0]] = trueValue
			continue
		}

		if forType(tags, scope) {
			thisData[jsonTag[0]] = trueValue
		}
	}
	return thisData
}

func forType(tags []string, scope string) bool {
	for _, v := range tags {
		if v == scope {
			return true
		}
	}
	return false
}
