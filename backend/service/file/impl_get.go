package file

import (
	"bytes"
	"crypto/sha256"
	"github.com/bsthun/gut"
	"io"
	"os"
)

func (r *Service) Get(path string, writer io.Writer) *gut.ErrorInstance {
	// * Construct absolute path
	absolutePath := r.filepath.AbsPath(path)

	// * Check file corruption
	if er := r.fileflag.Corrupted(path); er != nil {
		return er
	}

	// * Check signed checksum
	sum, er := r.fileflag.SumGet(path)
	if er != nil {
		return er
	}
	if sum != nil {
		// * Check file path
		val, err := r.pogreb.Sum.Get(sum)
		if err != nil || val == nil {
			return gut.Err(false, "file record not found")
		}
		if !bytes.Equal(val, []byte(path)) {
			return gut.Err(false, "file path mismatch")
		}
	}

	// * Open the file
	file, err := os.Open(absolutePath)
	if err != nil {
		return gut.Err(false, "unable to open file", err)
	}
	defer func() {
		_ = file.Close()
	}()

	// * Initialize hash
	hash := sha256.New()
	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			break
		}
		if n > 0 {
			hash.Write(buffer[:n])
			_, err := writer.Write(buffer[:n])
			if err != nil {
				return gut.Err(false, err.Error())
			}
		}
	}

	// * Compare content hash and xattr hash
	actual := hash.Sum(nil)
	if sum != nil {
		if !bytes.Equal(actual, sum) {
			if er := r.fileflag.CorruptedSet(path, true); er != nil {
				return er
			}
			return gut.Err(false, "invalid file hash")
		}
	} else {
		// * Get path value
		val, err := r.pogreb.Sum.Get(actual)
		if err != nil || val == nil {
			return gut.Err(false, "file record not found")
		}

		// * Check file path
		if !bytes.Equal(val, []byte(path)) {
			return gut.Err(false, "file path mismatch")
		}
	}

	return nil
}
