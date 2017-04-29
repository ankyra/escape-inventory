package model

import (
    "strings"
    "strconv"
)

type SemanticVersion struct {
    versionParts []string
}

func NewSemanticVersion(v string) *SemanticVersion {
    return &SemanticVersion{
        versionParts: strings.Split(v, "."),
    }
}

func (s *SemanticVersion) IncrementSmallest() error {
    lastIx := len(s.versionParts) - 1
    last := s.versionParts[lastIx]
    lastI, err := strconv.Atoi(last)
    if err != nil {
        return err
    }
    lastI += 1
    s.versionParts[lastIx] = strconv.Itoa(lastI)
    return nil
}

func (s *SemanticVersion) ToString() string {
    return strings.Join(s.versionParts, ".")
}

func (s *SemanticVersion) LessOrEqual(o *SemanticVersion) bool {
    ix := 0
    for true {
        if len(s.versionParts) == ix {
            return len(o.versionParts) != ix
        }
        if len(o.versionParts) == ix {
            return false
        }
        mine := s.versionParts[ix]
        theirs := o.versionParts[ix]

        mineInt, mineIntErr := strconv.Atoi(mine)
        if mineIntErr != nil {
            _, theirsIntErr := strconv.Atoi(theirs)
            return theirsIntErr == nil
        }
        theirsInt, theirsIntErr := strconv.Atoi(theirs)
        if theirsIntErr != nil {
            return false
        }
        if mineInt < theirsInt {
            return true
        }
        if mineInt > theirsInt {
            return false
        }
        ix += 1
    }
    return false
}
