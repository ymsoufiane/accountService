package response

var MessageErrInternal string
var MessageBadeRequest string
func init(){
//500 Internal Server Error
MessageErrInternal="500 Internal Server Error"
MessageBadeRequest="Bad Request !!"

}

type ErrorUniqueField struct{
	Name string
	Message string 
}