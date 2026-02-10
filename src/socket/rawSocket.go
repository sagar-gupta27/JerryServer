//this file handles creation of socket connection from scratch by call OS via syscall
//there is list of every available syscall for every OS structure on ``https://cs.opensource.google/go/go/+/refs/tags/go1.25.6:src/syscall/zerrors_linux_amd64.go;l=9``
package socket

import ("fmt"
	  "syscall"
	  "httpServer/src/utils"
	  "strings"
	)
// Request represents an HTTP request
type Request struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    string
	Raw     string
}

type Response struct {
	ConnFD int 
	writer *utils.ResponseWriter
}
// this is polymorphism for route handler function
type RouteHandlerFunc func(*Request, *Response)

type Server struct {
	socketFD int
	ipAddr string
	routes map[string]RouteHandlerFunc
}

func NewServerInstance()*Server{
	return &Server{
		routes: make(map[string]RouteHandlerFunc)}
}


// defining the StartServer method in Servre clasee
func (s *Server) StartServer(PortNo int){
	// raw socket creation, this return socketFD. FD(file descriptor) =integer given by OS, it's like a ticket given by os, so OS recognize this Socket from his FD 
	//AF_INET this refers to ipv4 domain
	//AF_INET6 refers to ipv6 domain
	//SOCK_STREAM = tcp connection type
	// SOCK_DGRAM = udp connection type, and 0 means use default protocal in this case the default protocal is TCP
	socketFD, err := syscall.Socket(syscall.AF_INET,syscall.SOCK_STREAM,0)
	s.socketFD  =socketFD
	//check for error
	if err != nil{
		panic(fmt.Sprintf("socket creation failed:= %v", err))
	}
	//add defer to close the socket when the programme exits
	defer syscall.Close(s.socketFD)
	// Set socket options.
	// Enable SO_REUSEADDR to allow the program to reuse the same port
	// every time it starts. Setting 1 means true for this option.
	err = syscall.SetsockoptInt(s.socketFD, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	//error handling
	if err != nil {
		panic(fmt.Sprintf("SetsockoptInt failed:= %v", err))
	}


	//bind socket to address and port
	address := syscall.SockaddrInet4{
		Port: PortNo,
		Addr: [4]byte{0, 0, 0, 0} }

	err = syscall.Bind(s.socketFD, &address)
	if err != nil {
		panic(fmt.Sprintf("Bind failed: %v", err))
	}
	
	// start listening for calls 
	//the scond argument for max pending request, pending request more than 10 will be rejected
	err = syscall.Listen(s.socketFD, 10)
	if err != nil {
		panic(fmt.Sprintf("Listen failed: %v", err))
	}

	// server  started to listen requests
	fmt.Printf("server is started successfully!\n") 
	fmt.Printf("Visit: http://localhost:%d\n", PortNo)


	//this inifinite loop for handleConnection forever
	for {
		clientFD, clientAddr, err := syscall.Accept(s.socketFD)

		if err != nil {
			fmt.Printf("request Accept failed: %v\n", err)
			continue
		}
		//handle this funtion in go routine for concurrency
		go s.handleConnection(clientFD, clientAddr)
	}
}

func (s *Server)  HandleRouteFunc(path string, function RouteHandlerFunc){
	s.routes[path] = function
}


func(s *Server) handleConnection(clientFD int, clientAddr syscall.Sockaddr){
	//close the connection when function is about to exit 
	defer syscall.Close(clientFD)

	// Extract client IP
	// addrArray := clientAddr.(*syscall.SockaddrInet4)
	// clientIP := fmt.Sprintf("%d.%d.%d.%d",addrArray.Addr[0],  addrArray.Addr[1],addrArray.Addr[2], addrArray.Addr[3]) 

	// reatelimit 


	// create buffer array for storing requests buffre
	// reainf request 8kb per sseinc
	buffer := make([]byte, 8192) 
	number, err := syscall.Read(clientFD, buffer)
	if err != nil {
		fmt.Printf("Reading fail := %v\n", err)
	}

	rawMessage := string(buffer[:number]) 

	//parsing http request
	req := s.ParseRequest(rawMessage)
fmt.Println(req.Path)
	res := &Response{
ConnFD : clientFD,
writer : utils.NewResponseWrite(clientFD)}
	routeHandler, exist := s.routes[req.Path]
	if !exist {
		res.writer.SendHTML(404, "<h1>403 - Not Found</h1>")
		return
	}
	routeHandler(req, res)

}

func(s *Server) ParseRequest(data string) *Request {

	lines := strings.Split(data, "\r\n")
	if len(lines) == 0 {
		return &Request{Raw: data}
	}
	parts := strings.Split(lines[0], " ")
	req := &Request{
		Method:  "GET",
		Path:    "/",
		Headers: make(map[string]string),
		Raw: data,
	}

	if len(parts) >= 2 {
		req.Method = parts[0]
		req.Path = parts[1]
	}

	// Parse headers
	for i := 1; i < len(lines); i++ {
		if lines[i] == "" {
			break // End of headers
		}
		parts := strings.SplitN(lines[i], ": ", 2)
		if len(parts) == 2 {
			req.Headers[parts[0]] = parts[1]
		}
	}

	return req
}

func (s *Server) ShutDown() {
	if s.socketFD > 0 {
		syscall.Close(s.socketFD)
	}
}

// response method for sending HTML
func (r *Response)  SendHTML(statusCode int, data string){
	r.writer.SendHTML(statusCode, data)
}


// response method for sendJSON 
func (r *Response)  SendJson(statusCode int, data string){
	r.writer.SendJson(statusCode, data)
}


func (r *Response)  InitChunkStream(statusCode int, contentType string){
	r.writer.InitChunkStream(statusCode, contentType)
}

func (r *Response) SendChunk(data []byte){
	r.writer.SendChunk(data)
}
func (r *Response) StopSendindChunks(){
	r.writer.StopSendindChunks()
}