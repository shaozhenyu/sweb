package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	listenPort = flag.String("port", "8899", "server listen port")
	syncFile   = flag.String("file", "", "transfer file")
	Host       = flag.String("host", "", "server host")
	syncSer    = flag.Bool("d", false, "server mode")
	syncFold   = flag.String("dir", "/tmp/gosync/", "recive sync fold ")
)

func main() {
	flag.Parse()
	if *syncSer {
		port := fmt.Sprintf(":%s", *listenPort)
		ln, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatal(err)
		}
		defer ln.Close()

		err = os.MkdirAll(*syncFold, 0755)
		if err != nil {
			log.Fatal(err)
		}

		Server(ln)
	} else {
		addr := fmt.Sprintf("%s:%s", *Host, *listenPort)
		clientSend(*syncFile, addr)
	}
}

func clientSend(filename string, addr string) {
	//filename = fmt.Sprintf("%s/%s", *syncFold, filename)
	finfo := getFileInfo(filename)
	newName := fmt.Sprintf("%s", finfo.fileName)

	cmdLine := fmt.Sprintf("upload %s %d %d %s", newName, finfo.fileSize, finfo.filePerm, finfo.fileMd5)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	conn.Write([]byte(cmdLine))
	conn.Write([]byte("\r\n"))

	fileHander, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	io.Copy(conn, fileHander)
	for {
		buffer := make([]byte, 1024)
		num, err := conn.Read(buffer)
		if err == nil && num > 0 {
			fmt.Println(string(buffer[:num]))
			break
		}
	}
}

func Server(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue
			}
			log.Fatal("network err", err)
		}
		go ServerHandler(conn)
	}
}

func ServerHandler(conn net.Conn) {
	defer conn.Close()
	status := 0
	var cmd *fileInfo
	var fSize int64
	var newFilename string
	var n int64

	for {
		buffer := make([]byte, 2048)
		num, err := conn.Read(buffer)
		numLen := int64(num)
		n = 0
		if status == 0 {
			n, cmd = cmdParse(buffer[:num])
			fmt.Println(string(cmd.fileName))
			newFilename = fmt.Sprintf("%s.newsync", cmd.fileName)
			fSize = cmd.fileSize
			status = 1
		}
		if status == 1 {
			last := numLen
			if fSize <= numLen-n {
				last = fSize + n
				status = 2
			}
			err = writeToFile(buffer[int(n):int(last)], newFilename, cmd.filePerm)
			if err != nil {
				log.Fatal(err)
			}
			fSize -= last - n
			if status == 2 {
				os.Remove(cmd.fileName)
				err = os.Rename(newFilename, cmd.fileName)
				if err != nil {
					log.Fatal(err)
				}

				//file change time
				fileHandle, err := os.Open(cmd.fileName)
				if err != nil {
					log.Fatal(err)
				}

				h := md5.New()
				io.Copy(h, fileHandle)
				newfMd5 := fmt.Sprintf("%x", h.Sum(nil))
				if newfMd5 == cmd.fileMd5 {
					sendInfo := fmt.Sprintf("%s sync success", cmd.fileName)
					conn.Write([]byte(sendInfo))
				} else {
					sendInfo := fmt.Sprintf("%s sync failed", cmd.fileName)
					conn.Write([]byte(sendInfo))
				}
			}
		}
	}
}

func writeToFile(data []byte, filename string, perm os.FileMode) error {
	writeFile, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, perm)
	if err != nil {
		return err
	}
	defer writeFile.Close()

	_, err = writeFile.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func cmdParse(buf []byte) (int64, *fileInfo) {
	i := len(buf) - 1
	if i < 2 {
		return 0, nil
	}
	fmt.Println(buf)
	if buf[i] == '\n' && buf[i-1] == '\r' {
		cmdLine := strings.Split(string(buf[:i-1]), " ")
		name := fmt.Sprintf("%s%s", *syncFold, cmdLine[1])
		size, _ := strconv.ParseInt(cmdLine[2], 10, 64)
		perm, _ := strconv.ParseInt(cmdLine[3], 10, 64)
		finfo := &fileInfo{
			fileName: name,
			fileSize: size,
			filePerm: os.FileMode(perm),
			fileMd5:  string(cmdLine[4]),
		}
		return int64(i + 1), finfo
	}
	return 0, nil
}

type fileInfo struct {
	fileName string
	fileSize int64
	filePerm os.FileMode
	fileMd5  string
}

func getFileInfo(name string) *fileInfo {
	info, err := os.Lstat(name)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		log.Fatal(err)
	}

	fi := &fileInfo{
		fileName: name,
		fileSize: info.Size(),
		filePerm: info.Mode().Perm(),
		fileMd5:  fmt.Sprintf("%x", h.Sum(nil)),
	}
	return fi
}
