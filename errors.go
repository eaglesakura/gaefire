package gaefire

type DatastoreError struct {
	message string
	errors  []error
}

func (it *DatastoreError) Error() string {
	return it.message
}
