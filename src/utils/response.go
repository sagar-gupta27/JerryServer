//this faile handles responding weather it's streaming of bytes or sending normal response

package utils
import ("fmt"
"syscall")

type ResponseWriter struct {
	ConnFD   int  // ConnFD is the file descriptor (integer ID). It represents the TCP socket connection to a specific client.
	toStream bool // Indicates whether the response should be sent as a stream of chunks.
}

//creates new response write for every connection
func NewResponseWrite(ConnFD int) *ResponseWriter {
	return &ResponseWriter {ConnFD: ConnFD}
}
//function for sending html files
func (rw *ResponseWriter)SendHTML(statusCode int, body string){
	statusMessage := GetStatusCode(statusCode)
	response := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode,statusMessage)
	response += "content-type:  text/html; charset=utf-8\r\n"
	response += fmt.Sprintf("Content-Length: %d\r\n", len(body))
	response += "Connection: close\r\n"
	response += "\r\n"
	response += body
// divide  the response in chunks of bytes 
	responseBytes := []byte(response)
	//then send it back to client
	_, err := syscall.Write(rw.ConnFD, responseBytes)
	if err != nil {
		fmt.Printf("response Failed to send:= %v", err)
	}
}

//function for sending JSON response
func (rw *ResponseWriter)SendJson(statusCode int, body string){
	statusMessage := GetStatusCode(statusCode)
	response := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode,statusMessage)
	response += "content-type:  application/json\r\n"
	response += fmt.Sprintf("Content-Length: %d\r\n", len(body))
	response += "Connection: close\r\n"
	response += "\r\n"
	response += body

	// divide  the response in chunks of bytes 
	responseBytes := []byte(response)
	//then send it back to client
	_, err := syscall.Write(rw.ConnFD, responseBytes)
	if err != nil {
		fmt.Printf("response Failed to send:= %v", err)
	}
}
