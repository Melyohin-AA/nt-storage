package storage

import (
	"crypto/aes"
	"encoding/binary"
	"io"
)

type Blob struct {
	Fid       string
	Key       []byte
	Name      string     `json:"name"`
	ModTime   int64      `json:"mtime"`
	Comment   string     `json:"comment"`
	TagValues []TagValue `json:"tags"`
}

func (blob *Blob) Write(writer io.Writer) error {
	if err := WriteString(writer, blob.Fid); err != nil {
		return err
	}
	if _, err := writer.Write(blob.Key); err != nil {
		return err
	}
	if err := WriteString(writer, blob.Name); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, blob.ModTime); err != nil {
		return err
	}
	if err := WriteString(writer, blob.Comment); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, int32(len(blob.TagValues))); err != nil {
		return err
	}
	for _, tv := range blob.TagValues {
		if err := tv.Write(writer); err != nil {
			return err
		}
	}
	return nil
}

func (blob *Blob) Read(reader io.Reader, tags map[int32]Tag) error {
	if err := ReadString(reader, &blob.Fid); err != nil {
		return err
	}
	blob.Key = make([]byte, aes.BlockSize)
	if err := ReadFixed(reader, blob.Key); err != nil {
		return err
	}
	if err := ReadString(reader, &blob.Name); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &blob.ModTime); err != nil {
		return err
	}
	if err := ReadString(reader, &blob.Comment); err != nil {
		return err
	}
	var tvCount int32
	if err := binary.Read(reader, binary.LittleEndian, &tvCount); err != nil {
		return err
	}
	blob.TagValues = make([]TagValue, tvCount)
	for i := int32(0); i < tvCount; i++ {
		if err := blob.TagValues[i].Read(reader, tags); err != nil {
			return err
		}
	}
	return nil
}
