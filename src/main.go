package main

import ("fmt"
"httpServer/src/socket"
"os/signal"
"os"
"httpServer/src/utils"
"httpServer/src/download"
"syscall"
)
const PORT_NO = 8080
func main(){
	srv := socket.NewServerInstance()
	srv.HandleRouteFunc("/", Home)
	srv.HandleRouteFunc("/download", download.ManageVideownload)
	srv.HandleRouteFunc("/stream", streamFile)
	// create a channel of one byte
	sigChan := make(chan os.Signal, 1)
	// if os notices any of these sends, sends to channel
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	//as channel recieve any of these channel it'll shutdown connection gracefully
	go func(){
		<-sigChan
		fmt.Println("\nShutting down")
		srv.ShutDown()
		os.Exit(0)
	}()

	srv.StartServer(PORT_NO)
}

func Home(req *socket.Request, res *socket.Response){
bytesArray, err := os.ReadFile("../index/index.html")
if err != nil {
	fmt.Printf("error in Home:=%v \n", err)
}
data := string(bytesArray)
res.SendHTML(200,data)
}

// func download(req *socket.Request, res *socket.Response){
// res.InitChunkStream(200, "text/plain")
// // opening video file
// f, err := os.Open("../World.mp4")
// if err != nil {
// 	fmt.Printf("failed to open Video file:= %v\n", err)
// }

// defer f.Close()

// buffer := make([]byte, 4096)
// for {
// 	n, err := f.Read(buffer)
// 	if n > 0 {
// 		res.SendChunk(buffer[:n])
// 	}
// 	if err != nil {
// 		break
// 	}
// }
// }

func streamFile(req *socket.Request, res *socket.Response)  {
	videoData, err := os.ReadFile("../World.mp4")
	if err != nil {
		fmt.Printf("Error while streaming File:= \n",err)
	}
	//if rnage is not available then will stream buffer as *kb
	rangeHeader := req.Headers["Range"]
	if rangeHeader == "" {
// as the reange header is empty will send whole video chunked transfer
res.InitChunkStream(200, "video/mp4")
var chunkSizeKB int = 8192
for i := 0; i < len(videoData);  i += chunkSizeKB {
	end := i + chunkSizeKB
	if end > len(videoData) {
		end = len(videoData)
	}
	res.SendChunk(videoData[i:end])
}
res.StopSendindChunks()
	}else {
		start, end := utils.ParseRange(rangeHeader, len(videoData))
			// Send 206 Partial Content
		response := fmt.Sprintf("HTTP/1.1 206 Partial Content\r\n")
		response += "Content-Type: video/mp4\r\n"
		response += fmt.Sprintf("Content-Range: bytes %d-%d/%d\r\n", start, end-1, len(videoData))
		response += fmt.Sprintf("Content-Length: %d\r\n", end-start)
		response += "Accept-Ranges: bytes\r\n"
		response += "\r\n"

		syscall.Write(res.ConnFD, []byte(response))
		syscall.Write(res.ConnFD, videoData[start:end])
	}
}
