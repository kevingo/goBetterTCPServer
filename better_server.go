package main

import (
    "flag"
    "fmt"
    "net"
    "syscall"
)

const maxRead = 25

func main() {
    flag.Parse()
    if flag.NArg() != 2 {
        println("oops, please specify host abd port")
        return;
    }

    hostAndPort := fmt.Sprintf("%s:%s", flag.Arg(0), flag.Arg(1))
    listener := initServer(hostAndPort)

    for {
        conn, err := listener.Accept()
        checkError(err)
        go handler(conn)
    }
}

func initServer(hostAndPort string) *net.TCPListener {
    serverAddr, err := net.ResolveTCPAddr("tcp", hostAndPort)
    checkError(err)
    listener, err := net.ListenTCP("tcp", serverAddr)
    checkError(err)
    fmt.Println("Listening to:", listener.Addr().String())
    return listener
}

func handler(conn net.Conn) {
    from := conn.RemoteAddr().String()
    println("Client Connection from: ", from);
    welcome(conn)

    for {
        println("handler")
        var buf []byte = make([]byte, maxRead+1)
        length, err := conn.Read(buf[0:maxRead])
        buf[maxRead] = 0

        switch err {
            case nil:
                handleMsg(length, err, buf)
            case syscall.EAGAIN:
                continue
            default:
                goto DISCONNECT
        }
    }

    DISCONNECT:
        err := conn.Close()
        fmt.Println("Close connection")
        checkError(err)
}

func handleMsg(length int, err error, msg []byte) {
    if length > 0 {
        print("<", length, ":")
        for i := 0; ; i++ {
            if msg[i] == 0 {
                break;
            }
            fmt.Printf("%c", msg[i])
        }
        print(">")
    }
}

func welcome(conn net.Conn) {
    msg := "Hello !!"
    conn.Write([]byte(msg))
}

func checkError(error error) {
    if error != nil {
        panic("ERROR: " + error.Error())
    }
}