package fileflag

import (
	"bytes"
	"github.com/bsthun/gut"
	"github.com/pkg/xattr"
)

func (r *Fileflag) SumSet(relativePath string, sum []byte) *gut.ErrorInstance {
	// * Xattr disabled
	if *r.config.EnableXattr == false {
		return nil
	}

	// * Convert path
	absolutePath := r.filepath.AbsPath(relativePath)

	// * Sign checksum
	hash := r.signature.GetHash()
	hash.Write(sum)
	signedSum := hash.Sum(nil)
	r.signature.PutHash(hash)

	// * Set file attributes
	err := xattr.Set(absolutePath, xattrSumTag, append(sum, signedSum...))
	if err != nil {
		return gut.Err(false, "unable to set file sum attributes", err)
	}

	return nil
}

func (r *Fileflag) SumGet(relativePath string) (r1 []byte, er *gut.ErrorInstance) {
	// * Xattr disabled
	if *r.config.EnableXattr == false {
		return nil, nil
	}

	defer func() {
		if er != nil {
			_ = r.CorruptedSet(relativePath, true)
		}
	}()
	// * Convert path
	absolutePath := r.filepath.AbsPath(relativePath)

	// * Get file attributes
	attr, err := xattr.Get(absolutePath, xattrSumTag)
	if err != nil {
		return nil, gut.Err(false, "unable to get file sum attributes", err)
	}

	// * Verify sum attributes
	hash := r.signature.GetHash()
	if len(attr) != hash.Size()*2 {
		return nil, gut.Err(false, "invalid sum attributes length", nil)
	}

	// * Split sum and signed sum
	sum := attr[:hash.Size()]
	signedSum := attr[hash.Size():]

	// * Verify signed checksum
	hash.Write(sum)
	expectedSignedSum := hash.Sum(nil)
	r.signature.PutHash(hash)

	// * Compare signed checksum
	if !bytes.Equal(signedSum, expectedSignedSum) {
		return nil, gut.Err(false, "invalid signed checksum", nil)
	}

	return sum, nil
}
