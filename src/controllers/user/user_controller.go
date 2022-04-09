package user

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_users-api/src/domain/users"
	"github.com/superbkibbles/realestate_users-api/src/services"
)

var UserController userControllerInterface = &userController{}

type userControllerInterface interface {
	Get(*gin.Context)
	GetByID(*gin.Context)
	Create(*gin.Context)
	UpdateUser(*gin.Context)
	UpdateUserPhoto(*gin.Context)
	Delete(*gin.Context)
	LikeProperty(*gin.Context)
	// Update(*gin.Context)
	// Delete(*gin.Context)
}

type userController struct{}

func getUserId(userIdParam string) (int64, rest_errors.RestErr) {
	userID, userErr := strconv.ParseInt(userIdParam, 10, 64)
	if userErr != nil {
		return 0, rest_errors.NewBadRequestErr("user ID should be a number")
	}
	return userID, nil
}

func (*userController) Get(c *gin.Context) {
	// Authenticatae Against auth api

	// Call service to get All Users
	users, err := services.UserService.Get()
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusFound, users.Marshal(false))
}

func (*userController) GetByID(c *gin.Context) {
	id, err := getUserId(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	user, err := services.UserService.GetByID(id)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, user.Marshal(false))
}

func (*userController) LikeProperty(c *gin.Context) {
	var likeProperty users.LikePrpertyReq

	if err := c.ShouldBindJSON(&likeProperty); err != nil {
		restErr := rest_errors.NewBadRequestErr("Invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	err := services.UserService.LikeProperty(likeProperty)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, "Successfully liked property")

}

func (*userController) Create(c *gin.Context) {
	var user users.UserForm

	if err := c.ShouldBind(&user); err != nil {
		restErr := rest_errors.NewBadRequestErr("Invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		if !strings.Contains(err.Error(), "http: no such file") {
			restErr := rest_errors.NewBadRequestErr("Invalid file type")
			c.JSON(restErr.Status(), restErr)
			fmt.Println(err.Error())
			return
		}
	}

	res, saveErr := services.UserService.Create(user, header, file)
	if saveErr != nil {
		c.JSON(saveErr.Status(), saveErr)
		return
	}

	c.JSON(http.StatusCreated, res.Marshal(false))
}

func (*userController) UpdateUser(c *gin.Context) {
	var user users.User
	userID, err := getUserId(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status(), err)
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestErr("Invalid JSON Body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	user.Id = userID

	updatedUser, err := services.UserService.UpdateUser(user)

	c.JSON(http.StatusOK, updatedUser)
}

func (*userController) UpdateUserPhoto(c *gin.Context) {
	userID, err := getUserId(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	var photoUpdateRequest users.UserPhotoUpdate

	if err := c.ShouldBind(&photoUpdateRequest); err != nil {
		restErr := rest_errors.NewBadRequestErr("Invalid Body JSON")
		c.JSON(restErr.Status(), restErr)
		return
	}

	photoUpdateRequest.UserID = userID

	file, header, photoErr := c.Request.FormFile("photo")
	if photoErr != nil {
		restErr := rest_errors.NewBadRequestErr("Invalid file type")
		c.JSON(restErr.Status(), restErr)
		return
	}

	user, err := services.UserService.UpdatePhoto(photoUpdateRequest, header, file)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, user)
	// Call service to update the pic
}

func (*userController) Delete(c *gin.Context) {
	userID, err := getUserId(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	if err := services.UserService.DeleteUser(userID); err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Deleted",
	})
}

// Test
// func ServeFile(c *gin.Context) {
// 	file, err := os.Open("datasources/images/f91be0d203344f37cfa4ae6f318255ea.pdf") //Create a file
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	defer file.Close()
// 	c.Writer.Header().Add("Content-type", "application/octet-stream")
// 	// c.Writer.Header().Add("Content-type", "application/pdf")
// 	_, err = io.Copy(c.Writer, file)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// }

// // DownloadAttachmentHandler ...
// func DownloadAttachmentHandler(c *gin.Context) {
//     var request api.DownloadAttachmentRequest
//     filePath = c.Param("file")
//     file, err := os.Open(filePath) //Create a file
//     if err != nil {
//        c.JSON(http.StatusNotFound, HTTPGenericResponse{
//         Code:    http.StatusInternalServerError,
//         Message: "文件加载失败:" + err.Error(),
//        })
//     return
//     }
//     defer file.Close()
// c.Header("Content-Disposition", "attachment; filename=a.tar")
//     c.Writer.Header().Add("Content-type", "application/octet-stream")
//     _, err = io.Copy(c.Writer, file)
//     if err != nil {
//         c.JSON(http.StatusNotFound, HTTPGenericResponse{
//             Code:    http.StatusInternalServerError,
//             Message: "文件加载失败:" + err.Error(),
//         })
//     return
//     }
// }

// c.File("datasources/images/f91be0d203344f37cfa4ae6f3182ddd55ea.pdf")
// c.String(200, "Not found")
