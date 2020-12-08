package indicators

type IndicatorError struct {
	msg string
	Err error
}

func (e IndicatorError) Error() string { return e.msg }

//IndicatorNotReadyError is an error thrown when an indicator needs more data to be used
type IndicatorNotReadyError struct {
	Msg string //description of error
	Len int    //the length needed before trying again
}

func (e IndicatorNotReadyError) Error() string { return e.Msg }
