package file

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"github.com/bsthun/gut"
	"github.com/strawst/strawhouse-go"
	"github.com/strawst/strawhouse-go/pb"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (r *Service) Upload(name string, directory string, attribute []byte, content io.Reader) (*string, []byte, *string, *gut.ErrorInstance) {
	// * Validate name
	if strings.HasPrefix(name, ".") {
		return nil, nil, nil, gut.Err(false, "invalid filename", nil)
	}
	if len(name) < 3 {
		return nil, nil, nil, gut.Err(false, "invalid filename", nil)
	}

	// * Normalize file name
	name = r.filepath.BaseName(name)

	// * Construct path
	relativeFilePath := filepath.Clean(filepath.Join(directory, name))
	absoluteFilePath := r.filepath.AbsPath(filepath.Clean(relativeFilePath))
	absoluteDirectoryPath := r.filepath.AbsPath(directory)

	// * Ensure directory
	if err := os.MkdirAll(absoluteDirectoryPath, 0700); err != nil {
		return nil, nil, nil, gut.Err(false, "unable to create directory", err)
	}

	// * Create file
	file, err := os.Create(absoluteFilePath)
	if err != nil {
		return nil, nil, nil, gut.Err(false, "unable to create file", err)
	}

	// * Calculate sha256 hash
	hash := sha256.New()
	fileBuffer := make([]byte, 1024)
	for {
		n, err := content.Read(fileBuffer)
		if err != nil {
			break
		}
		if _, err := file.Write(fileBuffer[:n]); err != nil {
			return nil, nil, nil, gut.Err(false, "unable to write file", err)
		}
		hash.Write(fileBuffer[:n])
	}
	sum := hash.Sum(nil)
	if err := file.Close(); err != nil {
		return nil, nil, nil, gut.Err(false, "unable to close file", err)
	}

	// * Check hash
	if val, err := r.pogreb.Sum.Get(sum); err != nil {
		return nil, nil, nil, gut.Err(false, "unable to check hash", err)
	} else {
		if val != nil {
			// Check if file already exists in other path
			pathVal := string(val)
			if pathVal != relativeFilePath {
				return nil, nil, nil, gut.Err(false, "file already exists in other path", nil)
			}

			// If file exists and not corrupted, deny upload
			if _, err := os.Stat(absoluteFilePath); err == nil {
				if er := r.fileflag.Corrupted(relativeFilePath); er == nil {
					return nil, nil, nil, gut.Err(false, "file already exist", nil)
				}
			}
		}
	}

	// * Save hash
	if err := r.pogreb.Sum.Put(sum, []byte(relativeFilePath)); err != nil {
		return nil, nil, nil, gut.Err(false, "unable to save hash", err)
	}

	// * Save log
	size := r.pogreb.Log.Count()
	sizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBytes, size)
	if err := r.pogreb.Log.Put(sizeBytes, sum); err != nil {
		return nil, nil, nil, gut.Err(false, "unable to save log", err)
	}

	// * Encode base64 hash
	encoded := base64.StdEncoding.EncodeToString(sum)
	encoded = encoded[:len(encoded)-1]
	r.signature.ReplaceClean(&encoded)

	// * Set file flag
	if er := r.fileflag.SumSet(relativeFilePath, sum); er != nil {
		return nil, nil, nil, er
	}

	// * Construct file flag
	if er := r.fileflag.CorruptedInit(relativeFilePath); er != nil {
		return nil, nil, nil, er
	}

	// * Fire event feed
	r.eventfeed.Fire(strawhouse.FeedTypeUpload, relativeFilePath, &pb.UploadFeedResponse{
		Name:      name,
		Directory: directory,
		Hash:      encoded,
		Attr:      attribute,
	})

	return &relativeFilePath, sum, &encoded, nil
}
