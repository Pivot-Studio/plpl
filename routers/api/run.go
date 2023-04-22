package api

import (
	"net/http"
	"plpl/service/compile_service"

	"github.com/gin-gonic/gin"
)

type postPLCodeReqest struct {
	Src string `json:"src" binding:"required"`
}
type postPLCodeRespone struct {
	CompileOut string `json:"compileOut"`
	Code       int    `json:"code"`
	RunOut     string `json:"runOut"`
	Session    string `json:"session"`
}

// @Summary
// @Produce  json
// @Param code of Pivot Language
// @Success 200 {object} postPLCodeRespone
// @Failure 500 {object} postPLCodeRespone
// @Router /run [post]
func PostPLCode(c *gin.Context) {

	var req *postPLCodeReqest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}

	compileOut, runOut, session := compile_service.Compile(req.Src)

	resp := postPLCodeRespone{
		CompileOut: compileOut,
		Code:       200,
		RunOut:     runOut,
		Session:    session,
	}
	c.JSON(http.StatusOK, resp)
	return
}
