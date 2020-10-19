package checksum

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kick-project/kick/internal/utils/errutils"
)

func Sha256SumFile(srcFile string, sumfile string) (sum string, err error) {
	srcIO, err := os.Open(srcFile)
	if err != nil {
		return "", fmt.Errorf("Failed to open file %s: %w", srcFile, err)
	}
	bytesum, err := Sha256Sum(srcIO)
	if err != nil {
		return "", fmt.Errorf("Failed to generate checksum: %w", err)
	}
	fname := filepath.Base(srcFile)
	sum = fmt.Sprintf("%x %s\n", bytesum, fname)
	dstIO, err := ioutil.TempFile("", "")
	if err != nil {
		return "", fmt.Errorf("Failed to open temporary file for %s: %s", sumfile, err)
	}
	_, err = dstIO.WriteString(sum)
	errutils.Epanic(err)
	dstIO.Close()
	err = os.Rename(dstIO.Name(), sumfile)
	errutils.Epanic(err)
	return sum, err
}

func Sha256Sum(rdr io.Reader) (bytesum []byte, err error) {
	hash := sha256.New()
	_, err = io.Copy(hash, rdr)
	errutils.Efatalf("Error copy bytes: %w", err)
	bytesum = hash.Sum(nil)
	return bytesum, nil
}

func VerifySha256sum(srcFile, sumfile string) (pass bool, sum string, err error) {
	srcAbs, err := filepath.Abs(srcFile)
	errutils.Efatalf("Can not get absolute path for %s: %v", srcFile, err)

	sumfileDir, err := filepath.Abs(filepath.Dir(sumfile))
	errutils.Efatalf("Can not get absolute path for %s: %v", sumfile, err)

	// Get original checksum
	file, err := os.Open(sumfile)
	errutils.Efatalf("Error opening file %s: %w", sumfile, err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	err = scanner.Err()
	errutils.Efatalf("Error scanning file: %w", err)

	var origsum string
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 2 {
			continue
		}
		origPath := filepath.Join(sumfileDir, fields[1])
		if srcAbs == origPath {
			origsum = fields[0]
		}
	}

	// Get checksum from source
	srcIO, err := os.Open(srcAbs)
	if err != nil {
		return false, "", fmt.Errorf("Failed to open file %s: %w", srcFile, err)
	}
	bytesum, err := Sha256Sum(srcIO)
	errutils.Efatalf("Can not get sha256sum for %s: %v", srcAbs, err)
	sum = fmt.Sprintf("%x", bytesum)

	if origsum != "" && sum == origsum {
		return true, sum, nil
	}
	return false, sum, nil
}
