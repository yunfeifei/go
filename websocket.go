package main
  
import (
    "log"
    "fmt"
    "net/http"
    "github.com/gorilla/websocket"
)

var (
    upgrader = websocket.Upgrader {
        // 读取存储空间大小
        ReadBufferSize:1024,
        // 写入存储空间大小
        WriteBufferSize:1024,
        // 允许跨域
        CheckOrigin: func(r *http.Request) bool {
            return true
        },
    }
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
    var (
        conn *websocket.Conn
        err error
        messageType int //websocket.BinaryMessage or websocket.TextMessage
        p []byte
    )
    // 完成http应答，在httpheader中放下如下参数
    if  conn, err = upgrader.Upgrade(w, r, nil);err != nil {
        return // 获取连接失败直接返回
    }

    for {
        // 只能发送Text, Binary 类型的数据,下划线意思是忽略这个变量.
        messageType, p, err = conn.ReadMessage()
        if err != nil {
            goto ERR // 跳转到关闭连接
        }
        if err = conn.WriteMessage(messageType, p); err != nil {
            goto ERR // 发送消息失败，关闭连接
        }
        if err = conn.WriteMessage(messageType, []byte("Server send end!")); err != nil {
            goto ERR // 发送消息失败，关闭连接
        }
    }

    ERR:
      conn.Close()
}

func main()  {
    fmt.Println("websocket start at 127.0.0.1:8080...")
    http.HandleFunc("/", wsHandler)
    err := http.ListenAndServe("127.0.0.1:8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe", err.Error())
    }
}