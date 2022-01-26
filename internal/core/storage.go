package core

type Storage struct {
	Account  AccountStorage
	Session  SessionStorage
	Project  ProjectStorage
	Member   MemberStorage
	Location LocationStorage
	Transfer TransferStorage
	Schedule ScheduleStorage
	Job      JobStorage
	History  HistoryStorage
}
