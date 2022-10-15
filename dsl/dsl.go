package dsl

import (
	"encoding/json"
	"reflect"

	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

type DSL map[string]interface{}

func (d DSL) String() string {
	v, _ := json.Marshal(d)
	return string(v)
}

func addValueForDSL(d DSL, field string, value interface{}) {
	var v = reflect.ValueOf(value)
	if v.IsValid() && !v.IsZero() {
		d[field] = value
	}
}

type LeafValue interface{}

type valueType struct {
	mType mapping.FieldType
	aType bool // is array type
}

func WithArrayType(isArrayType bool) func(ArrayTypeNode) {
	return func(n ArrayTypeNode) {
		n.isArrayType()

	}
}

func NewValueType(mType mapping.FieldType, isArrayType bool) *valueType {
	return &valueType{
		mType: mType,
		aType: isArrayType,
	}
}

func (v *valueType) isArrayType() bool {
	return v.aType
}

func (v *valueType) setArrayType(isArrayType bool) {
	v.aType = isArrayType
}

var EmptyValue LeafValue = nil

var emptyValue = EmptyValue

var EmptyDSL = DSL{}

var emptyDSL = EmptyDSL
