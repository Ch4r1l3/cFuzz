package controller

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type UserController struct{}

// swagger:model
type UserReq struct {
	// example: 123
	Username string `json:"username"`
	// example: 1233
	Password string `json:"password"`
}

// swagger:model
type LoginResp struct {
	// example: 111xxxx
	Token string `json:"token"`
}

// get current user status
func (uc *UserController) Status(c *gin.Context) {
	// swagger:operation GET /user/status user userStatus
	// get current user status
	//
	// get current user status
	// ---
	// produces:
	// - application/json
	//
	// responses:
	//   '202':
	//      schema:
	//        "$ref": "#/definitions/User"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	id := uint64(c.GetInt64("id"))
	var user models.User
	if models.GetObjectByID(&user, id) != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, user)
}

// list user
func (us *UserController) List(c *gin.Context) {
	// swagger:operation GET /user user listUser
	// list user
	//
	// list user
	// ---
	// produces:
	// - application/json
	//
	// responses:
	//   '202':
	//      schema:
	//        type: array
	//        items:
	//          "$ref": "#/definitions/User"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	users, err := models.GetNormalUser()
	if err != nil {
		utils.InternalErrorWithMsg(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, users)
}

// user login
func (us *UserController) Login(c *gin.Context) {
	// swagger:operation GET /user/login user userLogin
	// user Login
	//
	// user Login
	// ---
	// produces:
	// - application/json
	//
	// responses:
	//   '202':
	//      schema:
	//        "$ref": "#/definitions/LoginResp"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	var req UserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestWithMsg(c, err.Error())
		return
	}
	ok, err := models.VerifyUser(req.Username, req.Password)
	if err != nil {
		utils.InternalErrorWithMsg(c, err.Error())
		return
	}
	if ok {
		user, err := models.GetUserByUsername(req.Username)
		if err != nil {
			utils.DBError(c)
			return
		}
		token, err := models.CreateToken(user.ID, user.IsAdmin)
		if err != nil {
			utils.InternalErrorWithMsg(c, err.Error())
			return
		}
		http.SetCookie(c.Writer, &http.Cookie{
			Name:    "session",
			Value:   token,
			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour),
		})
		c.JSON(http.StatusOK, LoginResp{
			Token: token,
		})
	} else {
		utils.ForbiddenWithMsg(c, "login failed")
	}
}

func (us *UserController) Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "session",
		Value:   "",
		Expires: time.Now(),
	})
	c.String(http.StatusOK, "")
}

// createUser
func (us *UserController) Create(c *gin.Context) {
	// swagger:operation POST /user user createUser
	// create user
	//
	// create user
	// ---
	// produces:
	// - application/json
	//
	// responses:
	//   '202':
	//      schema:
	//        "$ref": "#/definitions/User"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	var req UserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestWithMsg(c, err.Error())
		return
	}
	ok, err := models.IsUsernameExists(req.Username)
	if err != nil {
		utils.DBError(c)
		return
	}
	if ok {
		utils.BadRequestWithMsg(c, "username exist")
		return
	}
	err = models.CreateUser(req.Username, req.Password, false)
	if err != nil {
		utils.InternalErrorWithMsg(c, err.Error())
		return
	}
	user, err := models.GetUserByUsername(req.Username)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, user)
}

// delete user
func (us *UserController) Delete(c *gin.Context) {
	// swagger:operation DELETE /user/{id} user deleteUser
	// create user
	//
	// create user
	// ---
	// produces:
	// - application/json
	//
	// responses:
	//   '202':
	//      schema:
	//        "$ref": "#/definitions/User"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	var user models.User
	err := getObject(c, &user)
	if err != nil {
		return
	}
	if err = models.DeleteObjectByID(&models.User{}, user.ID); err != nil {
		utils.InternalErrorWithMsg(c, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}
