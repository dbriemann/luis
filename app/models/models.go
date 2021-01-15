package models

import "gorm.io/gorm"

type FileType int
type PermissionType int

const (
	FileTypeUnknown FileType = 0
	FileTypeImage   FileType = 1
	FileTypeVideo   FileType = 2
	FileTypePDF     FileType = 3
	FileTypeTXT     FileType = 4

	PermissionTypeView PermissionType = 0
	PermissionTypeAdd  PermissionType = 1
)

// type Tag struct {
// TODO: think about a simple tag solution that works together
// with the permission management of collections.

// gorm.Model

// Key   string
// Value string
// }

type Comment struct {
	gorm.Model

	Text   string
	FileID uint // Belongs to 1 File.
}

type Star struct {
	gorm.Model

	FileID uint // Belongs to 1 File.
}

type Permission struct {
	gorm.Model

	Type         PermissionType
	CollectionID uint // Belongs to 1 Collection.
	UserID       uint // Belongs to 1 User.
}

type Collection struct {
	gorm.Model

	Name        string
	Description string
	Cover       File
	Files       []File       `gorm:"many2many:collection_files;"` // Has 0..n Files.
	Permissions []Permission // Has 0..n Permissions.
	OwnerID     uint         // Belongs to 1 User.
}

type File struct {
	gorm.Model

	Name         string
	Thumb        string
	Type         FileType
	Title        string
	Description  string
	Comments     []Comment    // Has 0..n Comments.
	Stars        []Star       // Has 0..n Stars.
	Collections  []Collection `gorm:"many2many:collection_files;"` // Is in to 0..n Collections.
	OwnerID      uint         // Belongs to 1 User as owner.
	CollectionID uint         // Belongs to 1 Collection as cover.

	// TODO acess/rights/special tags
}

type User struct {
	gorm.Model

	Email       string `gorm:"uniqueIndex"`
	Secret      string
	Name        string
	IsAdmin     bool
	Permissions []Permission // Has 0..n Permissions.
	Files       []File       `gorm:"foreignKey:OwnerID"` // Has 0..n Files.
	Collections []Collection `gorm:"foreignKey:OwnerID"` // Has 0..n Collections.
	// TODO: Avatar?
}
