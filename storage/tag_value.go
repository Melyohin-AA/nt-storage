package storage

import (
	"encoding/binary"
	"fmt"
	"io"
)

type TagValue struct {
	TagId int32       `json:"tagId"`
	Value interface{} `json:"value"`
}

func (tv *TagValue) Write(writer io.Writer) error {
	if err := binary.Write(writer, binary.LittleEndian, tv.TagId); err != nil {
		return err
	}
	var err error
	if tv.Value != nil {
		if str, isString := tv.Value.(string); isString {
			err = WriteString(writer, str)
		} else {
			err = binary.Write(writer, binary.LittleEndian, tv.Value)
		}
	}
	return err
}

func (tv *TagValue) Read(reader io.Reader, tags map[int32]Tag) error {
	if err := binary.Read(reader, binary.LittleEndian, &tv.TagId); err != nil {
		return err
	}
	tag, exists := tags[tv.TagId]
	if !exists {
		return fmt.Errorf("unknown tag id %d", tv.TagId)
	}
	var err error
	switch tag.Type {
	case TagTypeNone:
		tv.Value = nil
	case TagTypeBool:
		var value bool
		err = binary.Read(reader, binary.LittleEndian, &value)
		tv.Value = value
	case TagTypeInt:
		var value int64
		err = binary.Read(reader, binary.LittleEndian, &value)
		tv.Value = value
	case TagTypeFloat:
		var value float64
		err = binary.Read(reader, binary.LittleEndian, &value)
		tv.Value = value
	case TagTypeString:
		var value string
		err = ReadString(reader, &value)
		tv.Value = value
	}
	return err
}
