package storage

import (
	"encoding/binary"
	"fmt"
	"io"
)

func ReadFixed(reader io.Reader, buffer []byte) error {
	n, err := reader.Read(buffer)
	if err != nil {
		return err
	}
	if n < len(buffer) {
		return fmt.Errorf("not enough data to read: %d required, %d provided", len(buffer), n)
	}
	return nil
}

func WriteString(writer io.Writer, str string) error {
	if err := binary.Write(writer, binary.LittleEndian, int32(len(str))); err != nil {
		return err
	}
	_, err := writer.Write([]byte(str))
	return err
}

func ReadString(reader io.Reader, str *string) error {
	var size int32
	if err := binary.Read(reader, binary.LittleEndian, &size); err != nil {
		return err
	}
	buffer := make([]byte, size)
	if err := ReadFixed(reader, buffer); err != nil {
		return err
	}
	*str = string(buffer)
	return nil
}
