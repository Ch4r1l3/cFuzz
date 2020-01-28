package controller

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
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
type DeploymentSim struct {
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

// Summary Of Deployment
func (dc *DeploymentController) Summary(c *gin.Context) {
	// swagger:operation GET /deployment/summary deployment summaryDeployment
	// summary deployment
	//
	// summary deployment, list all id, name of deployment, and total count of deployment
	// ---
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//      schema:
	//        type: array
	//        items:
	//          "$ref": "#/definitions/DeploymentSim"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	var deployments []models.Deployment
	err := models.GetObjects(&deployments)
	if err != nil {
		utils.DBError(c)
		return
	}
	deploymentSims := []DeploymentSim{}
	for _, deployment := range deployments {
		deploymentSims = append(deploymentSims, DeploymentSim{
			ID:   deployment.ID,
			Name: deployment.Name,
		})
	}
	c.JSON(http.StatusOK, deploymentSims)
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
	//   '200':
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
		utils.BadRequest(c)
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
	c.JSON(http.StatusOK, deployment)
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
	//   '204':
	//      schema:
	//        "$ref": "#/definitions/Deployment"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	var uriReq UriIDReq
	err := c.ShouldBindUri(&uriReq)
	if err != nil {
		utils.BadRequest(c)
		return
	}
	var req DeploymentReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(c)
		return
	}
	var deployment models.Deployment
	err = models.GetObjectByID(&deployment, uriReq.ID)
	if err != nil {
		utils.DBError(c)
		return
	}
	deployment.Name = req.Name
	deployment.Content = req.Content
	if err = models.DB.Save(&deployment).Error; err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusNoContent, "")
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
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	var uriReq UriIDReq

	err := c.ShouldBindUri(&uriReq)
	if err != nil {
		utils.BadRequest(c)
		return
	}
	err = models.DeleteObjectByID(models.Deployment{}, uriReq.ID)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusNoContent, "")
}
