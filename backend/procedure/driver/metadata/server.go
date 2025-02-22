package metadata

import (
	"github.com/strawst/strawhouse-go/pb"
	"google.golang.org/grpc"
	"strawhouse-backend/common/config"
	"strawhouse-backend/common/pogreb"
	"strawhouse-backend/util/eventfeed"
	"strawhouse-backend/util/filepath"
)

type Server struct {
	pb.UnimplementedDriverMetadataServer
	Config    *config.Config
	Pogreb    *pogreb.Pogreb
	Filepath  *filepath.Filepath
	EventFeed *eventfeed.EventFeed
}

func Register(registrar *grpc.Server, config *config.Config, pogreb *pogreb.Pogreb, filepath *filepath.Filepath, eventfeed *eventfeed.EventFeed) {
	server := &Server{
		Config:    config,
		Pogreb:    pogreb,
		Filepath:  filepath,
		EventFeed: eventfeed,
	}

	pb.RegisterDriverMetadataServer(registrar, server)
}
