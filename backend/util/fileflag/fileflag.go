package fileflag

import (
	"github.com/strawst/strawhouse-go"
	"strawhouse-backend/common/config"
	"strawhouse-backend/util/filepath"
)

const xattrSumTag = "user.sh.sum"
const xattrFlagTag = "user.sh.flag"

type Fileflag struct {
	config    *config.Config
	filepath  *filepath.Filepath
	signature *strawhouse.Signature
}

func Init(config *config.Config, filepath *filepath.Filepath, signature *strawhouse.Signature) *Fileflag {
	return &Fileflag{
		config:    config,
		filepath:  filepath,
		signature: signature,
	}
}
