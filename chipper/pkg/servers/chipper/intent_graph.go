package server

import (
	"time"

	"github.com/BrutalHex/wire-pod/chipper/pkg/logger"
	"github.com/BrutalHex/wire-pod/chipper/pkg/vtt"
	pb "github.com/digital-dream-labs/api/go/chipperpb"
)

// StreamingIntentGraph handles intent graph request streams
func (s *Server) StreamingIntentGraph(stream pb.ChipperGrpc_StreamingIntentGraphServer) error {
	recvTime := time.Now()

	req, err := stream.Recv()
	if err != nil {
		logger.Println("Intent graph error")
		logger.Println(err)

		return err
	}

	if _, err = s.intentGraph.ProcessIntentGraph(
		&vtt.IntentGraphRequest{
			Time:       recvTime,
			Stream:     stream,
			Device:     req.DeviceId,
			Session:    req.Session,
			LangString: req.LanguageCode.String(),
			FirstReq:   req,
			AudioCodec: req.AudioEncoding,
			// Mode:
		},
	); err != nil {
		logger.Println("Intent graph error")
		logger.Println(err)
		return err
	}

	return nil
}
