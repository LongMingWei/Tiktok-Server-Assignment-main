package main

import (
	"context"
	"fmt"
	rpc "github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc/imservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	_ "github.com/go-sql-driver/mysql"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
)

var (
	dbClient = &SQLClient{} // Create an SQLClient with global visibility in the 'main' scope
)

func main() {
	ctx := context.Background()
	//time.Sleep(50 * time.Second)
	err := dbClient.InitClient(ctx, "user", "a", "tiktok_server_assignment-main-db-1", "allMessages")
	if err != nil {
		errMsg := fmt.Sprintf("failed to init SQL client, err: %v", err)
		log.Fatal(errMsg)
	}

	r, err := etcd.NewEtcdRegistry([]string{"etcd:2379"})
	if err != nil {
		log.Fatal(err)
	}

	svr := rpc.NewServer(new(IMServiceImpl), server.WithRegistry(r), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: "demo.rpc.server",
	}))

	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
