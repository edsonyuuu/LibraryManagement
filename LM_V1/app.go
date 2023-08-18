package LM_V1

import (
	"LibraryManagementV1/LM_V1/model"
	"LibraryManagementV1/LM_V1/router"
	"LibraryManagementV1/LM_V1/tools"
)

func Start() {
	defer func() {
		model.Close()
	}()

	model.MySql()
	tools.NewToken("")
	r := router.New()
	_ = r.Run(":8083")
}
