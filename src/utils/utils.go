//this file contains the COntentType of request and status code

package utils
import ("fmt")

const (
	ContentTypeHTML       = "text/html; charset=utf-8"
	ContentTypeJSON       = "application/json"
	ContentTypePlain      = "text/plain"
	ContentTypeJavaScript = "application/javascript"
	ContentTypeCSS        = "text/css"
)


const (
	StatusOK                  = 200
	StatusCreated             = 201
	StatusNoContent           = 204
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusTooManyRequests     = 429
	StatusInternalServerError = 500
)

func GetStatusCode(statusCode int)string{
statusCodes := make(map[int]string)
statusCodes[200] = "StatusOK"
statusCodes[201] = "StatusCreated"
statusCodes[204] = "StatusNoContent"
statusCodes[400] = "StatusBadRequest"
statusCodes[401] = "StatusUnauthorized"
statusCodes[403] = "StatusForbidden"
statusCodes[404] = "StatusNotFound"
statusCodes[429] = "StatusTooManyRequests"
statusCodes[500] = "StatusInternalServerError"
fmt.Println(statusCodes[200])
status, isExist := statusCodes[statusCode]
if isExist {
	return status
}
return "UNKNOWN"
}
// parseRange parses HTTP Range header
// Example: "bytes=0-1023" → start=0, end=1024
func  ParseRange(rangeHeader string , fileSize int) (int,int){
	var start, end int

	fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)

	if end == 0 || end >= fileSize {
		end = fileSize
	} else{
		end++
	}
	return start, end
}