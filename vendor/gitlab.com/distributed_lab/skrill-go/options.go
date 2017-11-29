package skrill

type Option struct {
	Key   string
	Value string
}

func StatusURL(value string) Option {
	return Option{
		Key:   "status_url",
		Value: value,
	}
}

func Password(value string) Option {
	return Option{
		Key:   "password",
		Value: value,
	}
}

func Currency(value string) Option {
	return Option{
		Key:   "currency",
		Value: value,
	}
}
