package config

//Response is the response template of the api
type Response struct {
	//Message for the response incase of sucess
	Message *string
	//Error messaage incase of failures/errors
	Error *string
	//Data is the response payload of the api
	Data interface{}
}

//Error to create an error response
func Error(message string) Response {
	return Response{Error: &message}
}

//Success to create a success response
func Success(data interface{}, message string) Response {
	return Response{Message: &message, Data: data}
}
