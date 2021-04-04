package model

import (
	"time"

	"gorm.io/gorm"
)

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
	Versions []Versions `gorm:"many2many:template_version"`
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
	LastUpdate time.Time `gorm:"index"`
}

// Versions template versions
type Versions struct {
	gorm.Model
	ID      uint `gorm:"primaryKey;not null"`
	Version string
}
