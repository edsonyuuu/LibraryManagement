package LM_V1

import (
	"LibraryManagementV1/LM_V1/model"
	"LibraryManagementV1/LM_V1/router"
)

func Start() {
	defer func() {
		model.Close()
	}()

	model.MySql()
	r := router.New()
	_ = r.Run(":8083")
}
