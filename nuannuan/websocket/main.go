package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"nuannuan/database"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 用户连接池
var clients = make(map[string]*websocket.Conn)
var clientsMutex sync.Mutex

// RabbitMQ连接和通道
var rabbitConn *amqp.Connection
var rabbitChannel *amqp.Channel

// MySQL 数据库连接
var db *gorm.DB

// Message 模型
type Message struct {
	ID         uint   `gorm:"primary_key"`
	FromUserID string `gorm:"type:varchar(255);not null"`
	ToUserID   string `gorm:"type:varchar(255);not null"`
	Message    string `gorm:"type:text;not null"`
	IsRead     bool   `gorm:"default:false"`
	CreatedAt  time.Time
}

// 初始化RabbitMQ
func initRabbitMQ() {
	var err error
	// 连接到RabbitMQ TODO:替换成你的RabbitMQ IP
	rabbitConn, err = amqp.Dial("amqp://guest:guest@IP:端口/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	// 创建通道
	rabbitChannel, err = rabbitConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	// 声明队列
	_, err = rabbitChannel.QueueDeclare(
		"message_queue", // 队列名称
		true,            // 持久化
		false,           // 自动删除
		false,           // 独占
		false,           // 无等待
		nil,             // 参数
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}
}

// 初始化MySQL
func initMySQL() {
	db = database.GetDB()
}

// 存储消息到数据库
func saveMessage(fromUserID, toUserID, message string) error {
	msg := Message{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Message:    message,
		IsRead:     false,
		CreatedAt:  time.Now(),
	}
	return db.Create(&msg).Error
}

// 获取用户的未读消息
func getUnreadMessages(userID string) ([]map[string]string, error) {
	var messages []Message
	err := db.Where("to_user_id = ? AND is_read = ?", userID, false).Find(&messages).Error
	if err != nil {
		return nil, err
	}

	// 将消息转换为 map
	var result []map[string]string
	for _, msg := range messages {
		result = append(result, map[string]string{
			"fromUserID": msg.FromUserID,
			"message":    msg.Message,
		})
	}

	// 标记消息为已读
	err = db.Model(&Message{}).Where("to_user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

// WebSocket 服务
func startWebSocketServer() {
	r := gin.Default()
	r.GET("/ws", handleConnections)
	log.Println("WebSocket server started on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start WebSocket server: %v", err)
	}
}

// 处理 WebSocket 连接
func handleConnections(c *gin.Context) {
	// 从请求中获取用户ID（假设通过查询参数传递）
	userID := c.Query("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID is required"})
		return
	}

	// 升级为WebSocket连接
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer ws.Close()

	// 将用户连接添加到连接池
	clientsMutex.Lock()
	clients[userID] = ws
	clientsMutex.Unlock()

	log.Printf("User %s connected", userID)

	// 获取用户的未读消息并发送
	unreadMessages, err := getUnreadMessages(userID)
	if err != nil {
		log.Printf("Failed to get unread messages for %s: %v", userID, err)
	} else {
		for _, msg := range unreadMessages {
			err := ws.WriteJSON(msg)
			if err != nil {
				log.Printf("Error sending unread message to %s: %v", userID, err)
			}
		}
	}

	// 监听消息
	for {
		var msg struct {
			ToUserID string `json:"toUserID"` // 目标用户ID
			Message  string `json:"message"`  // 消息内容
		}
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("User %s disconnected: %v", userID, err)
			clientsMutex.Lock()
			delete(clients, userID)
			clientsMutex.Unlock()
			break
		}

		log.Printf("Received message from %s to %s: %s", userID, msg.ToUserID, msg.Message)

		// 查找目标用户连接
		clientsMutex.Lock()
		targetClient, ok := clients[msg.ToUserID]
		clientsMutex.Unlock()

		if ok {
			// 发送消息给目标用户
			err := targetClient.WriteJSON(map[string]string{
				"fromUserID": userID,
				"message":    msg.Message,
			})
			if err != nil {
				log.Printf("Error sending message to %s: %v", msg.ToUserID, err)
			}
		} else {
			log.Printf("User %s is offline, storing message in RabbitMQ", msg.ToUserID)
			// 将消息发送到RabbitMQ队列
			messageData := map[string]string{
				"fromUserID": userID,
				"toUserID":   msg.ToUserID,
				"message":    msg.Message,
			}
			messageBody, err := json.Marshal(messageData)
			if err != nil {
				log.Printf("Failed to marshal message: %v", err)
				continue
			}

			err = rabbitChannel.Publish(
				"",              // 交换机
				"message_queue", // 队列名称
				false,           // 强制
				false,           // 立即
				amqp.Publishing{
					ContentType: "application/json",
					Body:        messageBody,
				},
			)
			if err != nil {
				log.Printf("Failed to publish message to RabbitMQ: %v", err)
			}
		}
	}
}

// RabbitMQ 消息消费者服务
func startMessageConsumer() {
	// 消费消息
	msgs, err := rabbitChannel.Consume(
		"message_queue", // 队列名称
		"",              // 消费者名称
		true,            // 自动确认
		false,           // 独占
		false,           // 无等待
		false,           // 无额外参数
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Println("Waiting for messages...")

	// 处理消息
	for msg := range msgs {
		var messageData map[string]string
		err := json.Unmarshal(msg.Body, &messageData)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// 存储消息到数据库
		err = saveMessage(messageData["fromUserID"], messageData["toUserID"], messageData["message"])
		if err != nil {
			log.Printf("Failed to save message to database: %v", err)
		} else {
			log.Printf("Saved message from %s to %s: %s", messageData["fromUserID"], messageData["toUserID"], messageData["message"])
		}
	}
}

func main() {
	// 初始化RabbitMQ
	initRabbitMQ()
	defer rabbitConn.Close()
	defer rabbitChannel.Close()

	// 初始化MySQL
	initMySQL()
	defer db.Close()

	// 启动 WebSocket 服务
	go startWebSocketServer()

	// 启动 RabbitMQ 消息消费者服务
	startMessageConsumer()
}
