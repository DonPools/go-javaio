package javaio

import (
	"bytes"
	"log"
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
	log.Println(dec.ReadObject())
}
