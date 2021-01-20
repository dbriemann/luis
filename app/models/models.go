package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)

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

// Key   string
// Value string
// }

type User struct {
	ID        int64 `db:"id"`
	CreatedAt int64 `db:"created_at"`
	UpdatedAt int64 `db:"updated_at"`

	Email       string       `db:"email"`
	Secret      string       `db:"secret"`
	Name        string       `db:"name"`
	IsAdmin     bool         `db:"is_admin"`
	Permissions []Permission `db:"-"`
	Files       []File       `db:"-"`
	Collections []Collection `db:"-"`
	// TODO: Avatar?
}

func UserByEmail(db *sqlx.DB, email string) (*User, error) {
	u := &User{}
	err := db.Get(u, "SELECT * FROM users WHERE email = LOWER( ? )", email)
	return u, err
}

func (u *User) Update(db *sqlx.DB) error {
	_, err := db.Exec(`
		UPDATE users 
		SET updated_at = ?, email = LOWER( ? ), secret = ?, name = ?, is_admin = ?
		WHERE id = ?`,
		time.Now().Unix(), u.Email, u.Secret, u.Name, u.IsAdmin, u.ID)
	return err
}

func (u *User) FetchFiles(db *sqlx.DB) error {
	err := db.Select(&u.Files, "SELECT * FROM files WHERE owner_id = ?", u.ID)
	return err
}

func (u *User) Insert(db *sqlx.DB) error {
	res, err := db.Exec(`
		INSERT INTO users (
			created_at, updated_at, email, secret, name, is_admin
		)
		VALUES (?, ?, LOWER( ? ), ?, ?, ?)`,
		u.CreatedAt,
		u.UpdatedAt,
		u.Email,
		u.Secret,
		u.Name,
		u.IsAdmin,
	)
	if err == nil {
		u.ID, err = res.LastInsertId()
	}

	return err
}

func FileByID(db *sqlx.DB, id int64) (*File, error) {
	f := &File{}
	err := db.Get(f, "SELECT * FROM files WHERE id = ?", id)
	return f, err
}

type File struct {
	ID        int64 `db:"id"`
	CreatedAt int64 `db:"created_at"`
	UpdatedAt int64 `db:"updated_at"`

	Name         string       `db:"name"`
	Thumb        string       `db:"thumb"`
	Type         FileType     `db:"type"`
	Title        string       `db:"title"`
	Description  string       `db:"description"`
	Comments     []Comment    `db:"-"`
	Stars        []Star       `db:"-"`
	Collections  []Collection `db:"-"`
	OwnerID      int64        `db:"owner_id"`
	CollectionID int64        `db:"collection_id"`

	// TODO acess/rights/special tags
}

func (f *File) Insert(db *sqlx.DB) error {
	res, err := db.Exec(`
		INSERT INTO files (
			created_at, updated_at, name, thumb, type, title, description, owner_id
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		f.CreatedAt,
		f.UpdatedAt,
		f.Name,
		f.Thumb,
		f.Type,
		f.Title,
		f.Description,
		f.OwnerID,
	)
	if err == nil {
		f.ID, err = res.LastInsertId()
	}

	return err
}

type Comment struct {
	ID        int64 `db:"id"`
	CreatedAt int64 `db:"created_at"`
	UpdatedAt int64 `db:"updated_at"`

	Text   string `db:"text"`
	FileID int64  `db:"file_id"`
}

type Star struct {
	ID int64 `db:"id"`

	FileID int64
}

type Permission struct {
	ID        int64 `db:"id"`
	CreatedAt int64 `db:"created_at"`
	UpdatedAt int64 `db:"updated_at"`

	Type         PermissionType
	CollectionID int64
	UserID       int64
}

type Collection struct {
	ID        int64 `db:"id"`
	CreatedAt int64 `db:"created_at"`
	UpdatedAt int64 `db:"updated_at"`

	Name        string
	Description string
	Cover       File
	Files       []File
	Permissions []Permission
	OwnerID     int64
}
