package mft_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/t9t/gomft/mft"
)

func TestParseRecordHeader(t *testing.T) {
	b := readTestMft(t)
	header := mft.ParseRecordHeader(b)
	expected := mft.RecordHeader{
		Signature:             []byte{'F', 'I', 'L', 'E'},
		UpdateSequenceOffset:  48,
		UpdateSequenceSize:    3,
		LogFileSequenceNumber: 25695988020,
		RecordUsageNumber:     145,
		HardLinkCount:         1,
		FirstAttributeOffset:  56,
		Flags:                 []byte{0x01, 0x00},
		ActualSize:            480,
		AllocatedSize:         1024,
		BaseRecordReference:   []byte{0xA0, 0xB0, 0xC0, 0xD0, 0xE0, 0xF0, 0x10, 0x90},
		NextAttributeId:       8,
	}

	assert.Equal(t, expected, header)
}

func TestParseAttributes(t *testing.T) {
	b := readTestMft(t)
	attributeData := b[56:]
	attributes, err := mft.ParseAttributes(attributeData)
	require.Nilf(t, err, "error parsing attributes: %v", err)

	expectedAttributes := []mft.Attribute{
		mft.Attribute{Type: 16, Resident: true, Name: "", Flags: []byte{0, 0}, AttributeId: 0, Data: []byte{0x94, 0xF0, 0x48, 0x96, 0x5B, 0x2F, 0xCC, 0x1, 0x94, 0xF0, 0x48, 0x96, 0x5B, 0x2F, 0xCC, 0x1, 0x94, 0xF0, 0x48, 0x96, 0x5B, 0x2F, 0xCC, 0x1, 0x94, 0xF0, 0x48, 0x96, 0x5B, 0x2F, 0xCC, 0x1, 0x6, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}},
		mft.Attribute{Type: 48, Resident: true, Name: "", Flags: []byte{0, 0}, AttributeId: 3, Data: []byte{0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5, 0x0, 0x94, 0xF0, 0x48, 0x96, 0x5B, 0x2F, 0xCC, 0x1, 0x94, 0xF0, 0x48, 0x96, 0x5B, 0x2F, 0xCC, 0x1, 0x94, 0xF0, 0x48, 0x96, 0x5B, 0x2F, 0xCC, 0x1, 0x94, 0xF0, 0x48, 0x96, 0x5B, 0x2F, 0xCC, 0x1, 0x0, 0x0, 0xBC, 0x39, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xBC, 0x39, 0x0, 0x0, 0x0, 0x0, 0x6, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0x3, 0x24, 0x0, 0x4D, 0x0, 0x46, 0x0, 0x54, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}},
		mft.Attribute{Type: 128, Resident: false, Name: "", Flags: []byte{0, 0}, AttributeId: 1, Data: []byte{0x33, 0x20, 0xC8, 0x0, 0x0, 0x0, 0xC, 0x43, 0x22, 0xB5, 0x0, 0xBA, 0x5, 0x5C, 0x3, 0x43, 0x81, 0xDE, 0x0, 0x65, 0xCF, 0x47, 0x4, 0x43, 0x84, 0xB3, 0x0, 0x5D, 0x8B, 0xEF, 0x9, 0x43, 0xB0, 0xE1, 0x0, 0x90, 0xB4, 0xB5, 0x18, 0x43, 0x0, 0xC8, 0x0, 0xF4, 0xEA, 0x13, 0x1, 0x43, 0x6, 0xC8, 0x0, 0x9A, 0x3A, 0x5A, 0xFE, 0x43, 0x12, 0xC8, 0x0, 0xF4, 0x7, 0x4D, 0xFE, 0x33, 0xF, 0xC8, 0x0, 0x23, 0xD4, 0xC0, 0x42, 0x62, 0x16, 0x54, 0x2, 0x95, 0x3, 0x0, 0x0, 0x0}},
		mft.Attribute{Type: 176, Resident: false, Name: "", Flags: []byte{0, 0}, AttributeId: 7, Data: []byte{0x41, 0x3A, 0xBE, 0x84, 0x83, 0x0, 0x0, 0x0}},
	}

	assert.Equal(t, expectedAttributes, attributes)
}

func readTestMft(t *testing.T) []byte {
	b, err := ioutil.ReadFile("test-mft.bin")
	require.Nilf(t, err, "unable to read test-mft.bin: %v", err)
	return b
}
