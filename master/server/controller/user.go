package controller

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

type UserController struct{}

// swagger:model
type UserReq struct {
	// example: 123
	Username string `json:"username" binding:"required"`
	// example: 1233
	Password string `json:"password" binding:"required"`
}

// swagger:model
type LoginResp struct {
	// example: 111xxxx
	Token string `json:"token"`
}

// swagger:model
type UserListResp struct {
	Data []models.User `json:"data"`
	CountResp
}

// swagger:model
type UserUpdateReq struct {
	// example: 1233445
	OldPassword string `json:"oldPassword"`
	// example: 123456
	NewPassword string `json:"newPassword" binding:"required"`
}

// get current user status
func (uc *UserController) Info(c *gin.Context) {
	// swagger:operation GET /user/info user userStatus
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
	//        "$ref": "#/definitions/UserListResp"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	offset := c.GetInt("offset")
	limit := c.GetInt("limit")
	name := c.Query("name")
	users, count, err := models.GetNormalUserCombine(offset, limit, name)
	if err != nil {
		utils.InternalErrorWithMsg(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, UserListResp{
		Data: users,
		CountResp: CountResp{
			Count: count,
		},
	})
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
	validate := validator.New()
	errs := validate.Var(req.Password, "min=6,max=18,printascii")
	if errs != nil {
		utils.BadRequestWithMsg(c, errs.Error())
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
	c.JSON(http.StatusCreated, user)
}

// update user
func (us *UserController) Update(c *gin.Context) {
	// swagger:operation UPDATE /user/{id} user updateUser
	// update user
	//
	// update user
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
	if uint64(c.GetInt64("id")) != user.ID && !c.GetBool("isAdmin") {
		utils.Forbidden(c)
		return
	}
	var req UserUpdateReq
	if err = c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestWithMsg(c, err.Error())
		return
	}
	validate := validator.New()
	if !c.GetBool("isAdmin") {
		if req.OldPassword == "" {
			utils.BadRequestWithMsg(c, "oldpassword empty")
			return
		}
		if models.GetEncryptPassword(req.OldPassword, user.Salt) != user.Password {
			utils.BadRequestWithMsg(c, "oldpassword wrong")
			return
		}
	}
	errs := validate.Var(req.NewPassword, "min=6,max=18,printascii")
	if errs != nil {
		utils.BadRequestWithMsg(c, errs.Error())
		return
	}
	user.Password = models.GetEncryptPassword(req.NewPassword, user.Salt)
	if err = models.DB.Save(user).Error; err != nil {
		utils.DBError(c)
		return
	}
	c.String(http.StatusCreated, "")
}

// delete user
func (us *UserController) Delete(c *gin.Context) {
	// swagger:operation DELETE /user/{id} user deleteUser
	// delete user
	//
	// delete user
	// ---
	// produces:
	// - application/json
	//
	// responses:
	//   '204':
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
