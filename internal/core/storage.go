package core

type Storage interface {
	AccountStorage
	ProjectStorage
	LocationStorage
	JobStorage
	ScheduleStorage
	HistoryStorage
}
