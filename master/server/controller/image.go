package controller

import (
	"errors"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/master/server/service"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// swagger:model
type ImageReq struct {
	// example: test-image
	// required: true
	Name string `json:"name" binding:"required"`

	// example: true
	IsDeployment bool `json:"isDeployment"`

	// example: 123
	Content string `json:"content" binding:"required"`
}

// swagger:model
type ImageCombine struct {
	Data []models.Image `json:"data"`
	CountResp
}

type ImageController struct{}

func getImage(c *gin.Context) (*models.Image, error) {
	var image models.Image
	err := getObject(c, &image)
	if err != nil {
		return nil, err
	}
	if image.UserID != uint64(c.GetInt64("id")) && !c.GetBool("isAdmin") {
		utils.Forbidden(c)
		return nil, errors.New("no permission")
	}
	return &image, nil
}

// List Image
func (dc *ImageController) List(c *gin.Context) {
	// swagger:operation GET /image image listImage
	// list all image
	//
	// list all image
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
	// - name: name
	//   in: query
	//   type: string
	//
	// responses:
	//   '200':
	//      schema:
	//        "$ref": "#/definitions/ImageCombine"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	var images []models.Image
	count, err := getList(c, &images)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, ImageCombine{
		Data: images,
		CountResp: CountResp{
			Count: count,
		},
	})
}

// Retrieve Image
func (dc *ImageController) Retrieve(c *gin.Context) {
	// swagger:operation GET /image/{id} image retrieveImage
	// retrieve image
	//
	// retrieve image
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
	//          "$ref": "#/definitions/Image"
	//   '400':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '404':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	image, err := getImage(c)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, *image)
}

// Create Image
func (dc *ImageController) Create(c *gin.Context) {
	// swagger:operation POST /image image createImage
	// create image
	//
	// create image
	// ---
	// produces:
	// - application/json
	//
	// parameters:
	// - name: ImageReq
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/ImageReq"
	//
	// responses:
	//   '201':
	//      schema:
	//        "$ref": "#/definitions/Image"
	//   '400':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	var req ImageReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequestWithMsg(c, err.Error())
		return
	}
	image := models.Image{
		Name:         req.Name,
		IsDeployment: req.IsDeployment,
		Content:      req.Content,
		UserID:       uint64(c.GetInt64("id")),
	}
	err = service.CreateImage(&image)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusCreated, image)
}

// Update Image
func (dc *ImageController) Update(c *gin.Context) {
	// swagger:operation PUT /image/{id} image updateImage
	// update image
	//
	// update image
	// ---
	// produces:
	// - application/json
	//
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   type: integer
	// - name: ImageReq
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/ImageReq"
	//
	// responses:
	//   '201':
	//      schema:
	//        "$ref": "#/definitions/Image"
	//   '400':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '404':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	var req ImageReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequestWithMsg(c, err.Error())
		return
	}
	image, err := getImage(c)
	if err != nil {
		return
	}
	if service.IsImageReferred(image.ID) {
		utils.BadRequestWithMsg(c, "this image is being used by task")
		return
	}
	image.Name = req.Name
	image.Content = req.Content
	image.IsDeployment = req.IsDeployment
	if err = service.UpdateImage(image); err != nil {
		utils.DBError(c)
		return
	}
	c.String(http.StatusCreated, "")
}

func (dc *ImageController) Destroy(c *gin.Context) {
	// swagger:operation DELETE /image/{id} image deleteImage
	// delete image
	//
	// delete image
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
	//        "$ref": "#/definitions/Image"
	//   '400':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '404':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	image, err := getImage(c)
	if err != nil {
		return
	}
	if service.IsImageReferred(image.ID) {
		utils.BadRequestWithMsg(c, "this image is being used by task")
		return
	}
	err = service.DeleteObjectByID(models.Image{}, image.ID)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.String(http.StatusNoContent, "")
}
