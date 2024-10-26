package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leijeng/huo-core/common/consts"
)

func GetReqId(c *gin.Context) string {
	reqId := c.GetString(consts.REQ_ID)
	if reqId == "" {
		reqId = uuid.NewString()
		c.Set(consts.REQ_ID, reqId)
	}
	return reqId
}
