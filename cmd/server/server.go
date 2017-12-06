package main

import (
	"archive/zip"
	"bytes"
	"log"
	"net"
	"path/filepath"

	pb "github.com/sago35/grpczip-sample"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type grpcServer struct {
}

func (s *grpcServer) Grpczip(ctx context.Context, r *pb.Request) (*pb.Response, error) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	for _, file := range r.GetFiles() {
		f, err := w.Create(filepath.Base(file.GetFilename()))
		if err != nil {
			return nil, err
		}
		_, err = f.Write(file.GetData())
		if err != nil {
			return nil, err
		}
	}
	err := w.Close()
	if err != nil {
		return nil, err
	}

	res := &pb.Response{
		ZipFile: &pb.File{
			Filename: r.GetZipFilename(),
			Data:     buf.Bytes(),
		},
	}
	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", ":11111")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer(grpc.MaxMsgSize(0xFFFFFFFF))

	pb.RegisterGrpczipServer(srv, &grpcServer{})
	srv.Serve(lis)
}
