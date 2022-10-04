package model

import (
	"crypto/sha256"
	"fmt"

	"github.com/keboola/go-utils/pkg/deepcopy"
)

type TemplateRepositoryType string

const (
	RepositoryTypeDir = `dir`
	RepositoryTypeGit = `git`
)

type TemplateRepositories struct {
	asSlice []TemplateRepository
	asMap   map[string]TemplateRepository
}

type TemplateRepository struct {
	Type TemplateRepositoryType `json:"type" validate:"oneof=dir git"`
	Name string                 `json:"name" validate:"required,max=40"`
	Url  string                 `json:"url" validate:"required"`
	Ref  string                 `json:"ref,omitempty" validate:"required_if=Type git"`
}

// String returns human-readable name of the repository.
func (r TemplateRepository) String() string {
	if r.Type == RepositoryTypeDir {
		return fmt.Sprintf("dir:%s", r.Url)
	}
	return fmt.Sprintf("%s:%s", r.Url, r.Ref)
}

// Hash returns unique identifier of the repository.
func (r TemplateRepository) Hash() string {
	hash := fmt.Sprintf("%s:%s:%s", r.Type, r.Url, r.Ref)
	sha := sha256.Sum256([]byte(hash))
	return string(sha[:])
}

type TemplateRef interface {
	Repository() TemplateRepository
	WithRepository(TemplateRepository) TemplateRef
	TemplateId() string
	Version() string
	Name() string
	FullName() string
}

type templateRef struct {
	repository TemplateRepository
	templateId string // for example "my-template"
	version    string // for example "v1"
}

func NewTemplateRepositories() *TemplateRepositories {
	return &TemplateRepositories{asMap: make(map[string]TemplateRepository)}
}

func (v *TemplateRepositories) Add(repo TemplateRepository) {
	v.asSlice = append(v.asSlice, repo)
	v.asMap[repo.Name] = repo
}

func (v *TemplateRepositories) Get(name string) (TemplateRepository, bool) {
	repo, found := v.asMap[name]
	return repo, found
}

func (v *TemplateRepositories) All() []TemplateRepository {
	return deepcopy.Copy(v.asSlice).([]TemplateRepository)
}

func NewTemplateRef(repository TemplateRepository, templateId string, version string) TemplateRef {
	return templateRef{
		repository: repository,
		templateId: templateId,
		version:    version,
	}
}

func (r templateRef) Repository() TemplateRepository {
	return r.repository
}

func (r templateRef) WithRepository(repository TemplateRepository) TemplateRef {
	r.repository = repository
	return r
}

func (r templateRef) TemplateId() string {
	return r.templateId
}

func (r templateRef) Version() string {
	return r.version
}

// Name without repository, for example "my-template/v1.
func (r templateRef) Name() string {
	return fmt.Sprintf("%s/%s", r.templateId, r.version)
}

// FullName with repository, for example "keboola/my-template/v1.
func (r templateRef) FullName() string {
	return fmt.Sprintf("%s/%s/%s", r.repository.Name, r.templateId, r.version)
}
