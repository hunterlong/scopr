package scopr

import (
	"encoding/json"
	"io"
	"reflect"
	"strings"
)

type Scopr struct {
	data  interface{}
	scope string
}

func New(obj interface{}, scope string) Scopr {
	return Scopr{data: obj, scope: scope}
}

// MarshalJSON for Scopr
func (s Scopr) MarshalJSON() ([]byte, error) {
	svc := reflect.ValueOf(s.data)
	if svc.Kind() == reflect.Slice {
		alldata := make([]map[string]interface{}, 0)
		for i := 0; i < svc.Len(); i++ {
			objIndex := svc.Index(i)
			if objIndex.Kind() == reflect.Ptr {
				objIndex = objIndex.Elem()
			}
			alldata = append(alldata, SafeJson(Scopr{objIndex, s.scope}, s.scope))
		}
		return json.Marshal(alldata)
	}
	return json.Marshal(SafeJson(svc.Interface(), s.scope))
}

// MarshalJSON for Scopr
func (s *Scopr) UnmarshalJSON(data []byte) error {
	if reflect.ValueOf(s.data).Kind() == reflect.Slice {
		alldata := make([]map[string]interface{}, 0, 1)
		sv := reflect.ValueOf(data)
		for i := 0; i < sv.Len(); i++ {
			alldata = append(alldata, SafeJson(sv.Index(i).Interface(), s.scope))
		}
		s.data = alldata
		return nil
	}
	s.data = SafeJson(data, s.scope)
	return nil
}

var EmptyShows = true

// A Decoder reads and decodes JSON values from an input stream.
type Decoder struct {
	r   io.Reader
	buf []byte
	err error
}

// An Encoder writes JSON values to an output stream.
type Encoder struct {
	w   io.Writer
	err error
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func Json(obj interface{}, scope string) ([]byte, error) {
	return json.Marshal(New(obj, scope))
}

func (enc *Encoder) Encode(v interface{}, scope string) error {
	jsonD, err := json.Marshal(New(v, scope))
	if err != nil {
		return err
	}
	_, err = enc.w.Write(jsonD)
	return err
}

func (enc *Encoder) Write(v interface{}) error {
	dd, _ := json.Marshal(v)
	_, err := enc.w.Write(dd)
	return err
}

func SafeJson(input interface{}, scope string) map[string]interface{} {
	thisData := make(map[string]interface{})
	t := reflect.TypeOf(input)
	elem := reflect.ValueOf(input)
	d, _ := json.Marshal(input)

	var raw map[string]*json.RawMessage
	json.Unmarshal(d, &raw)

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
