package tcp_proxy_middleware

import (
	"fmt"
	"go_gateway/dao"
	"go_gateway/public"
	"strings"
)

func TCPFlowLimitMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		serverInterface := c.Get("service")
		if serverInterface == nil {
			c.conn.Write([]byte("get service empty"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		//服务限流
		if serviceDetail.AccessControl.ServiceFlowLimit != 0 {
			serviceLimiter, err := public.FlowLimiterHandler.GetLimiter(
				public.FlowServicePrefix+serviceDetail.Info.ServiceName,
				float64(serviceDetail.AccessControl.ServiceFlowLimit))
			if err != nil {
				c.conn.Write([]byte(err.Error()))
				c.Abort()
				return
			}
			if !serviceLimiter.Allow() {
				c.conn.Write([]byte(err.Error()))
				c.Abort()
				return
			}

		}

		//ip:port
		splits := strings.Split(c.conn.RemoteAddr().String(), ":")
		clientIP := ""
		if len(splits) == 2 {
			clientIP = splits[0]
		}

		//客户端限流
		if serviceDetail.AccessControl.ClientIPFlowLimit > 0 {
			clientLimiter, err := public.FlowLimiterHandler.GetLimiter(
				public.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+clientIP,
				float64(serviceDetail.AccessControl.ClientIPFlowLimit))
			if err != nil {
				c.conn.Write([]byte(err.Error()))
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				c.conn.Write([]byte(fmt.Sprintf("%v flow limit %v", clientIP,
					serviceDetail.AccessControl.ClientIPFlowLimit)))
				c.Abort()
				return
			}

		}

		c.Next()
	}
}
