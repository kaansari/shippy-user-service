package auth

import (
	"github.com/jinzhu/gorm"

	uuid "github.com/satori/go.uuid"
)

/*
"BeforeCreate Assign uuid is it is not passed by the client"
*/
func (model *User) BeforeCreate(scope *gorm.Scope) (err error) {

	if len(model.Id) == 0 {
		uuid := uuid.NewV4()
		scope.SetColumn("Id", uuid.String())
	}

	return

}
