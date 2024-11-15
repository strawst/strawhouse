package eventfeed

import (
	"github.com/strawst/strawhouse-go"
	"github.com/strawst/strawhouse-go/pb"
	"strings"
	"sync"
)

type EventFeed struct {
	mu         sync.RWMutex
	increment  uint64
	getFeed    map[string]map[uint64]func(struct{})
	uploadFeed map[string]map[uint64]func(*pb.UploadFeedResponse)
	deleteFeed map[string]map[uint64]func(struct{})
}

func Init() *EventFeed {
	return &EventFeed{
		getFeed:    make(map[string]map[uint64]func(struct{})),
		uploadFeed: make(map[string]map[uint64]func(*pb.UploadFeedResponse)),
		deleteFeed: make(map[string]map[uint64]func(struct{})),
	}
}

func (r *EventFeed) Bind(typ strawhouse.FeedType, dir string, handler func(resp any)) uint64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	inc := r.increment
	r.increment++
	switch typ {
	case strawhouse.FeedTypeGet:
		if _, ok := r.getFeed[dir]; !ok {
			r.getFeed[dir] = make(map[uint64]func(struct{}))
		}
		r.getFeed[dir][inc] = any(handler).(func(struct{}))
	case strawhouse.FeedTypeUpload:
		if _, ok := r.uploadFeed[dir]; !ok {
			r.uploadFeed[dir] = make(map[uint64]func(*pb.UploadFeedResponse))
		}
		r.uploadFeed[dir][inc] = func(resp *pb.UploadFeedResponse) {
			handler(resp)
		}
	case strawhouse.FeedTypeDelete:
		if _, ok := r.deleteFeed[dir]; !ok {
			r.deleteFeed[dir] = make(map[uint64]func(struct{}))
		}
		r.deleteFeed[dir][inc] = any(handler).(func(struct{}))
	default:
		panic("invalid feed type")
	}
	return inc
}

func (r *EventFeed) Unbind(typ strawhouse.FeedType, dir string, id uint64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	switch typ {
	case strawhouse.FeedTypeGet:
		if handlers, ok := r.getFeed[dir]; ok {
			delete(handlers, id)
			if len(handlers) == 0 {
				delete(r.getFeed, dir)
			}
		}
	case strawhouse.FeedTypeUpload:
		if handlers, ok := r.uploadFeed[dir]; ok {
			delete(handlers, id)
			if len(handlers) == 0 {
				delete(r.uploadFeed, dir)
			}
		}
	case strawhouse.FeedTypeDelete:
		if handlers, ok := r.deleteFeed[dir]; ok {
			delete(handlers, id)
			if len(handlers) == 0 {
				delete(r.deleteFeed, dir)
			}
		}
	default:
		panic("invalid feed type")
	}
}

func (r *EventFeed) Fire(typ strawhouse.FeedType, dir string, request *pb.UploadFeedResponse) {
	// TODO: Optimize event binding indexing, not to iterate over all directories
	r.mu.Lock()
	defer r.mu.Unlock()
	if typ == strawhouse.FeedTypeGet {
		panic("not implemented")
	} else if typ == strawhouse.FeedTypeUpload {
		for d := range r.uploadFeed {
			if strings.HasPrefix(dir, d) {
				for _, h := range r.uploadFeed[d] {
					h(request)
				}
			}
		}
	} else if typ == strawhouse.FeedTypeDelete {
		panic("not implemented")
	}
}
