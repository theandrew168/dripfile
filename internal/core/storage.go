package core

type Storage struct {
	Account  AccountStorage
	Session  SessionStorage
	Project  ProjectStorage
	Location LocationStorage
	Transfer TransferStorage
	Schedule ScheduleStorage
	History  HistoryStorage
}
