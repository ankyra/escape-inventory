package model

import (
    "github.com/ankyra/escape-registry/dao"
)

type UserError struct {
    Msg string
}

func NewUserError(err error) error {
    if dao.IsNotFound(err) {
        return err
    }
    return UserError{Msg: err.Error()}
}

func (u UserError) Error() string {
    return u.Msg
}

func IsUserError(err error) bool {
    switch err.(type) {
    case UserError:
        return true
    }
    return false
}
