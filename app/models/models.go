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
	err := db.Select(&u.Files, "SELECT * FROM files WHERE owner_id = ? ORDER BY date DESC", u.ID)
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

//[
// MetaData example data.. TODO :)
// {
//   "image": {
//     "name": "/app/tmp/9eFSVHsaPB",
//     "format": "JPEG",
//     "formatDescription": "Joint Photographic Experts Group JFIF format",
//     "mimeType": "image/jpeg",
//     "class": "DirectClass",
//     "geometry": {
//       "width": 640,
//       "height": 480,
//       "x": 0,
//       "y": 0
//     },
//     "resolution": {
//       "x": "300",
//       "y": "300"
//     },
//     "printSize": {
//       "x": "2.1333333333333333037",
//       "y": "1.6000000000000000888"
//     },
//     "units": "PixelsPerInch",
//     "type": "TrueColor",
//     "endianess": "Undefined",
//     "colorspace": "sRGB",
//     "depth": 8,
//     "baseDepth": 8,
//     "channelDepth": {
//       "red": 8,
//       "green": 8,
//       "blue": 8
//     },
//     "pixels": 307200,
//     "imageStatistics": {
//       "all": {
//         "min": "0",
//         "max": "255",
//         "mean": "120.794",
//         "standardDeviation": "48.999",
//         "kurtosis": "2.04928",
//         "skewness": "0.518196"
//       }
//     },
//     "channelStatistics": {
//       "red": {
//         "min": "10",
//         "max": "255",
//         "mean": "152.362",
//         "standardDeviation": "49.3103",
//         "kurtosis": "-0.459491",
//         "skewness": "0.267462"
//       },
//       "green": {
//         "min": "0",
//         "max": "255",
//         "mean": "132.752",
//         "standardDeviation": "49.7499",
//         "kurtosis": "-0.338396",
//         "skewness": "0.29157"
//       },
//       "blue": {
//         "min": "0",
//         "max": "255",
//         "mean": "77.2686",
//         "standardDeviation": "47.918",
//         "kurtosis": "2.32993",
//         "skewness": "1.31755"
//       }
//     },
//     "renderingIntent": "Perceptual",
//     "gamma": 0.454545,
//     "chromaticity": {
//       "redPrimary": {
//         "x": 0.64,
//         "y": 0.33
//       },
//       "greenPrimary": {
//         "x": 0.3,
//         "y": 0.6
//       },
//       "bluePrimary": {
//         "x": 0.15,
//         "y": 0.06
//       },
//       "whitePrimary": {
//         "x": 0.3127,
//         "y": 0.329
//       }
//     },
//     "backgroundColor": "#FFFFFF",
//     "borderColor": "#DFDFDF",
//     "matteColor": "#BDBDBD",
//     "transparentColor": "#000000",
//     "interlace": "None",
//     "intensity": "Undefined",
//     "compose": "Over",
//     "pageGeometry": {
//       "width": 640,
//       "height": 480,
//       "x": 0,
//       "y": 0
//     },
//     "dispose": "Undefined",
//     "iterations": 0,
//     "compression": "JPEG",
//     "quality": 75,
//     "orientation": "TopLeft",
//     "properties": {
//       "date:create": "2019-01-12T02:51:28+00:00",
//       "date:modify": "2019-01-12T02:51:28+00:00",
//       "exif:ColorSpace": "1",
//       "exif:ComponentsConfiguration": "1, 2, 3, 0",
//       "exif:Contrast": "0",
//       "exif:CustomRendered": "0",
//       "exif:DateTime": "2008:11:01 21:15:07",
//       "exif:DateTimeDigitized": "2008:10:22 16:28:39",
//       "exif:DateTimeOriginal": "2008:10:22 16:28:39",
//       "exif:DigitalZoomRatio": "0/100",
//       "exif:ExifImageLength": "480",
//       "exif:ExifImageWidth": "640",
//       "exif:ExifOffset": "268",
//       "exif:ExifVersion": "48, 50, 50, 48",
//       "exif:ExposureBiasValue": "0/10",
//       "exif:ExposureMode": "0",
//       "exif:ExposureProgram": "2",
//       "exif:ExposureTime": "4/300",
//       "exif:FileSource": "3",
//       "exif:Flash": "16",
//       "exif:FlashPixVersion": "48, 49, 48, 48",
//       "exif:FNumber": "59/10",
//       "exif:FocalLength": "24/1",
//       "exif:FocalLengthIn35mmFilm": "112",
//       "exif:GainControl": "0",
//       "exif:GPSAltitudeRef": "0",
//       "exif:GPSDateStamp": "2008:10:23",
//       "exif:GPSImgDirectionRef": null,
//       "exif:GPSInfo": "926",
//       "exif:GPSLatitude": "43/1, 28/1, 281400000/100000000",
//       "exif:GPSLatitudeRef": "N",
//       "exif:GPSLongitude": "11/1, 53/1, 645599999/100000000",
//       "exif:GPSLongitudeRef": "E",
//       "exif:GPSMapDatum": "WGS-84   ",
//       "exif:GPSSatellites": "06",
//       "exif:GPSTimeStamp": "14/1, 27/1, 724/100",
//       "exif:ImageDescription": "                               ",
//       "exif:InteroperabilityOffset": "896",
//       "exif:ISOSpeedRatings": "64",
//       "exif:LightSource": "0",
//       "exif:Make": "NIKON",
//       "exif:MakerNote": "78, 105, 107, 111, 110, 0, ..., 0, 0, 0, ",
//       "exif:MaxApertureValue": "29/10",
//       "exif:MeteringMode": "5",
//       "exif:Model": "COOLPIX P6000",
//       "exif:Orientation": "1",
//       "exif:ResolutionUnit": "2",
//       "exif:Saturation": "0",
//       "exif:SceneCaptureType": "0",
//       "exif:SceneType": "1",
//       "exif:Sharpness": "0",
//       "exif:Software": "Nikon Transfer 1.1 W",
//       "exif:SubjectDistanceRange": "0",
//       "exif:thumbnail:Compression": "6",
//       "exif:thumbnail:InteroperabilityIndex": "R98",
//       "exif:thumbnail:InteroperabilityVersion": "48, 49, 48, 48",
//       "exif:thumbnail:JPEGInterchangeFormat": "4548",
//       "exif:thumbnail:JPEGInterchangeFormatLength": "6702",
//       "exif:thumbnail:ResolutionUnit": "2",
//       "exif:thumbnail:XResolution": "72/1",
//       "exif:thumbnail:YResolution": "72/1",
//       "exif:UserComment": "65, 83, 67, 73, 73, 0, ...., 32, 32, 32, 0",
//       "exif:WhiteBalance": "0",
//       "exif:XResolution": "300/1",
//       "exif:YCbCrPositioning": "1",
//       "exif:YResolution": "300/1",
//       "jpeg:colorspace": "2",
//       "jpeg:sampling-factor": "2x1,1x1,1x1",
//       "MicrosoftPhoto:Rating": "0",
//       "signature": "55cbf121d52110cda7c785d97bf02f6a31bd0f5ac44c06f9b2f70c9c7d00ade4"
//     },
//     "profiles": {
//       "exif": {
//         "length": "11256"
//       },
//       "xmp": {
//         "length": "4000"
//       }
//     },
//     "artifacts": {
//       "filename": "/app/tmp/9eFSVHsaPB"
//     },
//     "tainted": false,
//     "filesize": "162KB",
//     "numberPixels": "307K",
//     "pixelsPerSecond": "15.36MB",
//     "userTime": "0.010u",
//     "elapsedTime": "0:01.020",
//     "version": "ImageMagick 6.9.5-9 Q16 x86_64 2016-10-21 http://www.imagemagick.org"
//   }
// }
//]
type MetaData struct {
	Image MetaImage `json:"image"`
}

type MetaImage struct {
	Properties MetaProperties `json:"properties"`
}

type MetaProperties struct {
	DateTimeOriginal string `json:"exif:DateTimeOriginal"`
}

type EXIF struct {
	Date int64
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
	Date         int64        `db:"date"`
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
			created_at, updated_at, name, date, thumb, type, title, description, owner_id
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		f.CreatedAt,
		f.UpdatedAt,
		f.Name,
		f.Date,
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
