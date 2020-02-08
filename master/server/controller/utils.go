package controller

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func getObject(c *gin.Context, obj interface{}) error {
	var uriReq UriIDReq
	err := c.ShouldBindUri(&uriReq)
	if err != nil {
		utils.BadRequestWithMsg(c, err.Error())
		return err
	}
	err = models.GetObjectByID(obj, uriReq.ID)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			utils.NotFound(c)
			return err
		}
		utils.DBError(c)
		return err
	}
	return nil
}

func getList(c *gin.Context, objs interface{}) (int, error) {
	var count int
	var err error
	offset := c.GetInt("offset")
	limit := c.GetInt("limit")
	name := c.Query("name")
	userID := uint64(c.GetInt64("id"))
	isAdmin := c.GetBool("isAdmin")
	count, err = models.GetObjectCombine(objs, offset, limit, name, userID, isAdmin)
	if err != nil {
		utils.DBError(c)
	}
	return count, err
}
