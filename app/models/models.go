package models

import "gorm.io/gorm"

type FileType int

const (
	FileTypeUnknown FileType = 0
	FileTypeImage   FileType = 1
	FileTypeVideo   FileType = 2
	FileTypePDF     FileType = 3
	FileTypeTXT     FileType = 4
)

type Tag struct {
	gorm.Model

	Key   string
	Value string
}

type Comment struct {
	gorm.Model

	Text   string
	FileID uint // Foreign key (File).
}

type Star struct {
	gorm.Model

	FileID uint // Foreign key (File).
}

type Permission struct {
	// TODO: Role?
}

type Collection struct {
	gorm.Model

	Name        string
	Description string
	Cover       File
	Files       []File `gorm:"many2many:collection_files;"`
	OwnerID     uint   // Foreign key (Owner/User).
}

type File struct {
	gorm.Model

	Path         string
	Name         string
	Ext          string
	Type         FileType
	Title        string
	Description  string
	Comments     []Comment    // Has-many relationship.
	Stars        []Star       // Has-many relationship.
	Tags         []Tag        `gorm:"many2many:file_tags;"`
	Collections  []Collection `gorm:"many2many:collection_files;"`
	OwnerID      uint         // Foreign key (Owner/User).
	CollectionID uint         // Foreign key (Collection/Cover)

	// TODO acess/rights/special tags
}
type User struct {
	gorm.Model

	Email       string `gorm:"uniqueIndex"`
	Secret      string
	Name        string
	IsAdmin     bool
	Files       []File       `gorm:"foreignKey:OwnerID"`
	Collections []Collection `gorm:"foreignKey:OwnerID"`
	// TODO: Avatar?
}
