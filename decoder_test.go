package javaio

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDecoder(t *testing.T) {
	_, err := NewDecoder(bytes.NewReader([]byte{0xac, 0xed, 0x00, 0x05}))
	assert.NoError(t, err)

	_, err = NewDecoder(bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x05}))
	assert.Error(t, err)

	_, err = NewDecoder(bytes.NewReader([]byte{0xac, 0xed, 0x00, 0x00}))
	assert.Error(t, err)
}

func TestDecoder_ReadObject(t *testing.T) {
	r := bytes.NewReader([]byte{
		0xac, 0xed, 0x00, 0x05, 0x73, 0x72, 0x00, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x01, 0x02, 0x00, 0x02, 0x49, 0x00, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x4c,
		0x00, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x74, 0x00, 0x06, 0x4c, 0x4c, 0x69, 0x73, 0x74, 0x3b, 0x78,
		0x70, 0x00, 0x00, 0x00, 0x11, 0x73, 0x71, 0x00, 0x7e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x13, 0x70,
		0x71, 0x00, 0x7e, 0x00, 0x03,
	})
	dec, err := NewDecoder(r)
	assert.NoError(t, err)
	dec.RegisterType("List", reflect.TypeOf(List{}))

	obj1, err := dec.ReadObject()
	assert.NoError(t, err)
	obj2, err := dec.ReadObject()
	assert.NoError(t, err)

	list2 := &List{
		Value: 19,
	}
	list1 := &List{
		Value: 17,
		Next:  list2,
	}

	assert.Equal(t, list1, obj1)
	assert.Equal(t, list2, obj2)
}

type ClassA struct {
	Hello *String
}

type ClassB struct {
	super ClassA
	value int32
}

func (ClassA) ClassName() string {
	return "A"
}

func (ClassB) ClassName() string {
	return "B"
}

func (b *ClassB) Super() interface{} {
	return &b.super
}

func (b *ClassB) ReadObject(dec *Decoder) error {
	err := dec.ReadBinary(&b.value)
	if err != nil {
		return err
	}
	return nil
}

func TestDecoder_ReadObjectWithCustomReadObject(t *testing.T) {
	r := bytes.NewReader([]byte{
		0xAC, 0xED, 0x00, 0x05, 0x73, 0x72, 0x00, 0x01, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x01, 0x03, 0x00, 0x00, 0x78, 0x72, 0x00, 0x01, 0x41, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x01, 0x02, 0x00, 0x01, 0x4C, 0x00, 0x05, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x74, 0x00, 0x12, 0x4C,
		0x6A, 0x61, 0x76, 0x61, 0x2F, 0x6C, 0x61, 0x6E, 0x67, 0x2F, 0x53, 0x74, 0x72, 0x69, 0x6E, 0x67,
		0x3B, 0x78, 0x70, 0x74, 0x00, 0x05, 0x77, 0x6F, 0x72, 0x6C, 0x64, 0x77, 0x04, 0x00, 0x00, 0x00,
		0x01, 0x78,
	})
	dec, err := NewDecoder(r)
	assert.NoError(t, err)
	dec.RegisterType("B", reflect.TypeOf(ClassB{}))
	object, err := dec.ReadObject()
	assert.NoError(t, err)
	b := object.(*ClassB)
	assert.Equal(t, &String{Value: "world"}, b.super.Hello)
	assert.Equal(t, int32(1), b.value)
}
