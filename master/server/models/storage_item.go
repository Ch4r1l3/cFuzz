package models

import (
	"os"
)

// item that record some file info
// swagger:model
type StorageItem struct {
	// example: 1
	ID uint64 `gorm:"primary_key" json:"id"`

	// example: afl
	Name string `json:"name" sql:"type:varchar(255) NOT NULL UNIQUE"`

	Path string `json:"path"`

	// example: fuzzer
	Type string `json:"type"`

	// example: true
	ExistsInImage bool `json:"existsInImage"`

	// if upload file is zip and type is not corpus, this field specefiy the path of file like target
	// example: test/target
	RelPath string `json:"relPath"`

	// example: 1
	UserID uint64 `json:"userID" sql:"type:integer REFERENCES user(id) ON DELETE CASCADE"`
}

func (s *StorageItem) Delete() error {
	if !s.ExistsInImage {
		os.RemoveAll(s.Path)
	}
	return nil
}

// storage item types
const (
	Fuzzer = "fuzzer"
	Target = "target"
	Corpus = "corpus"
)

var StorageItemTypes = []string{Fuzzer, Target, Corpus}

func IsStorageItemTypeValid(mtype string) bool {
	switch mtype {
	case
		Fuzzer,
		Target,
		Corpus:
		return true
	}
	return false
}
