package storage

import (
	"encoding/binary"
	"io"
)

type TagCategory struct {
	Id   int32
	Name string
}

func (tc *TagCategory) Write(writer io.Writer) error {
	if err := binary.Write(writer, binary.LittleEndian, tc.Id); err != nil {
		return err
	}
	return WriteString(writer, tc.Name)
}

func (tc *TagCategory) Read(reader io.Reader) error {
	if err := binary.Read(reader, binary.LittleEndian, &tc.Id); err != nil {
		return err
	}
	return ReadString(reader, &tc.Name)
}
