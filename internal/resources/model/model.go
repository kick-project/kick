package model

import (
	"time"

	"github.com/kick-project/kick/internal/resources/errs"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

//
// Models
//

// Repo a set of templates
type Repo struct {
	gorm.Model
	ID       uint `gorm:"primaryKey;not null"`
	Name     string
	URL      string `gorm:"index:,unique"`
	Desc     string
	Template []Template `gorm:"many2many:repo_template"`
}

// Template a template definition
type Template struct {
	gorm.Model
	ID       uint `gorm:"primaryKey;not null"`
	Name     string
	URL      string `gorm:"index:,unique"`
	Desc     string
	Repo     []Repo `gorm:"many2many:repo_template"`
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
	ID         uint     `gorm:"primaryKey;not null"`
	Version    string   `gorm:"index:idx_version_template,unique"`
	TemplateID uint     `gorm:"index:idx_version_template,unique"`
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

	errs.FatalF("Can not initialize an ORM database: %v", err)

	err = db.AutoMigrate(
		&Repo{},
		&Versions{},
		&Template{},
		&Installed{},
		&Sync{},
	)
	errs.FatalF("can not migrate database: %v", err)

	// Insert base repo
	m := &Repo{
		Name: "local",
		URL:  "none",
		Desc: "Locally defined templates",
	}
	result := db.Clauses(clause.Insert{Modifier: "OR IGNORE"}).Create(m)
	if result.Error != nil {
		errs.FatalF("can not insert root record into database: %v", result.Error)
	}

	return db
}
