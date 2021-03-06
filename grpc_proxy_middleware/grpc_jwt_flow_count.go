package grpc_proxy_middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_gateway/dao"
	"go_gateway/public"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

func GrpcJwtFlowCountMiddleware(serviceDetail *dao.ServiceDetail) func(srv interface{}, ss grpc.ServerStream,
	info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errors.New("miss metadata from context")
		}

		appInfos := md.Get("app")
		//!!!
		fmt.Println("----222------")
		fmt.Println(md)
		fmt.Println(appInfos)
		if len(appInfos) == 0 {
			if err := handler(srv, ss); err != nil {
				log.Printf("RPC failed with error %v\n", err)
				return err
			}
			return nil
		}
		appInfo := &dao.App{}
		if err := json.Unmarshal([]byte(appInfos[0]), appInfo); err != nil {
			return err
		}

		appCounter, err := public.FlowCounterHandler.GetCounter(
			public.FlowAppPrefix + appInfo.AppID)
		if err != nil {
			return err
		}
		appCounter.Increase()
		if appInfo.Qpd > 0 && appCounter.TotalCount > appInfo.Qpd {
			return errors.New(fmt.Sprintf("租户日请求量限流 limit:%v current:%v",
				appInfo.Qpd, appCounter.TotalCount))
		}
		if err := handler(srv, ss); err != nil {
			log.Printf("GrpcJwtFlowCountMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}
