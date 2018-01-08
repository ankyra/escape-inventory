package handlers

import (
	"github.com/ankyra/escape-inventory/dao/types"
	"github.com/gorilla/mux"
	. "gopkg.in/check.v1"
)

const (
	FeedURL            = "/api/v1/project/__feed"
	ProjectFeedURL     = "/api/v1/project/{project}/feed"
	projectFeedTestURL = "/api/v1/project/project/feed"
)

/*
	FeedHandler
*/

func (s *suite) feedMuxWithProvider(provider *feedHandlerProvider) *mux.Router {
	return s.GetMuxForHandler("GET", FeedURL, provider.FeedHandler)
}

func (s *suite) Test_FeedHandler_happy_path(c *C) {
	var capturedPage int
	feedPayload := []*types.FeedEvent{
		types.NewCreateProjectEvent("project"),
		types.NewReleaseEvent("project", "name", "1.0", "test-user"),
	}
	provider := &feedHandlerProvider{
		GetFeedPage: func(page int) ([]*types.FeedEvent, error) {
			capturedPage = page
			return feedPayload, nil
		},
	}
	resp := s.testGET(c, s.feedMuxWithProvider(provider), FeedURL)
	s.ExpectSuccessResponse_with_JSON(c, resp, feedPayload)
	c.Assert(capturedPage, Equals, 0)
}

func (s *suite) Test_FeedHandler_happy_path_with_page(c *C) {
	var capturedPage int
	feedPayload := []*types.FeedEvent{
		types.NewCreateProjectEvent("project"),
		types.NewReleaseEvent("project", "name", "1.0", "test-user"),
	}
	provider := &feedHandlerProvider{
		GetFeedPage: func(page int) ([]*types.FeedEvent, error) {
			capturedPage = page
			return feedPayload, nil
		},
	}
	resp := s.testGET(c, s.feedMuxWithProvider(provider), FeedURL+"?page=100")
	s.ExpectSuccessResponse_with_JSON(c, resp, feedPayload)
	c.Assert(capturedPage, Equals, 100)
}

func (s *suite) Test_FeedHandler_happy_path_invalid_page(c *C) {
	var capturedPage int
	feedPayload := []*types.FeedEvent{
		types.NewCreateProjectEvent("project"),
		types.NewReleaseEvent("project", "name", "1.0", "test-user"),
	}
	provider := &feedHandlerProvider{
		GetFeedPage: func(page int) ([]*types.FeedEvent, error) {
			capturedPage = page
			return feedPayload, nil
		},
	}
	resp := s.testGET(c, s.feedMuxWithProvider(provider), FeedURL+"?page=asbasd")
	s.ExpectSuccessResponse_with_JSON(c, resp, feedPayload)
	c.Assert(capturedPage, Equals, 0)
}

func (s *suite) Test_FeedHandler_fails_if_GetFeedPage_fails(c *C) {
	provider := &feedHandlerProvider{
		GetFeedPage: func(page int) ([]*types.FeedEvent, error) {
			return nil, types.Unauthorized
		},
	}
	resp := s.testGET(c, s.feedMuxWithProvider(provider), FeedURL)
	s.ExpectErrorResponse(c, resp, 401, "")
}

/*
	ProjectFeedHandler
*/

func (s *suite) projectFeedMuxWithProvider(provider *feedHandlerProvider) *mux.Router {
	return s.GetMuxForHandler("GET", ProjectFeedURL, provider.ProjectFeedHandler)
}

func (s *suite) Test_ProjectFeedHandler_happy_path(c *C) {
	var capturedPage int
	var capturedProject string
	feedPayload := []*types.FeedEvent{
		types.NewCreateProjectEvent("project"),
		types.NewReleaseEvent("project", "name", "1.0", "test-user"),
	}
	provider := &feedHandlerProvider{
		GetProjectFeed: func(project string, page int) ([]*types.FeedEvent, error) {
			capturedProject = project
			capturedPage = page
			return feedPayload, nil
		},
	}
	resp := s.testGET(c, s.projectFeedMuxWithProvider(provider), projectFeedTestURL)
	s.ExpectSuccessResponse_with_JSON(c, resp, feedPayload)
	c.Assert(capturedPage, Equals, 0)
	c.Assert(capturedProject, Equals, "project")
}

func (s *suite) Test_ProjectFeedHandler_happy_path_with_page(c *C) {
	var capturedPage int
	var capturedProject string
	feedPayload := []*types.FeedEvent{
		types.NewCreateProjectEvent("project"),
		types.NewReleaseEvent("project", "name", "1.0", "test-user"),
	}
	provider := &feedHandlerProvider{
		GetProjectFeed: func(project string, page int) ([]*types.FeedEvent, error) {
			capturedProject = project
			capturedPage = page
			return feedPayload, nil
		},
	}
	resp := s.testGET(c, s.projectFeedMuxWithProvider(provider), projectFeedTestURL+"?page=100")
	s.ExpectSuccessResponse_with_JSON(c, resp, feedPayload)
	c.Assert(capturedPage, Equals, 100)
	c.Assert(capturedProject, Equals, "project")
}

func (s *suite) Test_ProjectFeedHandler_fails_if_GetProjectFeed_fails(c *C) {
	provider := &feedHandlerProvider{
		GetProjectFeed: func(project string, page int) ([]*types.FeedEvent, error) {
			return nil, types.Unauthorized
		},
	}
	resp := s.testGET(c, s.projectFeedMuxWithProvider(provider), projectFeedTestURL+"?page=100")
	s.ExpectErrorResponse(c, resp, 401, "")
}
