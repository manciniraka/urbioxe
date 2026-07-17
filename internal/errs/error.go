package errs

type AppError struct {
	Code    int
	Message string
	Err     error
}
