/*
Copyright 2017 Ankyra

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
