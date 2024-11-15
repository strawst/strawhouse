package fileflag

import (
	"github.com/strawst/strawhouse-go"
	"strawhouse-backend/util/filepath"
)

const xattrSumTag = "user.sh.sum"
const xattrFlagTag = "user.sh.flag"

type Fileflag struct {
	filepath  *filepath.Filepath
	signature *strawhouse.Signature
}

func Init(filepath *filepath.Filepath, signature *strawhouse.Signature) *Fileflag {
	return &Fileflag{
		filepath:  filepath,
		signature: signature,
	}
}
