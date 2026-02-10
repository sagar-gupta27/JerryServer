//this file handles streaming of chunks of memory over the network
package utils

import ("fmt"
"syscall")

//this is for starting the stream 
func (rw *ResponseWriter) InitChunkStream(statusCode int, contentType string){
	statusMessage := GetStatusCode(statusCode)
	response := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusMessage)
	response += fmt.Sprintf("Content-Type: %s\r\n", contentType)
	response += "Transfer-Encoding: chunked\r\n" // KEY: No Content-Length!
	response += "Connection: keep-alive\r\n"
	response += "\r\n"

	// divide  the response in chunks of bytes 
	responseBytes := []byte(response)
	//then send it back to client
	_, err := syscall.Write(rw.ConnFD, responseBytes)
	if err != nil {
		fmt.Printf("response Failed to send:= %v", err)
	}
	rw.toStream = true
}

//this this is for sending one chunks of data
func (rw *ResponseWriter) SendChunk(data []byte){
	if !rw.toStream {
		fmt.Println("stream is not active for this")
		return
	}
	//length in bytes
	size := len(data)
	chunk := fmt.Sprintf("%X\r\n", size)
	chunk += string(data) + "\r\n"

	//// divide  the response in chunks of bytes 
	responseBytes := []byte(chunk)
	_, err := syscall.Write(rw.ConnFD, responseBytes)
	if err != nil {
		fmt.Printf("response Failed to send:= %v", err)
	}
}
//this this is for sending one chunks of data
func (rw *ResponseWriter) StopSendindChunks(){
	if !rw.toStream {
		fmt.Println("stream is not active for this")
		return
	}
	//send zero-length chunk to end
	finalChunk := "0\r\n\r\n"
	_, err := syscall.Write(rw.ConnFD, []byte(finalChunk))
	if err != nil {
		fmt.Printf("response Failed to send:= %v", err)
	}
	rw.toStream  =false
}



