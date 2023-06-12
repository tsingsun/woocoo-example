package main

import (
	"fmt"
	"github.com/tsingsun/woocoo"
	"github.com/tsingsun/woocoo-example/grpc/polaris/hellopb"
	"github.com/tsingsun/woocoo/rpc/grpcx"
	"log"
	"net/http"

	_ "github.com/tsingsun/woocoo/contrib/polarismesh"
)

func main() {
	app := woocoo.New()
	client := grpcx.NewClient(app.AppConfiguration().Sub("grpc"))
	conn, err := client.Dial("")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	hello := hellopb.NewHelloServiceClient(conn)
	http.HandleFunc("/echo", func(rw http.ResponseWriter, r *http.Request) {
		for i := 0; i < 10; i++ {
			resp, err := hello.SayHello(r.Context(), &hellopb.HelloRequest{Greeting: "woocoo"})
			if err != nil {
				log.Printf("[error] fail to say, err is %v", err)
				rw.WriteHeader(http.StatusInternalServerError)
				_, _ = rw.Write([]byte(fmt.Sprintf("[error] fail to say, err is %v", err)))
				_, _ = rw.Write([]byte(err.Error()))
				_, _ = rw.Write([]byte("\n"))
				continue
			}
			rw.WriteHeader(http.StatusOK)
			_, _ = rw.Write([]byte(resp.Reply))
			_, _ = rw.Write([]byte("\n"))
		}
	})

	log.Printf("start run web server, port : %d", 12000)

	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", 12000), nil); err != nil {
		log.Fatalf("[ERROR]fail to run webServer, err is %v", err)
	}
}
