package fileflag

import (
	"github.com/bsthun/gut"
	"github.com/pkg/xattr"
)

func (r *Fileflag) Corrupted(relativePath string) *gut.ErrorInstance {
	// * Xattr disabled
	if *r.config.EnableXattr == false {
		return nil
	}

	// * Convert path
	absolutePath := r.filepath.AbsPath(relativePath)
	flag, err := xattr.Get(absolutePath, xattrFlagTag)
	if err != nil {
		return gut.Err(false, "unable to get file flag attributes", err)
	}
	if flag[0]&0b00001000 != 0 {
		return gut.Err(false, "file corrupted")
	}
	return nil
}

func (r *Fileflag) CorruptedInit(relativePath string) *gut.ErrorInstance {
	// * Xattr disabled
	if *r.config.EnableXattr == false {
		return nil
	}

	// * Convert path
	absolutePath := r.filepath.AbsPath(relativePath)
	flag := make([]byte, 4)
	err := xattr.Set(absolutePath, xattrFlagTag, flag)
	if err != nil {
		return gut.Err(false, "unable to init file flag attributes", err)
	}
	return nil
}

func (r *Fileflag) CorruptedSet(relativePath string, status bool) *gut.ErrorInstance {
	// * Xattr disabled
	if *r.config.EnableXattr == false {
		return nil
	}

	// * Convert path
	absolutePath := r.filepath.AbsPath(relativePath)
	flag, err := xattr.Get(absolutePath, xattrFlagTag)
	if err != nil {
		return gut.Err(false, "unable to get file flag attributes", err)
	}
	if status {
		flag[0] |= 0b00001000
	} else {
		flag[0] &= 0b11110111
	}
	err = xattr.Set(absolutePath, xattrFlagTag, flag)
	if err != nil {
		return gut.Err(false, "unable to set file flag attributes", err)
	}
	return nil
}
