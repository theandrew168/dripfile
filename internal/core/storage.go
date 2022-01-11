package core

type Storage interface {
	AccountStorage
	ProjectStorage
	LocationStorage
	TransferStorage
	ScheduleStorage
	HistoryStorage
}
