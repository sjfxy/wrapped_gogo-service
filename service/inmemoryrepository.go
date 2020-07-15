package service

import (
	"errors"
	"github.com/cloudnativego/gogo-engine"
	"strings"
)

type inMemoryRepository struct {
	matches []gogo.Match
}

// NewRepository creates a new in-momory match repository
func NewInMemoryRepository() *inMemoryRepository {
	repo := &inMemoryRepository{}
	repo.matches = []gogo.Match{}
	return repo
}
func (repo *inMemoryRepository) addMatch(match gogo.Match) (err error) {
	repo.matches = append(repo.matches, match)
	return err
}
func (repo *inMemoryRepository) getMatches() (matches []gogo.Match, err error) {
	matches = repo.matches
	return
}
func (repo *inMemoryRepository) getMatch(id string) (match gogo.Match, err error) {
	found := false
	for _, target := range repo.matches {
		if strings.Compare(target.ID, id) == 0 {
			match = target
			found = true
		}
	}
	if !found {
		err = errors.New("Cloud not found match in repository")
	}
	return match, err
}
func (repo *inMemoryRepository) updateMatch(id string, match gogo.Match) (err error) {
	found := false
	for k, v := range repo.matches {
		if strings.Compare(v.ID, id) == 0 {
			repo.matches[k] = match
			found = true
		}
	}
	if !found {
		err = errors.New("Cloud not found match in repository")
	}
	return
}

// add update delte getklist
// append
// for tra-,ayr
// forrpo[;]= match
// reyirn
//
