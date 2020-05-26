package inmemory

import (
	"fmt"

	. "github.com/fgrehm/brinfo/core"
)

type contentSourceRepo struct {
	sourcesByID   map[string]*ContentSource
	sourcesByHost map[string]*ContentSource
}

func NewContentSourceRepo() ContentSourceRepo {
	return &contentSourceRepo{
		sourcesByID:   map[string]*ContentSource{},
		sourcesByHost: map[string]*ContentSource{},
	}
}

func (r *contentSourceRepo) Register(cs *ContentSource) error {
	if cs.ID == "" {
		return fmt.Errorf("No ID provided: %+v", cs)
	}
	if cs.Host == "" {
		return fmt.Errorf("No Host provided: %+v", cs)
	}

	if _, ok := r.sourcesByID[cs.ID]; ok {
		return fmt.Errorf("Duplicated ID provided: %+v", cs)
	}
	if _, ok := r.sourcesByHost[cs.Host]; ok {
		return fmt.Errorf("Duplicated Host provided: %+v", cs)
	}

	r.sourcesByID[cs.ID] = cs
	r.sourcesByHost[cs.Host] = cs

	return nil
}

func (r *contentSourceRepo) FindByID(id string) (*ContentSource, error) {
	cs, ok := r.sourcesByID[id]
	if !ok {
		return nil, fmt.Errorf("Content source not found: %s", id)
	}

	return cs, nil
}

func (r *contentSourceRepo) FindByHost(host string) (*ContentSource, error) {
	cs, ok := r.sourcesByHost[host]
	if !ok {
		return nil, fmt.Errorf("Content source not found: %s", host)
	}

	return cs, nil
}

func (r *contentSourceRepo) GetByHost(host string) (*ContentSource, error) {
	cs, ok := r.sourcesByHost[host]
	if !ok {
		return nil, nil
	}
	return cs, nil
}
