package users

import "mime/multipart"

type UserPhotoUpdate struct {
	UserID int64                 `form:"user_id" binding:"-"`
	Photo  *multipart.FileHeader `form:"photo" binding:"-"`
}
