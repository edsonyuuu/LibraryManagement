package LM_V2

import (
	"LibraryManagementV1/LM_V2/model"
	"LibraryManagementV1/LM_V2/router"
)

func Start() {
	defer func() {
		model.Close()
	}()

	model.MySql()
	r := router.New()
	_ = r.Run(":8083")
}
