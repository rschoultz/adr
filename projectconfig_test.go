package main

import (
	"testing"
)

func TestAppendProject(t *testing.T) {
	adrConfig := AdrConfig{nil}
	adrConfig, _ = addProject("hello", adrConfig)
	adrConfig, _ = addProject("is/it/me/you/are/looking/for", adrConfig)

	count := len(adrConfig.Projects)
	if count != 2 {
		t.Errorf("Projects not 2, rather %d", count)
	}
}

func TestAppendProjectOnlyOnce(t *testing.T) {
	adrConfig := AdrConfig{nil}
	adrConfig, _ = addProject("hello", adrConfig)
	adrConfig, _ = addProject("hello", adrConfig)

	count := len(adrConfig.Projects)
	if count != 1 {
		t.Errorf("Projects not 1, rather %d", count)
	}
}

func TestAppendProjectOnSubPathShouldNotWork(t *testing.T) {
	adrConfig := AdrConfig{nil}
	adrConfig, _ = addProject("hello", adrConfig)
	adrConfig, _ = addProject("hello/is/it/me/you/are/looking/for", adrConfig)

	count := len(adrConfig.Projects)
	if count != 1 {
		t.Errorf("Projects not 1, rather %d", count)
	}
}

func TestAppendProjectOnSuperPathShouldNotWork(t *testing.T) {
	adrConfig := AdrConfig{nil}
	adrConfig, _ = addProject("hello/is/it/me/you/are/looking/for", adrConfig)
	adrConfig, _ = addProject("hello", adrConfig)

	count := len(adrConfig.Projects)
	if count != 1 {
		t.Errorf("Projects not 1, rather %d", count)
	}
}
