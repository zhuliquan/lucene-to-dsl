package dsl

import (
	"encoding/json"
	"reflect"

	mapping "github.com/zhuliquan/es-mapping"
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

// indicate whether is node array data type
type ArrayTypeNode interface {
	IsArrayType() bool
	SetArrayType(arrayType bool)
}

type ValueType interface {
	ArrayTypeNode
	FieldTypeNode() mapping.FieldType
}

type valueType struct {
	mType mapping.FieldType
	aType bool // is array type
}

func WithArrayType(isArrayType bool) func(ArrayTypeNode) {
	return func(n ArrayTypeNode) {
		n.SetArrayType(isArrayType)
	}
}

func NewValueType(mType mapping.FieldType, isArrayType bool) *valueType {
	return &valueType{
		mType: mType,
		aType: isArrayType,
	}
}

func (v *valueType) IsArrayType() bool {
	return v.aType
}

func (v *valueType) SetArrayType(isArrayType bool) {
	v.aType = isArrayType
}

func (v *valueType) FieldTypeNode() mapping.FieldType {
	return v.mType
}

var EmptyValue LeafValue = nil

var EmptyDSL = DSL{}
