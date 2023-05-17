package memory

type Info struct{}

func (info Info) Validate() error {
	return nil
}

type FileServer struct {
	info Info
	data map[string][]byte
}

func New(info Info) (*FileServer, error) {
	fs := FileServer{
		info: info,
		data: make(map[string][]byte),
	}

	return &fs, nil
}
