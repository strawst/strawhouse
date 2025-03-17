package strawhouse

import (
	"github.com/bsthun/gut"
	"io"
)

type PluginCallback interface {
	Get(path string, writer io.Writer) *gut.ErrorInstance
	Upload(name string, directory string, attribute []byte, file io.Reader) (*string, []byte, *string, *gut.ErrorInstance)
	Bind(typ FeedType, dir string, handler func(resp any)) uint64
	Unbind(typ FeedType, dir string, id uint64)
}

type Plugin interface {
	Load(callback PluginCallback)
	Unload()
}
