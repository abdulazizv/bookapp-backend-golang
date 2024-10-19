package models

import "mime/multipart"

type File struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type FileResponse struct {
	Url string `json:"file_url"`
}
