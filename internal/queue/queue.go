package queue

type Publisher interface {
}

type Subscriber interface {
}

func NewPostgresPublisher() Publisher {
	return nil
}

func NewPostgresSubscriber() Subscriber {
	return nil
}
