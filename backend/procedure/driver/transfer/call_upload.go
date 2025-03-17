package transfer

import (
	"bytes"
	"context"
	"github.com/strawst/strawhouse-go/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (r *Server) FileUpload(ctx context.Context, req *pb.UploadRequest) (*emptypb.Empty, error) {
	// * Create io.Reader
	reader := bytes.NewReader(req.Content)

	// * Call upload service
	_, _, _, er := r.File.Upload(req.Name, req.Directory, req.Attribute, reader)
	if er != nil {
		return nil, er
	}

	return new(emptypb.Empty), nil
}
