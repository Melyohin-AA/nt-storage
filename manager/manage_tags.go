package manager

import (
	"encoding/json"
	"nt-storage/storage"
)

func (m *Manager) GetTagData() (string, int) {
	if output, code := m.requireDataFile(); code != 0 {
		return output, code
	}
	data := struct {
		TagCats map[int32]storage.TagCategory `json:"tagCats"`
		Tags    map[int32]storage.Tag         `json:"tags"`
	}{
		TagCats: m.wad.DataFile.TagCats(),
		Tags:    m.wad.DataFile.Tags(),
	}
	encoded, err := json.Marshal(data)
	if err != nil {
		return err.Error(), 500
	}
	return string(encoded), 200
}

func (m *Manager) AddTagCat(tagCat storage.TagCategory) (string, int) {
	if output, code := m.requireDataFile(); code != 0 {
		return output, code
	}
	var err error
	if tagCat, err = m.wad.DataFile.AddTagCat(tagCat); err != nil {
		return err.Error(), 500
	}
	encoded, err := json.Marshal(tagCat)
	if err != nil {
		return err.Error(), 500
	}
	return string(encoded), 201
}

func (m *Manager) AddTag(tag storage.Tag) (string, int) {
	if output, code := m.requireDataFile(); code != 0 {
		return output, code
	}
	var err error
	if tag, err = m.wad.DataFile.AddTag(tag); err != nil {
		return err.Error(), 500
	}
	encoded, err := json.Marshal(tag)
	if err != nil {
		return err.Error(), 500
	}
	return string(encoded), 201
}
