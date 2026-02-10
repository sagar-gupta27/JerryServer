package download 

import ("fmt"
"syscall"
"httpServer/src/socket"
"os")

func ManageVideownload(req *socket.Request, res *socket.Response){
	videoData, err :=os.ReadFile("../World.mp4")
	if err != nil {
		fmt.Printf("error while ManageVideownload:= %v \n", err)
	}

	response := "HTTP/1.1 200 OK\r\n"
	response += "Content-Type: video/mp4\r\n"
	response += "Content-Disposition: attachment; filename=\"demo-video.mp4\"\r\n"
	response += fmt.Sprintf("Content-Length: %d\r\n", len(videoData))
	response += "\r\n"
	syscall.Write(res.ConnFD, []byte(response))
	var chunkSizeKB int = 8192
	for i := 0; i < len(videoData);  i += chunkSizeKB {
	end := i + chunkSizeKB
	if end > len(videoData) {
		end = len(videoData)
	}
	syscall.Write(res.ConnFD,videoData[i:end])
}
}