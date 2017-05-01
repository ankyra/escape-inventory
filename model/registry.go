package model

import (
    "github.com/ankyra/escape-registry/dao"
    . "github.com/ankyra/escape-registry/dao/types"
)


func Registry(typ, name string) ([]string, error) {
    types, err := dao.GetReleaseTypes()
    if err != nil {
        return nil, err
    }
    if typ == "" {
        return types, nil
    }
    typeFound := false
    for _, t := range types {
        if t == typ {
            typeFound = true
        }
    }
    if !typeFound {
        return nil, NotFound
    }
    if name == "" {
        return dao.GetApplicationsByType(typ)
    }
    app, err := dao.GetApplication(typ, name)
    if err != nil {
        return nil, err
    }
    return app.FindAllVersions()
}
