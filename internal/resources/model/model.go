package model

import (
	"time"

	"github.com/kick-project/kick/internal/utils/errutils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

//
// Models
//

// Global global template defintion
type Global struct {
	gorm.Model
	ID     uint `gorm:"primaryKey;not null"`
	Name   string
	URL    string `gorm:"index:,unique"`
	Desc   string
	Master []Master `gorm:"many2many:global_master"`
}

// Master a set of templates
type Master struct {
	gorm.Model
	ID       uint `gorm:"primaryKey;not null"`
	Name     string
	URL      string `gorm:"index:,unique"`
	Desc     string
	Template []Template `gorm:"many2many:master_template"`
}

// Template a template definition
type Template struct {
	gorm.Model
	ID       uint `gorm:"primaryKey;not null"`
	Name     string
	URL      string `gorm:"index:,unique"`
	Desc     string
	Master   []Master `gorm:"many2many:master_template"`
	Versions []Versions
}

// Installed a table of installed templates
type Installed struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey;not null"`
	Handle   string `gorm:"index:,unique;index:idx_installed_handle_origin_url,unique"`
	Template string
	Origin   string `gorm:"index:idx_installed_handle_origin_url,unique"`
	URL      string `gorm:"index:idx_installed_handle_origin_url,unique"`
	VcsRef   string
	Desc     string
	Time     time.Time
}

// Sync last sync times
type Sync struct {
	gorm.Model
	ID         uint      `gorm:"primaryKey;not null"`
	Key        string    `gorm:"index:,unique"`
	LastUpdate time.Time `gorm:"index;column:lastupdate"`
}

// Versions template versions
type Versions struct {
	gorm.Model
	ID         uint `gorm:"primaryKey;not null"`
	Version    string
	TemplateID uint
	Template   Template `gorm:"foreignKey:TemplateID"`
}

//
//
//

// Options options create model
type Options struct {
	File string
}

// CreateModel new way of creating a schema
func CreateModel(opts *Options) (db *gorm.DB) {
	dia := sqlite.Open(opts.File)
	db, err := gorm.Open(dia, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	errutils.Efatalf("Can not initialize an ORM database: %v", err)

	err = db.AutoMigrate(
		&Global{},
		&Master{},
		&Versions{},
		&Template{},
		&Installed{},
		&Sync{},
	)
	errutils.Efatalf("can not migrate database: %v", err)
	result := db.Clauses(clause.Insert{Modifier: "OR IGNORE"}).Create(&Master{
		Name: "local",
		URL:  "none",
		Desc: "This template is generated locally",
	})

	if result.Error != nil {
		errutils.Efatalf("can not insert root record into database: %v", result.Error)
	}
	return db
}
