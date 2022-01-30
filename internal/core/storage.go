package core

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
