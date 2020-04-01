package main

import (
	"testing"
)

func TestAppendProject(t *testing.T) {
	adrConfig := AdrConfig{nil}
	adrConfig, _ = addProject(adrConfig, "hello")
	adrConfig, _ = addProject(adrConfig, "is/it/me/you/are/looking/for")

	count := len(adrConfig.Projects)
	if count != 2 {
		t.Errorf("Projects not 2, rather %d", count)
	}
}

func TestAppendProjectOnlyOnce(t *testing.T) {
	adrConfig := AdrConfig{nil}
	adrConfig, _ = addProject(adrConfig, "hello")
	adrConfig, _ = addProject(adrConfig, "hello")

	count := len(adrConfig.Projects)
	if count != 1 {
		t.Errorf("Projects not 1, rather %d", count)
	}
}

func TestAppendProjectOnSubPathShouldNotWork(t *testing.T) {
	adrConfig := AdrConfig{nil}
	adrConfig, _ = addProject(adrConfig, "hello")
	adrConfig, _ = addProject(adrConfig, "hello/is/it/me/you/are/looking/for")

	count := len(adrConfig.Projects)
	if count != 1 {
		t.Errorf("Projects not 1, rather %d", count)
	}
}

func TestAppendProjectOnSuperPathShouldNotWork(t *testing.T) {
	adrConfig := AdrConfig{nil}
	adrConfig, _ = addProject(adrConfig, "hello/is/it/me/you/are/looking/for")
	adrConfig, _ = addProject(adrConfig, "hello")

	count := len(adrConfig.Projects)
	if count != 1 {
		t.Errorf("Projects not 1, rather %d", count)
	}
}
