package dto

import "mime/multipart"

type FileData struct {
	File   multipart.File
	Header *multipart.FileHeader
}
