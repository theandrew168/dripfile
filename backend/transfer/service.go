package transfer

type Service interface {
	Run(id string) error
}
