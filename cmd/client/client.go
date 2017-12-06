package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/net/context"

	pb "github.com/sago35/grpczip-sample"

	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage : %s [zipFilename] [input...]", os.Args[0])
	}
	zipFilename := os.Args[1]
	files := os.Args[2:]

	conn, err := grpc.Dial("127.0.0.1:11111", grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(0xFFFFFFFF)))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewGrpczipClient(conn)

	req, err := makeRequest(zipFilename, files)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Grpczip(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	z := resp.GetZipFile()
	w, err := os.Create(z.GetFilename())
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	r := bytes.NewReader(z.GetData())
	size, err := io.Copy(w, r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s (%d bytes)\n", z.GetFilename(), size)
}

func makeRequest(zipFilename string, files []string) (*pb.Request, error) {
	ret := &pb.Request{
		ZipFilename: zipFilename,
		Files:       nil,
	}

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		ret.Files = append(ret.Files, &pb.File{
			Filename: file,
			Data:     content,
		})
	}
	return ret, nil
}
