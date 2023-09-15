package service

import (
	g "InstantMessaging/app/global"
	"InstantMessaging/app/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"gopkg.in/fatih/set.v0"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	// 存储用户连接的映射关系
	clientMap = make(map[int64]*model.Node)
	// 读写锁，用于保护 clientMap 的并发读写
	rwLocker sync.RWMutex
)

// Chat 处理 WebSocket 连接的函数 需要 ：发送者ID ，接受者ID ，消息类型，发送的内容，发送类型
func Chat(writer http.ResponseWriter, request *http.Request) {
	//1.  获取参数 并 检验 token 等合法性
	//token := query.Get("token")
	query := request.URL.Query() // 获取查询参数
	Id := query.Get("userId")    // 获取用户ID
	userId, _ := strconv.ParseInt(Id, 10, 64)

	isValida := true //checkToken()
	// 设置UpGrader
	var upgrader = websocket.Upgrader{
		//token 校验， 校验 Origin 头部，用于安全校验
		CheckOrigin: func(r *http.Request) bool {
			return isValida
		},
	}

	//2. 升级协议，返回WebSocket连接
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 初始化节点信息
	currentTime := uint64(time.Now().Unix())
	node := &model.Node{
		Conn:          conn,
		Addr:          conn.RemoteAddr().String(), //客户端地址
		HeartbeatTime: currentTime,                //心跳时间
		LoginTime:     currentTime,                //登录时间
		DataQueue:     make(chan []byte, 50),
		GroupSets:     set.New(set.ThreadSafe),
	}
	//3. 用户关系
	//4. 将用户 ID 与节点绑定，并加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	//5.完成发送逻辑
	go sendProc(node)
	//6.完成接受逻辑
	go recvProc(node)
	//7.加入在线用户到缓存
	SetUserOnlineInfo("online_"+Id, []byte(node.Addr), time.Duration(viper.GetInt("timeout.RedisOnlineTime"))*time.Hour)
	sendMsg(userId, []byte("Welcome to CQUPT"))

}

// sendProc 处理发送消息的逻辑
func sendProc(node *model.Node) {
	for {
		// 监听DataQueue通道，一旦有数据可读取，就执行以下逻辑
		select {
		case data := <-node.DataQueue:
			fmt.Println("[ws]sendProc >>>> msg :", string(data))
			// 使用WebSocket连接发送消息
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// recvProc 处理接收消息的逻辑
func recvProc(node *model.Node) {
	for {
		// 从WebSocket连接中读取消息
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		// 处理接收到的数据
		handleReceivedData(node, data)
	}
}

// handleReceivedData 处理接收到的数据
func handleReceivedData(node *model.Node, data []byte) {
	msg := model.Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	//心跳检测 msg.Media == -1 || msg.Type == 3
	if msg.Type == 3 {
		currentTime := uint64(time.Now().Unix())
		node.Heartbeat(currentTime)
	} else {
		// 调用dispatch函数进行处理
		dispatch(data)
		fmt.Println("[ws] recvProc <<<<< ", string(data))
	}
}

// dispatch 处理后端调度逻辑
func dispatch(data []byte) {
	msg := model.Message{}
	msg.CreateTime = uint64(time.Now().Unix())
	err := json.Unmarshal(data, &msg) // 解析接收到的JSON数据
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1: //私信
		fmt.Println("dispatch  data :", string(data))
		sendMsg(msg.TargetId, data) // 发送消息给指定目标
	case 2: //群发
		sendGroupMsg(msg.TargetId, data) //发送的群ID, 消息内容, 发送消息给指定群组

	}
}

// sendGroupMsg 群发消息
func sendGroupMsg(targetId int64, msg []byte) {
	fmt.Println("开始群发消息")
	// 根据群组ID查找成员用户ID列表
	userIds := SearchUserByGroupId(uint(targetId))
	for i := 0; i < len(userIds); i++ {
		//排除给自己的
		if targetId != int64(userIds[i]) {
			sendMsg(int64(userIds[i]), msg) // 发送消息给群组中的其他成员
		}
	}
}

// sendMsg 发送消息
func sendMsg(targetId int64, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[targetId] // 根据目标ID查找对应的节点
	rwLocker.RUnlock()

	jsonMsg := model.Message{}
	err := json.Unmarshal(msg, &jsonMsg) // 解析消息内容
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx := context.Background()
	targetIdStr := strconv.Itoa(int(targetId))
	userIdStr := strconv.Itoa(int(jsonMsg.UserId))
	jsonMsg.CreateTime = uint64(time.Now().Unix())

	r, err := g.Rdb.Get(ctx, "online_"+userIdStr).Result() // 检查用户是否在线
	if err != nil {
		fmt.Println(err)
	}
	if r != "" {
		if ok {
			// 将消息放入目标节点的DataQueue通道中，等待发送
			fmt.Println("sendMsg >>> targetId: ", targetId, "  msg:", string(msg))
			node.DataQueue <- msg
		}
	}

	// 将消息按照一定的顺序存储到 Redis 缓存中，以便后续可以按照特定的顺序获取消息记录。
	var key string
	if targetId > jsonMsg.UserId {
		key = "msg_" + userIdStr + "_" + targetIdStr
	} else {
		key = "msg_" + targetIdStr + "_" + userIdStr
	}
	res, err := g.Rdb.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		fmt.Println(err)
	}
	score := float64(cap(res)) + 1                                                  // 计算新消息的分数。分数的计算是在已有消息记录的基础上加一
	ress, err := g.Rdb.ZAdd(ctx, key, &redis.Z{Score: score, Member: msg}).Result() // 将消息记录添加到有序集合中，同时指定了新消息的分数
	//res, err := utils.Red.Do(ctx, "zadd", key, 1, jsonMsg).Result() //备用 后续拓展 记录完整msg
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ress)
}

// RedisMsg 获取缓存里面的消息
func RedisMsg(userIdA int64, userIdB int64, start int64, end int64, isRev bool) []string {
	ctx := context.Background()
	userIdStr := strconv.Itoa(int(userIdA))
	targetIdStr := strconv.Itoa(int(userIdB))
	var key string
	if userIdA > userIdB {
		key = "msg_" + targetIdStr + "_" + userIdStr
	} else {
		key = "msg_" + userIdStr + "_" + targetIdStr
	}

	var rels []string
	var err error
	if isRev {
		rels, err = g.Rdb.ZRange(ctx, key, start, end).Result() // 从Redis中获取消息记录
	} else {
		rels, err = g.Rdb.ZRevRange(ctx, key, start, end).Result() // 从Redis中获取消息记录（反向）
	}
	if err != nil {
		fmt.Println(err) // 没有找到消息记录
	}
	return rels
}
