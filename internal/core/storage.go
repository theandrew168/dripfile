package core

type Storage struct {
	Account  AccountStorage
	Project  ProjectStorage
	Location LocationStorage
	Transfer TransferStorage
	Schedule ScheduleStorage
	History  HistoryStorage
}
