package mem

import (
	. "github.com/ankyra/escape-inventory/dao/types"
)

func (a *dao) GetApplications(project string) (map[string]*Application, error) {
	result := map[string]*Application{}
	for _, app := range a.namespaces[project] {
		result[app.App.Name] = app.App
	}
	return result, nil
}

func (a *dao) GetApplication(project, name string) (*Application, error) {
	prj, ok := a.namespaces[project]
	if !ok {
		return nil, NotFound
	}
	result, ok := prj[name]
	if !ok {
		return nil, NotFound
	}
	return result.App, nil
}

func (a *dao) FindAllVersions(app *Application) ([]string, error) {
	application := a.apps[app]
	versions := []string{}
	if application == nil {
		return versions, nil
	}
	for _, r := range application.Releases {
		versions = append(versions, r.Release.Version)
	}
	return versions, nil
}

func (a *dao) AddApplication(app *Application) error {
	apps, ok := a.namespaces[app.Project]
	if !ok {
		return NotFound
	}
	_, ok = apps[app.Name]
	if ok {
		return AlreadyExists
	}
	apps[app.Name] = &application{app, map[string]*release{}}
	a.apps[app] = apps[app.Name]
	a.applicationHooks[app] = NewHooks()
	return nil
}
func (a *dao) UpdateApplication(app *Application) error {
	apps, ok := a.namespaces[app.Project]
	if !ok {
		return NotFound
	}
	proj, ok := apps[app.Name]
	if !ok {
		return NotFound
	}
	apps[app.Name] = &application{app, proj.Releases}
	a.apps[app] = apps[app.Name]
	return nil
}

func (a *dao) GetApplicationHooks(app *Application) (Hooks, error) {
	apps, ok := a.namespaces[app.Project]
	if !ok {
		return nil, NotFound
	}
	unit, ok := apps[app.Name]
	if !ok {
		return nil, NotFound
	}
	return a.applicationHooks[unit.App], nil
}

func (a *dao) SetApplicationHooks(app *Application, hooks Hooks) error {
	apps, ok := a.namespaces[app.Project]
	if !ok {
		return NotFound
	}
	unit, ok := apps[app.Name]
	if !ok {
		return NotFound
	}
	a.applicationHooks[unit.App] = hooks
	return nil
}

func (a *dao) GetDownstreamHooks(app *Application) ([]*Hooks, error) {
	result := []*Hooks{}
	for downstream, subs := range a.subscriptions {
		for _, sub := range subs {
			if sub.Name == app.Name && sub.Project == app.Project {
				hooks := a.applicationHooks[downstream]
				result = append(result, &hooks)
			}
		}
	}
	return result, nil
}

func (a *dao) SetApplicationSubscribesToUpdatesFrom(app *Application, upstream []*Application) error {
	a.subscriptions[app] = upstream
	return nil
}
