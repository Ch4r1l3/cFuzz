package controller

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

// swagger:model
type DeploymentReq struct {
	// example: test-image
	// required: true
	Name string `json:"name" binding:"required"`

	// example: 123
	Content string `json:"content"`
}

// swagger:model
type DeploymentSimp struct {
	// example: 1
	ID uint64 `json:"id"`

	// example: test-image
	Name string `json:"name"`
}

type DeploymentController struct{}

// List Deployment
func (dc *DeploymentController) List(c *gin.Context) {
	// swagger:operation GET /deployment deployment listDeployment
	// list all deployment
	//
	// list all deployment
	// ---
	// produces:
	// - application/json
	//
	// parameters:
	// - name: offset
	//   in: query
	//   type: integer
	// - name: limit
	//   in: query
	//   type: integer
	//
	// responses:
	//   '200':
	//      schema:
	//        type: array
	//        items:
	//          "$ref": "#/definitions/Deployment"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	var err error
	var deployments []models.Deployment
	if !c.GetBool("pagination") {
		err = models.GetObjects(&deployments)
	} else {
		offset := c.GetInt("offset")
		limit := c.GetInt("limit")
		err = models.GetObjectsPagination(&deployments, offset, limit)
	}
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, deployments)
}

// Count of Deployment
func (dc *DeploymentController) Count(c *gin.Context) {
	// swagger:operation GET /deployment/count deployment countDeployment
	// count of deployment
	//
	// count of deployment
	// ---
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//      schema:
	//        "$ref": "#/definitions/CountResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	count, err := models.GetCount(&models.Deployment{})
	if err != nil {
		utils.DBError(c)
	}
	c.JSON(http.StatusOK, CountResp{
		Count: count,
	})
}

// Simplification List Of Deployment
func (dc *DeploymentController) SimpList(c *gin.Context) {
	// swagger:operation GET /deployment/simplist deployment simlistDeployment
	// simplification list of deployment
	//
	// simplification deployment, list all id, name of deployment
	// ---
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//      schema:
	//        type: array
	//        items:
	//          "$ref": "#/definitions/DeploymentSimp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	var deployments []models.Deployment
	var err error
	if !c.GetBool("pagination") {
		err = models.GetObjects(&deployments)
	} else {
		offset := c.GetInt("offset")
		limit := c.GetInt("limit")
		err = models.GetObjectsPagination(&deployments, offset, limit)
	}
	if err != nil {
		utils.DBError(c)
		return
	}
	deploymentSimps := []DeploymentSimp{}
	for _, deployment := range deployments {
		deploymentSimps = append(deploymentSimps, DeploymentSimp{
			ID:   deployment.ID,
			Name: deployment.Name,
		})
	}
	c.JSON(http.StatusOK, deploymentSimps)
}

// Retrieve Deployment
func (dc *DeploymentController) Retrieve(c *gin.Context, id uint64) {
	// swagger:operation GET /deployment/{id} deployment retrieveDeployment
	// retrieve deployment
	//
	// retrieve deployment
	// ---
	// produces:
	// - application/json
	//
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   type: integer
	//
	// responses:
	//   '200':
	//      schema:
	//        type: array
	//        items:
	//          "$ref": "#/definitions/Deployment"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '404':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	var deployment models.Deployment
	err := models.GetObjectByID(&deployment, id)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			utils.NotFound(c)
			return
		}
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, deployment)
}

// Create Deployment
func (dc *DeploymentController) Create(c *gin.Context) {
	// swagger:operation POST /deployment deployment createDeployment
	// create deployment
	//
	// create deployment
	// ---
	// produces:
	// - application/json
	//
	// parameters:
	// - name: DeploymentReq
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/DeploymentReq"
	//
	// responses:
	//   '201':
	//      schema:
	//        "$ref": "#/definitions/Deployment"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	var req DeploymentReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequestWithMsg(c, err.Error())
		return
	}
	deployment := models.Deployment{
		Name:    req.Name,
		Content: req.Content,
	}
	err = models.InsertObject(&deployment)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusCreated, deployment)
}

// Update Deployment
func (dc *DeploymentController) Update(c *gin.Context) {
	// swagger:operation PUT /deployment/{id} deployment updateDeployment
	// update deployment
	//
	// update deployment
	// ---
	// produces:
	// - application/json
	//
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   type: integer
	// - name: DeploymentReq
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/DeploymentReq"
	//
	// responses:
	//   '201':
	//      schema:
	//        "$ref": "#/definitions/Deployment"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '404':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	var uriReq UriIDReq
	err := c.ShouldBindUri(&uriReq)
	if err != nil {
		utils.BadRequestWithMsg(c, err.Error())
		return
	}
	var req DeploymentReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequestWithMsg(c, err.Error())
		return
	}
	var deployment models.Deployment
	err = models.GetObjectByID(&deployment, uriReq.ID)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			utils.NotFound(c)
			return
		}
		utils.DBError(c)
		return
	}
	deployment.Name = req.Name
	deployment.Content = req.Content
	if err = models.DB.Save(&deployment).Error; err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusCreated, "")
}

func (dc *DeploymentController) Destroy(c *gin.Context) {
	// swagger:operation DELETE /deployment/{id} deployment deleteDeployment
	// delete deployment
	//
	// delete deployment
	// ---
	// produces:
	// - application/json
	//
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   type: integer
	//
	// responses:
	//   '204':
	//      schema:
	//        "$ref": "#/definitions/Deployment"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '404':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	var uriReq UriIDReq

	err := c.ShouldBindUri(&uriReq)
	if err != nil {
		utils.BadRequestWithMsg(c, err.Error())
		return
	}
	err = models.DeleteObjectByID(models.Deployment{}, uriReq.ID)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			utils.NotFound(c)
			return
		}
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusNoContent, "")
}
