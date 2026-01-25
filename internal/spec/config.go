package spec

type Config struct {
	OutputPath string
	Form       Form
}

type Form struct {
	SomeFlag    string
	AnotherFlag string
}
