package application

import "github.com/superbkibbles/realestate_users-api/src/controllers/user"

func mapUrls() {
	// Get all users
	router.GET("/api/users", user.UserController.Get)
	// Create new user
	router.POST("/api/users", user.UserController.Create)
	// Get user by id
	router.GET("/api/users/:user_id", user.UserController.GetByID)
	// update user
	router.PATCH("/api/users/:user_id", user.UserController.UpdateUser)
	// Update User Photo
	router.PATCH("/api/users/photo/:user_id", user.UserController.UpdateUserPhoto)
	// deactivated user
	router.DELETE("/api/users/:user_id", user.UserController.Delete)
}
