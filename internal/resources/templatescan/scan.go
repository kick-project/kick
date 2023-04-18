package templatescan

import (
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/kick-project/kick/internal/resources/config/configtemplate"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/marshal"
	"github.com/kick-project/kick/internal/resources/model"
	"github.com/kick-project/kick/internal/resources/modeline"
	"gorm.io/gorm"
)

var ErrInvalidFileType = errors.New("Invalid file type")

//go:embed scan_label.sql
var QueryScanLabel string

//go:embed scan_option.sql
var QueryScanOption string

type Scan struct {
	DB *gorm.DB
}

// Run scans root directory for files containing mode lines.
func (s Scan) Run(root string, lines int) (err error) {
	defer func() {
		// Work around. GORM places an entry with a 0 ID entry in bridging tables.
		s.DB.Exec("DELETE FROM file_option WHERE option_id=0")
		s.DB.Exec("DELETE FROM file_label WHERE label_id=0")
	}()

	err = s.scanBase(root, lines)
	if err != nil {
		return
	}

	err = s.readConf(root)
	return
}

// reads the template configuration
func (s Scan) readConf(root string) (err error) {
	p := filepath.Join(root, ".kick.yml")
	info, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return fmt.Errorf("unable to stat file %s: %w", p, err)
		}
	}
	if info.IsDir() {
		return ErrInvalidFileType
	}
	mlt := &model.Base{}
	tx := s.DB.Where("base = ?", root).First(mlt)
	if tx.RowsAffected == 0 {
		return nil
	}

	conf := &configtemplate.TemplateMain{}
	err = marshal.FromFile(conf, p)
	if err != nil {
		return err
	}

	for file, labels := range conf.Labels {
		path := filepath.FromSlash(file)
		mlfs := []*model.File{}
		tx = s.DB.Where("base_id = ? AND (file = ? OR file LIKE ?)", mlt.ID, path, path+string(filepath.Separator)+"%").Find(&mlfs)
		if tx.RowsAffected == 0 {
			continue
		}
		for _, entry := range mlfs {
			for _, l := range labels {
				//nolint
				err = s.many2ManyFileLabel(entry, l)
			}
		}
	}
	return nil
}

func (s Scan) scanBase(root string, lines int) error {
	mlt, err := s.fetchTemplate(root)
	if err != nil {
		return err
	}

	fileSystem := os.DirFS(root)
	var real string
	return fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if path == "." {
			return nil
		}
		if path == ".git" && d.IsDir() {
			return fs.SkipDir
		}
		mlf := s.fetchFile(mlt.ID, path)

		if d.Type().IsRegular() {
			real = filepath.Join(root, path)
			ml, err := modeline.Parse(real, nil, lines)
			if err != nil {
				return nil
			}
			s.updateModeInfo(ml, mlf)
		}
		return nil
	})
}

func (s Scan) updateModeInfo(ml *modeline.ModeLine, mlf *model.File) {
	if ml != nil {
		for _, o := range ml.GetOptions() {
			Option := &model.Option{Option: o}
			err := s.DB.Model(mlf).Association("Option").Append(Option)
			if err != nil {
				errs.Fatal(err)
			}
		}
		for _, l := range ml.GetLabel() {
			Label := &model.Label{Label: l}
			err := s.DB.Model(mlf).Association("Label").Append(Label)
			if err != nil {
				errs.Fatal(err)
			}
		}
	}
}

func (s Scan) fetchFile(id uint, path string) *model.File {
	mlf := &model.File{}
	tx := s.DB.Where("base_id = ? AND file = ?", id, path).First(mlf)
	if tx.RowsAffected == 0 {
		mlf.File = path
		mlf.BaseID = id
		s.DB.Create(mlf)
	}
	return mlf
}

func (s Scan) fetchTemplate(root string) (*model.Base, error) {
	mlt := &model.Base{}
	tx := s.DB.Where("base = ?", root).Find(mlt)
	if tx.RowsAffected == 0 {
		mlt.Base = root
		tx = s.DB.Create(mlt)
		if tx.Error != nil {
			return nil, tx.Error
		}
	}
	return mlt, nil
}

// workaround to perform many2many association.
func (s Scan) many2ManyFileLabel(file *model.File, label string) error {
	l := model.Label{Label: label}
	tx := s.DB.Where("label = ?", label).First(&l)
	if tx.RowsAffected == 0 {
		tx2 := s.DB.Create(&l)
		if tx2.Error != nil {
			return tx2.Error
		}
	}
	tx = s.DB.Exec("INSERT OR IGNORE INTO file_label (file_id, label_id) VALUES (?, ?)", file.ID, l.ID)
	return tx.Error
}
