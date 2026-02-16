package spec

// ServerTypeFlag is an enum and can be one of the following values
type ServerTypeFlag string

const (
	REST    ServerTypeFlag = "rest"
	GRPC    ServerTypeFlag = "grpc"
	GraphQL ServerTypeFlag = "graphql"
)

type Config struct {
	OutputPath string
	Default    bool
	Form       Form
}

type Form struct {
	Name           string
	ServerTypeFlag ServerTypeFlag
	DatabaseFlag   bool
	MakefileFlag   bool
}
