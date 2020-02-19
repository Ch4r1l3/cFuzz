package controller

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/master/server/service"
	"github.com/Ch4r1l3/cFuzz/master/server/service/kubernetes"
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
	user, err := service.GetUserByID(id)
	if err != nil {
		utils.DBError(c)
		return
	}
	if user == nil {
		utils.NotFound(c)
		return
	}
	c.JSON(http.StatusOK, *user)
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
	users, count, err := service.GetNormalUserCombine(offset, limit, name)
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
	// swagger:operation POST /user/login user userLogin
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
	ok, err := service.VerifyUser(req.Username, req.Password)
	if err != nil {
		utils.InternalErrorWithMsg(c, err.Error())
		return
	}
	if ok {
		user, err := service.GetUserByUsername(req.Username)
		if err != nil {
			utils.DBError(c)
			return
		}
		token, err := service.CreateToken(user.ID, user.IsAdmin)
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
	// swagger:operation GET /user/logout user userLogout
	// user Logout
	//
	// user Logout
	// ---
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//      description: "logout"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "session",
		Value:   "",
		Path:    "/",
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
	if service.IsUserExistsByUsername(req.Username) {
		utils.BadRequestWithMsg(c, "username exist")
		return
	}
	validate := validator.New()
	errs := validate.Var(req.Password, "min=6,max=18,printascii")
	if errs != nil {
		utils.BadRequestWithMsg(c, errs.Error())
		return
	}
	err := service.CreateUser(req.Username, req.Password, false)
	if err != nil {
		utils.InternalErrorWithMsg(c, err.Error())
		return
	}
	user, err := service.GetUserByUsername(req.Username)
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
	if uint64(c.GetInt64("id")) == user.ID {
		if req.OldPassword == "" {
			utils.BadRequestWithMsg(c, "oldpassword empty")
			return
		}
		if utils.GetEncryptPassword(req.OldPassword, user.Salt) != user.Password {
			utils.BadRequestWithMsg(c, "oldpassword wrong")
			return
		}
	}
	errs := validate.Var(req.NewPassword, "min=6,max=18,printascii")
	if errs != nil {
		utils.BadRequestWithMsg(c, errs.Error())
		return
	}
	user.Password = utils.GetEncryptPassword(req.NewPassword, user.Salt)
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
	var tasks []models.Task
	if err = service.GetObjectsByUserID(&tasks, user.ID); err != nil {
		utils.InternalErrorWithMsg(c, err.Error())
	}
	for _, t := range tasks {
		if t.IsRunning() {
			kubernetes.DeleteContainerByTaskID(t.ID)
		}
	}
	if err = service.DeleteUserByID(user.ID); err != nil {
		utils.InternalErrorWithMsg(c, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}
