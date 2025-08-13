package storage

import (
	"encoding/binary"
	"fmt"
	"io"
	"sync"
)

const DataFileVersion = int32(1)

type DataFile struct {
	mx      sync.RWMutex
	tagCats map[int32]TagCategory
	tags    map[int32]Tag
	blobs   map[string]*Blob
}

func NewDataFile() *DataFile {
	return &DataFile{
		tagCats: make(map[int32]TagCategory),
		tags:    make(map[int32]Tag),
		blobs:   make(map[string]*Blob),
	}
}

// tag cats

func (df *DataFile) GetTagCat(id int32) (TagCategory, bool) {
	df.mx.RLock()
	defer df.mx.RUnlock()
	tc, exists := df.tagCats[id]
	return tc, exists
}

func (df *DataFile) TagCatCount() int {
	df.mx.RLock()
	defer df.mx.RUnlock()
	return len(df.tagCats)
}

func (df *DataFile) TagCats() map[int32]TagCategory {
	return df.tagCats
}

func (df *DataFile) AddTagCat(tagCat TagCategory) (TagCategory, error) {
	df.mx.Lock()
	defer df.mx.Unlock()
	tagCat.Id = int32(len(df.tagCats))
	if _, exists := df.tagCats[tagCat.Id]; exists {
		return tagCat, fmt.Errorf("tag category with id %d already exists", tagCat.Id)
	}
	df.tagCats[tagCat.Id] = tagCat
	return tagCat, nil
}

// tags

func (df *DataFile) GetTag(id int32) (Tag, bool) {
	df.mx.RLock()
	defer df.mx.RUnlock()
	tag, exists := df.tags[id]
	return tag, exists
}

func (df *DataFile) TagCount() int {
	df.mx.RLock()
	defer df.mx.RUnlock()
	return len(df.tags)
}

func (df *DataFile) Tags() map[int32]Tag {
	return df.tags
}

func (df *DataFile) AddTag(tag Tag) (Tag, error) {
	df.mx.Lock()
	defer df.mx.Unlock()
	tag.Id = int32(len(df.tags))
	if _, exists := df.tags[tag.Id]; exists {
		return tag, fmt.Errorf("tag with id %d already exists", tag.Id)
	}
	if _, exists := df.tagCats[tag.CatId]; !exists {
		return tag, fmt.Errorf("tag category with id %d does not exist", tag.CatId)
	}
	df.tags[tag.Id] = tag
	return tag, nil
}

// blobs

func (df *DataFile) GetBlob(fid string) *Blob {
	df.mx.RLock()
	defer df.mx.RUnlock()
	return df.blobs[fid]
}

func (df *DataFile) BlobCount() int {
	df.mx.RLock()
	defer df.mx.RUnlock()
	return len(df.blobs)
}

func (df *DataFile) EnumerateBlobs(f func(*Blob) bool) {
	df.mx.RLock()
	defer df.mx.RUnlock()
	for _, blob := range df.blobs {
		if !f(blob) {
			break
		}
	}
}

func (df *DataFile) AddBlob(blob *Blob) error {
	df.mx.Lock()
	defer df.mx.Unlock()
	if _, exists := df.blobs[blob.Fid]; exists {
		return fmt.Errorf("blob with fid '%s' already exists", blob.Fid)
	}
	df.blobs[blob.Fid] = blob
	return nil
}

func (df *DataFile) RemoveBlob(fid string) {
	df.mx.Lock()
	defer df.mx.Unlock()
	delete(df.blobs, fid)
}

// io

func (df *DataFile) Write(writer io.Writer) error {
	df.mx.RLock()
	defer df.mx.RUnlock()
	// Version
	if err := binary.Write(writer, binary.LittleEndian, DataFileVersion); err != nil {
		return err
	}
	// Tag categories
	if err := binary.Write(writer, binary.LittleEndian, int32(len(df.tagCats))); err != nil {
		return err
	}
	for _, tc := range df.tagCats {
		if err := tc.Write(writer); err != nil {
			return err
		}
	}
	// Tags
	if err := binary.Write(writer, binary.LittleEndian, int32(len(df.tags))); err != nil {
		return err
	}
	for _, tag := range df.tags {
		if err := tag.Write(writer); err != nil {
			return err
		}
	}
	// Blobs
	if err := binary.Write(writer, binary.LittleEndian, int32(len(df.blobs))); err != nil {
		return err
	}
	for _, blob := range df.blobs {
		if err := blob.Write(writer); err != nil {
			return err
		}
	}
	return nil
}

func (df *DataFile) Read(reader io.Reader) error {
	df.mx.Lock()
	defer df.mx.Unlock()
	// Version
	var version int32
	if err := binary.Read(reader, binary.LittleEndian, &version); err != nil {
		return err
	}
	if version != DataFileVersion {
		return fmt.Errorf("unsupported version %d of data file, %d needed", version, DataFileVersion)
	}
	// Tag categories
	var tagCatCount int32
	if err := binary.Read(reader, binary.LittleEndian, &tagCatCount); err != nil {
		return err
	}
	df.tagCats = make(map[int32]TagCategory)
	for i := int32(0); i < tagCatCount; i++ {
		var tc TagCategory
		if err := tc.Read(reader); err != nil {
			return err
		}
		df.tagCats[tc.Id] = tc
	}
	// Tags
	var tagCount int32
	if err := binary.Read(reader, binary.LittleEndian, &tagCount); err != nil {
		return err
	}
	df.tags = make(map[int32]Tag)
	for i := int32(0); i < tagCount; i++ {
		var tag Tag
		if err := tag.Read(reader); err != nil {
			return err
		}
		df.tags[tag.Id] = tag
	}
	// Blobs
	var blobCount int32
	if err := binary.Read(reader, binary.LittleEndian, &blobCount); err != nil {
		return err
	}
	df.blobs = make(map[string]*Blob)
	for i := int32(0); i < blobCount; i++ {
		blob := new(Blob)
		if err := blob.Read(reader, df.tags); err != nil {
			return err
		}
		df.blobs[blob.Fid] = blob
	}
	return nil
}
