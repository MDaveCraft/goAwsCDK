package interfaces

import "lambda-func/types"

type IDbCLient interface {
	InsertUser(user types.User) error
	UserExists(username string) (bool, error)
	GetUser(username string) (types.User, error)
}
