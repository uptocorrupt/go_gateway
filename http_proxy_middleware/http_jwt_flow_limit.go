package http_proxy_middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go_gateway/dao"
	"go_gateway/middleware"
	"go_gateway/public"
)

//单机并发数拦截，并非日请求量限流

func HTTPJwtFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}
		appInfo := appInterface.(*dao.App)

		if appInfo.Qps != 0 {
			clientLimiter, err := public.FlowLimiterHandler.GetLimiter(
				public.FlowAppPrefix+appInfo.AppID+"_"+c.ClientIP(),
				float64(appInfo.Qps))
			if err != nil {
				middleware.ResponseError(c, 5001, err)
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				middleware.ResponseError(c, 5002, errors.New(fmt.Sprintf(
					"%v flow limit %v",
					c.ClientIP(), appInfo.Qps)))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
