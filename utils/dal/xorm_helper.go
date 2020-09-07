package dal

import (
	"encoding/json"
	"reflect"
)

/*
若要使用xorm存储和查询jsonb类型的数据，需要使用该type
xorm规定，若要使用自定义的数据类型，需要实现以下方法
FromDB ToDB Unmarshal Marshal String
*/

type Jsonb []byte

func (j *Jsonb) FromDB(bytes []byte) error {
	*j = Jsonb(bytes)
	return nil
}

func (j *Jsonb) ToDB() ([]byte, error) {
	return []byte(*j), nil
}

func (j *Jsonb) Marshal(v interface{}) (err error) {
	*j, err = json.Marshal(v)
	return
}

func (j *Jsonb) Unmarshal(v interface{}) error {
	return json.Unmarshal([]byte(*j), v)
}

func (j *Jsonb) String() string {
	return string([]byte(*j))
}

func NewJsonb(v interface{}) (Jsonb, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return Jsonb(b), nil
}

func ExtractQueryFilter(v interface{}) map[string]interface{} {
	var result = make(map[string]interface{})

	f := func(src reflect.Value) bool {
		switch src.Kind() {
		case reflect.String:
			if src.String() == "" {
				return false
			} else {
				return true
			}
		case reflect.Int:
			if src.Int() > 0 {
				return true
			} else {
				return false
			}
		case reflect.Float64:
			if src.Float() == 0 {
				return false
			} else {
				return true
			}
		default:
			return false
		}
	}
	vx := reflect.ValueOf(v)
	for i := 0; i < vx.Type().NumField(); i++ {
		key := vx.Type().Field(i).Tag.Get("json")
		if key != "" {
			field := vx.Field(i)
			if f(field) {
				result[key] = field.Interface()
			}
		}
	}
	return result
}

func GetMapKeys(src map[string]interface{}) (result []string) {
	for k := range src {
		result = append(result, k)
	}
	return
}
