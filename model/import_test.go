package model

import (
	. "gopkg.in/check.v1"
)

type importSuite struct{}

var _ = Suite(&importSuite{})

func (s *importSuite) Test_Import_Empty(c *C) {
    releases := []map[string]interface{}{}
    err := Import(releases)
    c.Assert(err, IsNil)
}

func (s *importSuite) Test_Import(c *C) {
    releases := []map[string]interface{}{
        map[string]interface{}{
            "type": "archive",
            "name": "import-test",
            "version": "1",
        },
    }
    err := Import(releases)
    c.Assert(err, IsNil)
    metadata, err := GetReleaseMetadata("archive-import-test-v1")
    c.Assert(err, IsNil)
    c.Assert(metadata.GetName(), Equals, "import-test")
}

func (s *importSuite) Test_Import_Ignore_Existing(c *C) {
    releases := []map[string]interface{}{
        map[string]interface{}{
            "type": "archive",
            "name": "import-exists-test",
            "version": "1",
        },
        map[string]interface{}{
            "type": "archive",
            "name": "import-exists-test",
            "version": "1",
        },
    }
    err := Import(releases)
    c.Assert(err, IsNil)
    metadata, err := GetReleaseMetadata("archive-import-exists-test-v1")
    c.Assert(err, IsNil)
    c.Assert(metadata.GetName(), Equals, "import-exists-test")
}
