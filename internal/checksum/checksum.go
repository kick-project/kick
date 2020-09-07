package checksum

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/crosseyed/prjstart/internal/utils"
)

func Sha256sum(srcFile string, sumfile string) (sum string, err error) {
	bytesum, err := sha256sum(srcFile)

	fname := filepath.Base(srcFile)
	sum = fmt.Sprintf("%x %s\n", bytesum, fname)
	dstIO, err := ioutil.TempFile("", "")
	if err != nil {
		return "", fmt.Errorf("Failed to open temporary file for %s: %s", sumfile, err)
	}
	dstIO.WriteString(sum)
	dstIO.Close()
	os.Rename(dstIO.Name(), sumfile)
	return sum, err
}

func sha256sum(srcFile string) (bytesum []byte, err error) {
	srcIO, err := os.Open(srcFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file %s: %s", srcFile, err)
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, srcIO); err != nil {
		log.Fatal(err)
	}
	bytesum = hash.Sum(nil)
	return bytesum, nil
}

func VerifySha256sum(srcFile, sumfile string) (pass bool, sum string, err error) {
	srcAbs, err := filepath.Abs(srcFile)
	utils.ChkErr(err, utils.Efatalf, "Can not get absoulte path for %s: %v", srcFile, err)

	sumfileDir, err := filepath.Abs(filepath.Dir(sumfile))
	utils.ChkErr(err, utils.Efatalf, "Can not get absoulte path for %s: %v", sumfile, err)

	// Get original checksum
	file, err := os.Open(sumfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

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
	bytesum, err := sha256sum(srcAbs)
	utils.ChkErr(err, utils.Efatalf, "Can not get sha256sum for %s: %v", srcAbs, err)
	sum = fmt.Sprintf("%x", bytesum)

	if origsum != "" && sum == origsum {
		return true, sum, nil
	}
	return false, sum, nil
}
