package model

import (
	"github.com/ankyra/escape-registry/dao"
	"github.com/ankyra/escape-registry/dao/types"
)

func GetDownstreamDependencies(project, name, version string) ([]*types.Dependency, error) {
	releaseId := name + "-" + version
	release, err := ResolveReleaseId(project, releaseId)
	if err != nil {
		return nil, err
	}
	return dao.GetDownstreamDependencies(release)
}

func GetDownstreamDependenciesByGroups(project, name, version string, readGroups []string) ([]*types.Dependency, error) {
	releaseId := name + "-" + version
	release, err := ResolveReleaseId(project, releaseId)
	if err != nil {
		return nil, err
	}
	return dao.GetDownstreamDependenciesByGroups(release, readGroups)
}

type DependencyGraphNode struct {
	Id    string `json:"id"`
	Label string `json:"label"`
	Type  string `json:"type"`
}
type DependencyGraphEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"`
}

type DependencyGraph struct {
	Nodes []*DependencyGraphNode `json:"nodes"`
	Edges []*DependencyGraphEdge `json:"edges"`
}

func (d *DependencyGraph) AddNode(id, typ string) {
	d.Nodes = append(d.Nodes, &DependencyGraphNode{
		Id:    typ + id,
		Label: id,
		Type:  typ,
	})
}
func (d *DependencyGraph) AddEdge(from, to, typ string) {
	d.Edges = append(d.Edges, &DependencyGraphEdge{
		From: from,
		To:   to,
		Type: typ,
	})
}

type DownstreamDependenciesResolver func(*types.Release) ([]*types.Dependency, error)

func GetDependencyGraph(project, name, version string, downstreamFunc DownstreamDependenciesResolver) (*DependencyGraph, error) {
	result := &DependencyGraph{
		Nodes: []*DependencyGraphNode{},
		Edges: []*DependencyGraphEdge{},
	}
	releaseId := name + "-" + version
	release, err := ResolveReleaseId(project, releaseId)
	if err != nil {
		return nil, err
	}
	mainId := release.Metadata.GetQualifiedReleaseId()
	result.AddNode(mainId, "main")

	upstream, err := dao.GetDependencies(release)
	if err != nil {
		return nil, err
	}

	if downstreamFunc == nil {
		downstreamFunc = dao.GetDownstreamDependencies
	}
	downstream, err := downstreamFunc(release)
	if err != nil {
		return nil, err
	}

	for _, ext := range release.Metadata.Extends {
		result.AddNode(ext.ReleaseId, "release")
		result.AddEdge(mainId, ext.ReleaseId, "extends")
	}
	for _, dep := range upstream {
		id := dep.Project + "/" + dep.Application + "-v" + dep.Version
		result.AddNode(id, "upstream")
		result.AddEdge(mainId, id, "upstream")
	}
	for _, dep := range downstream {
		id := dep.Project + "/" + dep.Application + "-v" + dep.Version
		result.AddNode(id, "downstream")
		result.AddEdge(id, mainId, "downstream")
	}
	return result, nil
}
