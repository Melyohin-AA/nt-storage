package storage

import (
	"encoding/binary"
	"fmt"
	"io"
)

type TagType byte

const TagTypeCount = 5

const (
	TagTypeNone   TagType = 0
	TagTypeBool   TagType = 1
	TagTypeInt    TagType = 2
	TagTypeFloat  TagType = 3
	TagTypeString TagType = 4
)

type Tag struct {
	Id    int32   `json:"id"`
	CatId int32   `json:"catId"`
	Label string  `json:"label"`
	Type  TagType `json:"type"`
}

func (tag *Tag) Write(writer io.Writer) error {
	if err := binary.Write(writer, binary.LittleEndian, tag.Id); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, tag.CatId); err != nil {
		return err
	}
	if err := WriteString(writer, tag.Label); err != nil {
		return err
	}
	return binary.Write(writer, binary.LittleEndian, tag.Type)
}

func (tag *Tag) Read(reader io.Reader) error {
	if err := binary.Read(reader, binary.LittleEndian, &tag.Id); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &tag.CatId); err != nil {
		return err
	}
	if err := ReadString(reader, &tag.Label); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &tag.Type); err != nil {
		return err
	}
	if tag.Type >= TagTypeCount {
		return fmt.Errorf("unknown tag type %d", tag.Type)
	}
	return nil
}
