package LM_V3

import (
	"LibraryManagementV1/LM_V3/model"
	"LibraryManagementV1/LM_V3/router"
	"LibraryManagementV1/LM_V3/tools"
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
