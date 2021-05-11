package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/kick-project/kick/internal/resources/errs"
)

// Move moves a file from src to dest
func Move(src, dst string) error {
	err := os.Rename(src, dst)
	if err == nil {
		return err
	}

	m := &mover{}
	// Get information on files
	sInfo, _ := Stat(src)
	dInfo, _ := Stat(dst)
	err = m.move(sInfo, dInfo)
	if err != nil {
		err := m.rollback()
		errs.LogF(`error rollingback %v`, err)
		return err
	}

	err = m.cleanup()
	return err
}

type mover struct {
	trashList    []Info
	rollBackList []Info
}

// flagRemoval flag a file for cleanup
func (m *mover) flagRemoval(path Info) {
	m.trashList = append(m.trashList, path)
}

// flagRollback flag a file for rollback if needed
func (m *mover) flagRollback(path Info) {
	m.rollBackList = append(m.rollBackList, path)
}

// move
func (m *mover) move(sInfo, dInfo Info) error {
	// Source doesn't exist
	if !sInfo.Exists() {
		return fmt.Errorf(`source file does not exists`)
	}

	// Mis-matching types
	if dInfo.Exists() && sInfo.IsDir() != dInfo.IsDir() {
		return fmt.Errorf(`source and destination exists but are have mis-matching file types`)
	}

	if sInfo.IsDir() {
		err := m.dirMove(sInfo, dInfo)
		if err != nil {
			return err
		}
	}
	err := m.fileMove(sInfo, dInfo)
	if err != nil {
		return err
	}
	return err
}

// fileMove move files
func (m *mover) fileMove(src, dst Info) error {
	// Move files using copy
	_, err := Copy(src.Path(), dst.Path())
	if err != nil {
		return err
	}
	m.flagRemoval(src)
	m.flagRollback(dst)
	return nil
}

// dirMove move directories
func (m *mover) dirMove(src, dst Info) error {
	// TODO: match mode and permissions
	err := os.Mkdir(dst.Path(), 0755)
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(".")
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.Name() == "." || f.Name() == ".." {
			continue
		}
		newSrc := filepath.Join(src.Path(), f.Name())
		newDest := filepath.Join(dst.Path(), f.Name())

		sInfo, _ := Stat(newSrc)
		dInfo, _ := Stat(newDest)
		err := m.move(sInfo, dInfo)
		if err != nil {
			return err
		}
	}
	m.flagRemoval(src)
	m.flagRollback(dst)

	return nil
}

// rollback copied files
func (m *mover) rollback() error {
	m.unlink(m.trashList)
	return nil
}

// cleanup clean up source files
func (m *mover) cleanup() error {
	m.unlink(m.trashList)
	return nil
}

// unlink
func (m *mover) unlink(unlinks []Info) {
	sort.Slice(unlinks, func(i, j int) bool {
		return len(unlinks[j].Abs()) < len(unlinks[i].Abs())
	})
	for _, u := range unlinks {
		if u.Exists() {
			err := os.Remove(u.Abs())
			errs.LogF(`can not remove file %s: %v`, u.Path(), err)
		}
	}
}
