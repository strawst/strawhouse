package plugin

import (
	"github.com/bsthun/gut"
	"github.com/strawst/strawhouse-go"
	"io"
	"strawhouse-backend/service/file"
	"strawhouse-backend/util/eventfeed"
)

func (r *Service) Plugger() strawhouse.PluginCallback {
	return &Plugger{
		file:      r.file,
		eventfeed: r.eventfeed,
	}
}

type Plugger struct {
	file      *file.Service
	eventfeed *eventfeed.EventFeed
}

func (r *Plugger) Get(path string, writer io.Writer) *gut.ErrorInstance {
	return r.file.Get(path, writer)
}

func (r *Plugger) Upload(name string, directory string, attribute []byte, file io.Reader) (*string, []byte, *string, *gut.ErrorInstance) {
	return r.file.Upload(name, directory, attribute, file)
}

func (r *Plugger) Bind(typ strawhouse.FeedType, dir string, handler func(resp any)) uint64 {
	return r.eventfeed.Bind(typ, dir, handler)
}

func (r *Plugger) Unbind(typ strawhouse.FeedType, dir string, id uint64) {
	r.eventfeed.Unbind(typ, dir, id)
}
