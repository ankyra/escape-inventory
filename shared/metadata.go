package shared

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
    "fmt"
    "os"
)

type ReleaseMetadata interface {
	ToJson() string
    ToDict() (map[string]interface{}, error)
	WriteJsonFile(string) error
	GetDirectories() []string
	GetVersionlessReleaseId() string
	GetReleaseId() string
	AddInputVariable(* variable)
	AddOutputVariable(* variable)
    SetVariableInContext(string, string)
    AddFileWithDigest(string, string)
    SetConsumes([]string)

	GetApiVersion() string
	GetBranch() string
	GetConsumes() []string
	GetDependencies() []string
	GetDescription() string
	GetErrands() map[string]*errand
	GetFiles() map[string]string
	GetRevision() string
	GetInputs() []*variable
	GetLogo() string
	GetMetadata() map[string]string
	GetName() string
	GetOutputs() []*variable
	GetPath() string
	GetPostBuild() string
	GetPostDestroy() string
	GetPreBuild() string
	GetPreDestroy() string
	GetProvides() []string
	GetTest() string
	GetType() string
	GetVersion() string
    GetVariableContext() map[string]string
}

type variable struct {
	Id           string                 `json:"id"`
	Type         string                 `json:"type"`
	Default      interface{}            `json:"default,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Friendly     string                 `json:"friendly,omitempty"`
	Visible      bool                   `json:"visible"`
	Options      map[string]interface{} `json:"options,omitempty"`
	Sensitive    bool                   `json:"sensitive,omitempty"`
	ProducesType string                 `json:"-"`
	Items        []interface{}          `json:"items"`
}

type errand struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Script      string      `json:"script"`
	Inputs      []*variable `json:"inputs"`
}

type releaseMetadata struct {
	ApiVersion  string             `json:"api_version"`
	Branch      string             `json:"branch"`
	Consumes    []string           `json:"consumes"`
	Depends     []string           `json:"depends"`
	Description string             `json:"description"`
	Errands     map[string]*errand `json:"errands"`
	Files       map[string]string  `json:"files", {}`
	Revision    string             `json:"git_revision"`
	Inputs      []*variable        `json:"inputs"`
	Logo        string             `json:"logo"`
	Metadata    map[string]string  `json:"metadata"`
	Name        string             `json:"name"`
	Outputs     []*variable        `json:"outputs"`
	Path        string             `json:"path"`
	PostBuild   string             `json:"post_build"`
	PostDestroy string             `json:"post_destroy"`
	PreBuild    string             `json:"pre_build"`
	PreDestroy  string             `json:"pre_destroy"`
	Provides    []string           `json:"provides"`
	Test        string             `json:"test"`
	Type        string             `json:"type"`
    VariableCtx map[string]string  `json:"variable_context"`
	Version     string             `json:"version"`
}

func NewEmptyReleaseMetadata() *releaseMetadata {
	return &releaseMetadata{
		ApiVersion: "1",
		Consumes:   []string{},
		Provides:   []string{},
		Depends:    []string{},
		Files:      map[string]string{},
		Metadata:   map[string]string{},
		Errands:    map[string]*errand{},
		Inputs:     []*variable{},
		Outputs:    []*variable{},
        VariableCtx: map[string]string{},
	}
}

func NewReleaseMetadata(typ, name, version string) *releaseMetadata {
    m := NewEmptyReleaseMetadata()
    m.Type = typ
    m.Name = name
    m.Version = version
    return m
}

func NewReleaseMetadataFromJsonString(content string) (*releaseMetadata, error) {
	result := NewEmptyReleaseMetadata()
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, err
	}
    if err := validate(result) ; err != nil {
        return nil, err
    }
	return result, nil
}

func NewReleaseMetadataFromFile(metadataFile string) (*releaseMetadata, error) {
	if !PathExists(metadataFile) {
		return nil, errors.New("Release metadata file " + metadataFile + " does not exist")
	}
	content, err := ioutil.ReadFile(metadataFile)
	if err != nil {
		return nil, err
	}
    return NewReleaseMetadataFromJsonString(string(content))
}

func validate(m *releaseMetadata) error {
    if m.Type == "" {
        return fmt.Errorf("Missing type field in release metadata")
    }
    if m.Name == "" {
        return fmt.Errorf("Missing name field in release metadata")
    }
    if m.Version == "" {
        return fmt.Errorf("Missing version field in release metadata")
    }
    return nil
}

func (m *releaseMetadata) GetApiVersion() string {
	return m.ApiVersion
}
func (m *releaseMetadata) GetBranch() string {
	return m.Branch
}
func (m *releaseMetadata) SetConsumes(c []string) {
	m.Consumes = c
}
func (m *releaseMetadata) GetConsumes() []string {
	return m.Consumes
}
func (m *releaseMetadata) GetDescription() string {
	return m.Description
}
func (m *releaseMetadata) GetErrands() map[string]*errand {
	result := map[string]*errand{}
	for key, val := range m.Errands {
		result[key] = val
	}
	return result
}
func (m *releaseMetadata) GetFiles() map[string]string {
	return m.Files
}
func (m *releaseMetadata) GetInputs() []*variable {
	result := []*variable{}
	for _, i := range m.Inputs {
		result = append(result, i)
	}
	return result
}
func (m *releaseMetadata) GetRevision() string {
	return m.Revision
}
func (m *releaseMetadata) GetLogo() string {
	return m.Logo
}
func (m *releaseMetadata) GetMetadata() map[string]string {
	return m.Metadata
}
func (m *releaseMetadata) GetName() string {
	return m.Name
}
func (m *releaseMetadata) GetOutputs() []*variable {
	result := []*variable{}
	for _, i := range m.Outputs {
		result = append(result, i)
	}
	return result
}
func (m *releaseMetadata) GetPath() string {
	return m.Path
}
func (m *releaseMetadata) GetPostBuild() string {
	return m.PostBuild
}
func (m *releaseMetadata) GetPostDestroy() string {
	return m.PostDestroy
}
func (m *releaseMetadata) GetPreBuild() string {
	return m.PreBuild
}
func (m *releaseMetadata) GetPreDestroy() string {
	return m.PreDestroy
}
func (m *releaseMetadata) GetProvides() []string {
	return m.Provides
}
func (m *releaseMetadata) GetTest() string {
	return m.Test
}
func (m *releaseMetadata) GetType() string {
	return m.Type
}
func (m *releaseMetadata) GetVersion() string {
	return m.Version
}
func (m *releaseMetadata) GetDependencies() []string {
	return m.Depends
}
func (m *releaseMetadata) GetVariableContext() map[string]string {
    if m.VariableCtx == nil {
        return map[string]string{}
    }
	return m.VariableCtx
}
func (m *releaseMetadata) SetVariableInContext(v string, ref string) {
    ctx := m.GetVariableContext()
    ctx[v] = ref
    m.VariableCtx = ctx
}
func (m *releaseMetadata) GetReleaseId() string {
	return m.Type + "-" + m.Name + "-v" + m.Version
}

func (m *releaseMetadata) GetVersionlessReleaseId() string {
	return m.Type + "-" + m.Name
}

func (m *releaseMetadata) AddInputVariable(input *variable) {
	m.Inputs = append(m.Inputs, input)
}
func (m *releaseMetadata) AddOutputVariable(output *variable) {
	m.Outputs = append(m.Outputs, output)
}

func (m *releaseMetadata) ToJson() string {
	str, err := json.MarshalIndent(m, "", "   ")
	if err != nil {
		panic(err)
	}
	return string(str)
}

func (m *releaseMetadata) ToDict() (map[string]interface{}, error) {
    asJson := []byte(m.ToJson())
    result := map[string]interface{}{}
    if err := json.Unmarshal(asJson, &result); err != nil {
        return nil, err
    }
    return result, nil
}

func (m *releaseMetadata) WriteJsonFile(path string) error {
	contents := []byte(m.ToJson())
	return ioutil.WriteFile(path, contents, 0644)
}

func (m *releaseMetadata) AddFileWithDigest(path, hexDigest string) {
	m.Files[path] = hexDigest
}


func (m *releaseMetadata) GetDirectories() []string {
	dirs := map[string]bool{}
	for file := range m.Files {
		dir, _ := filepath.Split(file)
		dirs[dir] = true
		root := ""
		for _, d := range strings.Split(dir, "/") {
			if d != "" {
                root += d + "/"
                dirs[root] = true
			}
		}
	}
	result := []string{}
	for d := range dirs {
		if d != "" {
			result = append(result, d)
		}
	}
	return result
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
