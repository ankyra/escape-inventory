package model

import (
    "github.com/ankyra/escape-registry/dao"
)


func Registry(typ, name string) ([]string, error) {
    println("Registry")
    println(typ)
    println(name)
    if typ == "" {
        return dao.GetReleaseTypes()
    }
    result := []string{}
    return result, nil
}
