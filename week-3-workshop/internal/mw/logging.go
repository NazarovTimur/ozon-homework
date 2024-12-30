package mw

import (
	"context"
	"log"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"google.golang.org/grpc"
)

// для превращения protbuf структур в json используем google.golang.org/protobuf/encoding/protojson пакет а не encoding/json
func Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	raw, _ := protojson.Marshal((req).(proto.Message))
	log.Printf("request: method: %v, req: %v\n", info.FullMethod, string(raw))

	if resp, err = handler(ctx, req); err != nil {
		log.Printf("response: method: %v, err: %v\n", info.FullMethod, err)
		return
	}

	rawResp, _ := protojson.Marshal((resp).(proto.Message))
	log.Printf("response: method: %v, resp: %v\n", info.FullMethod, string(rawResp))

	return resp, err
}
