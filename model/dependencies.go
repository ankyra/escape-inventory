package model

import (
	"github.com/ankyra/escape-inventory/dao"
	"github.com/ankyra/escape-inventory/dao/types"
)

func GetDownstreamDependencies(namespace, name, version string) ([]*types.Dependency, error) {
	releaseId := name + "-" + version
	release, err := ResolveReleaseId(namespace, releaseId)
	if err != nil {
		return nil, err
	}
	return dao.GetDownstreamDependencies(release)
}

func GetDownstreamDependenciesByGroups(namespace, name, version string, readGroups []string) ([]*types.Dependency, error) {
	releaseId := name + "-" + version
	release, err := ResolveReleaseId(namespace, releaseId)
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

func GetDependencyGraph(namespace, name, version string, downstreamFunc DownstreamDependenciesResolver) (*DependencyGraph, error) {
	result := &DependencyGraph{
		Nodes: []*DependencyGraphNode{},
		Edges: []*DependencyGraphEdge{},
	}
	releaseId := name + "-" + version
	release, err := ResolveReleaseId(namespace, releaseId)
	if err != nil {
		return nil, err
	}
	mainId := release.Metadata.GetQualifiedReleaseId()
	result.AddNode(mainId, "main")
	mainId = "main" + mainId

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

	for _, c := range release.Metadata.Consumes {
		result.AddNode(c.Name, "consumer")
		result.AddEdge("consumer"+c.Name, mainId, "consumes")
	}
	for _, c := range release.Metadata.Provides {
		result.AddNode(c.Name, "provider")
		result.AddEdge(mainId, "provider"+c.Name, "provides")
	}
	for _, dep := range upstream {
		id := dep.Project + "/" + dep.Application + "-v" + dep.Version
		typ := "upstream"
		label := "upstream"
		if dep.IsExtension {
			typ = "extension"
			label = "extends"
		}
		result.AddNode(id, typ)
		result.AddEdge(mainId, typ+id, label)
	}
	for _, dep := range downstream {
		id := dep.Project + "/" + dep.Application + "-v" + dep.Version
		typ := "downstream"
		label := "downstream"
		if dep.IsExtension {
			typ = "extension"
			label = "extends"
		}
		result.AddNode(id, typ)
		result.AddEdge(typ+id, mainId, label)
	}
	return result, nil
}

func GetDependencyGraphByGroups(namespace, name, version string, groups []string) (*DependencyGraph, error) {
	resolver := func(release *types.Release) ([]*types.Dependency, error) {
		return dao.GetDownstreamDependenciesByGroups(release, groups)
	}
	return GetDependencyGraph(namespace, name, version, resolver)
}
