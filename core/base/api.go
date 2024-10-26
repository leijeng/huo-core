package base

import (
	"github.com/gin-gonic/gin"
	"github.com/leijeng/huo-core/common/utils"
	"github.com/leijeng/huo-core/core/errs"
	"net/http"
)

type BaseApi struct {
}

func (e *BaseApi) GetReqId(c *gin.Context) string {
	return utils.GetReqId(c)
}

func (e *BaseApi) Error(c *gin.Context, err error) {
	resMsg(c, FAILURE, err.Error())
}

func (e *BaseApi) Fail(c *gin.Context, code int, msg string, data ...any) {
	resMsg(c, code, msg, data...)
}

func (e *BaseApi) Code(c *gin.Context, code int) {
	resMsg(c, code, "")
}

func (e *BaseApi) Err(c *gin.Context, err errs.IError) {
	errer(c, err)
}

func (e *BaseApi) Ok(c *gin.Context, data ...any) {
	ok(c, data...)
}

func (e *BaseApi) PureOk(c *gin.Context, data any) {
	pureJSON(c, data)
}

func (e *BaseApi) OkWithAbout(c *gin.Context, data any) {
	resMsgWithAbort(c, http.StatusOK, "OK", data)
}

func (e *BaseApi) ResCustom(c *gin.Context, opts ...Option) {
	result(c, opts...)
}

func (e *BaseApi) Page(c *gin.Context, list any, total int64, page, size int) {
	pageResp(c, list, total, page, size)
}
