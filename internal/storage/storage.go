package storage

import (
	"github.com/theandrew168/dripfile/internal/core"
)

// aggregation of core storage interfaces
type Storage struct {
	Project  Project
	Account  Account
	Session  Session
	Location Location
	Transfer Transfer
	Schedule Schedule
	Job      Job
	History  History
}

type Project interface {
	Create(project *core.Project) error
	Read(id string) (core.Project, error)
	Delete(project core.Project) error
}

type Account interface {
	Create(account *core.Account) error
	Read(id string) (core.Account, error)
	Update(account core.Account) error
	Delete(account core.Account) error

	ReadByEmail(email string) (core.Account, error)
	CountByProject(project core.Project) (int, error)
}

type Session interface {
	Create(session *core.Session) error
	Read(hash string) (core.Session, error)
	Delete(session core.Session) error
}

type Location interface {
	Create(location *core.Location) error
	Read(id string) (core.Location, error)
	Update(location core.Location) error
	Delete(location core.Location) error

	ReadManyByProject(project core.Project) ([]core.Location, error)
}

type Transfer interface {
	Create(transfer *core.Transfer) error
	Read(id string) (core.Transfer, error)
	Update(transfer core.Transfer) error
	Delete(transfer core.Transfer) error

	ReadManyByProject(project core.Project) ([]core.Transfer, error)
}

type Schedule interface {
	Create(schedule *core.Schedule) error
	Read(id string) (core.Schedule, error)
	Update(schedule core.Schedule) error
	Delete(schedule core.Schedule) error
}

type Job interface {
	// baseline CRUD ops all deal with one record
	Create(job *core.Job) error
	Read(id string) (core.Job, error)
	Update(job core.Job) error
	Delete(job core.Job) error
}

type History interface {
	Create(history *core.History) error
	Read(id string) (core.History, error)
	Update(history core.History) error
	Delete(history core.History) error
}
