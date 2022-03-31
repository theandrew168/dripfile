package database

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/pkg/core"
	"github.com/theandrew168/dripfile/pkg/database/postgres"
)

// aggregation of core storage interfaces
type Storage struct {
	Project  ProjectStorage
	Account  AccountStorage
	Session  SessionStorage
	Location LocationStorage
	Transfer TransferStorage
	Schedule ScheduleStorage
	Job      JobStorage
	History  HistoryStorage
}

func NewPostgresStorage(pool *pgxpool.Pool) Storage {
	s := Storage{
		Project:  postgres.NewProjectStorage(pool),
		Account:  postgres.NewAccountStorage(pool),
		Session:  postgres.NewSessionStorage(pool),
		Location: postgres.NewLocationStorage(pool),
		Transfer: postgres.NewTransferStorage(pool),
		Schedule: postgres.NewScheduleStorage(pool),
		Job:      postgres.NewJobStorage(pool),
		History:  postgres.NewHistoryStorage(pool),
	}
	return s
}

type ProjectStorage interface {
	Create(project *core.Project) error
	Read(id string) (core.Project, error)
	Delete(project core.Project) error
}

type AccountStorage interface {
	Create(account *core.Account) error
	Read(id string) (core.Account, error)
	Update(account core.Account) error
	Delete(account core.Account) error

	ReadByEmail(email string) (core.Account, error)
	CountByProject(project core.Project) (int, error)
}

type SessionStorage interface {
	Create(session *core.Session) error
	Read(hash string) (core.Session, error)
	Delete(session core.Session) error
	DeleteExpired() error
}

type LocationStorage interface {
	Create(location *core.Location) error
	Read(id string) (core.Location, error)
	Update(location core.Location) error
	Delete(location core.Location) error

	ReadManyByProject(project core.Project) ([]core.Location, error)
}

type TransferStorage interface {
	Create(transfer *core.Transfer) error
	Read(id string) (core.Transfer, error)
	Update(transfer core.Transfer) error
	Delete(transfer core.Transfer) error

	ReadManyByProject(project core.Project) ([]core.Transfer, error)
}

type ScheduleStorage interface {
	Create(schedule *core.Schedule) error
	Read(id string) (core.Schedule, error)
	Update(schedule core.Schedule) error
	Delete(schedule core.Schedule) error
}

type JobStorage interface {
	// baseline CRUD ops all deal with one record
	Create(job *core.Job) error
	Read(id string) (core.Job, error)
	Update(job core.Job) error
	Delete(job core.Job) error
}

type HistoryStorage interface {
	Create(history *core.History) error
	Read(id string) (core.History, error)

	ReadManyByProject(project core.Project) ([]core.History, error)
}
