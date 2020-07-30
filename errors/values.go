package errors

var (
	NotImplemented = New("not implemented")
	// TODO move to resolvers package!
	UnableToResolve = New("unable to resolve")
	EmptyArgs       = New("empty argument")
	NotFound        = New("not found")
	NotAuthorized   = New("not authorized")
	EmailUsed       = New("email already used")
)

func WrongType(expected, actual interface{}) error {
	return Errorf("wrong type: wanted %T, got %T", expected, actual)
}
