package network

import (
	"encoding/json"
	"log"
	"net"
	"sync"

	"github.com/quangtran666/godot-golangudp-ping-pong/internal/game"
	"github.com/quangtran666/godot-golangudp-ping-pong/internal/protocol"
)

// Server quản lý kết nối UDP và điều phối messages
type Server struct {
	conn    *net.UDPConn
	rooms   map[string]*game.Room   // Map lưu trữ các phòng chơi
	players map[string]*game.Player // Map lưu trữ các player
	mutex   sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		rooms:   make(map[string]*game.Room),
		players: make(map[string]*game.Player),
	}
}

// Start khởi động server và lắng nghe kết nối tới
func (server *Server) Start(address string) error {
	// Resolve địa chỉ UDP
	add, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	// Tạo socket UDP
	conn, err := net.ListenUDP("udp", add)
	if err != nil {
		return err
	}

	server.conn = conn
	log.Printf("Server started on %s", address)

	return server.handleConnections()
}

const MAX_BUFFER_SIZE = 8192 // Tăng kích thước buffer

func (server *Server) handleConnections() error {
	buffer := make([]byte, MAX_BUFFER_SIZE)

	for {
		n, remoteAddr, err := server.conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			continue
		}

		// Kiểm tra xem message có bị truncate không
		if n >= MAX_BUFFER_SIZE {
			log.Printf("Warning: Message too large, may be truncated")
			continue
		}

		// Chỉ xử lý phần data thực tế nhận được
		data := make([]byte, n)
		copy(data, buffer[:n])

		// Validate JSON trước khi xử lý
		if !json.Valid(data) {
			log.Printf("Invalid JSON received: %s", string(data))
			continue
		}

		go server.handleMessage(data, remoteAddr)
	}
}

// handleMessage xử lý các message nhận được
func (server *Server) handleMessage(data []byte, addr *net.UDPAddr) {
	var msg protocol.Message
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("Error unmarshalling message: %v, Raw data: %s", err, string(data))
		return
	}

	// Xử lý message dựa vào type
	switch msg.Type {
	case protocol.MsgTypeJoinGame:
		server.handleJoinGame(msg.Payload, addr)
	case protocol.MsgTypePlayerInput:
		server.handlePlayerInput(msg.Payload, addr)
	case protocol.MsgTypeBallUpdate:
		server.handleBallUpdate(msg.Payload, addr)
	}
}

// handleJoinGame xử lý khi player muốn tham gia game
func (server *Server) handleJoinGame(data json.RawMessage, addr *net.UDPAddr) {
	var joinMsg protocol.JoinGameMessage
	if err := json.Unmarshal(data, &joinMsg); err != nil {
		log.Printf("Error unmarshalling JoinGameMessage: %v", err)
		return
	}

	player := game.NewPlayer(joinMsg.PlayerID, addr)

	server.mutex.Lock()
	server.players[player.ID] = player

	// Tìm phòng trống hoặc tạo phòng mới
	var targetRoom *game.Room
	for _, room := range server.rooms {
		if !room.IsFull() {
			targetRoom = room
			break
		}
	}

	if targetRoom == nil {
		targetRoom = game.NewRoom(joinMsg.PlayerID + "_room")
		server.rooms[targetRoom.ID] = targetRoom
	}
	server.mutex.Unlock()

	// Nếu phòng đủ 2 người, bắt đầu game
	if targetRoom.AddPlayer(player) && targetRoom.IsFull() {
		startMsg := protocol.GameStartMessage{
			Player1ID: targetRoom.Player1.ID,
			Player2ID: targetRoom.Player2.ID,
		}
		log.Printf("Game started between %s and %s", targetRoom.Player1.ID, targetRoom.Player2.ID)
		server.broadcastToRoom(targetRoom, protocol.MsgTypeGameStart, startMsg)
	}
}

func (server *Server) handlePlayerInput(payload json.RawMessage, addr *net.UDPAddr) {
	var inputMsg protocol.PlayerInputMessage
	if err := json.Unmarshal(payload, &inputMsg); err != nil {
		log.Printf("Error unmarshalling PlayerInputMessage: %v", err)
		return
	}

	log.Printf("Player %s input: %f", inputMsg.PlayerID, inputMsg.Position)

	server.mutex.RLock()
	players, exists := server.players[inputMsg.PlayerID]
	server.mutex.RUnlock()

	if !exists || players.Room == nil {
		return
	}

	// Relay input tới player khác trong phòng
	if players.Room.Player1 != nil && players.Room.Player1.ID != inputMsg.PlayerID {
		server.sendToPlayer(players.Room.Player1, protocol.MsgTypePlayerInput, inputMsg)
	}
	if players.Room.Player2 != nil && players.Room.Player2.ID != inputMsg.PlayerID {
		server.sendToPlayer(players.Room.Player2, protocol.MsgTypePlayerInput, inputMsg)
	}
}

// broadcastToRoom gửi message tới tất cả players trong phòng
func (server *Server) broadcastToRoom(room *game.Room, msgType string, payload interface{}) {
	data, err := json.Marshal(protocol.Message{
		Type:    msgType,
		Payload: json.RawMessage(mustMarshal(payload)),
	})
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		return
	}

	// Gửi message tới cả 2 players
	if room.Player1 != nil {
		server.conn.WriteToUDP(data, room.Player1.Addr)
	}
	if room.Player2 != nil {
		server.conn.WriteToUDP(data, room.Player2.Addr)
	}
}

func (server *Server) sendToPlayer(player *game.Player, msgType string, payload interface{}) {
	data, err := json.Marshal(protocol.Message{
		Type:    msgType,
		Payload: json.RawMessage(mustMarshal(payload)),
	})
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		return
	}
	server.conn.WriteToUDP(data, player.Addr)
}

func (server *Server) handleBallUpdate(payload json.RawMessage, addr *net.UDPAddr) {
	var ballMsg protocol.BallUpdateMessage
	if err := json.Unmarshal(payload, &ballMsg); err != nil {
		log.Printf("Error unmarshalling BallUpdateMessage: %v", err)
		return
	}

	// Tìm player và room tương ứng dựa vào địa chỉ UDP
	var targetRoom *game.Room
	server.mutex.RLock()
	for _, room := range server.rooms {
		if (room.Player1 != nil && room.Player1.Addr.String() == addr.String()) ||
			(room.Player2 != nil && room.Player2.Addr.String() == addr.String()) {
			targetRoom = room
			break
		}
	}
	server.mutex.RUnlock()

	if targetRoom == nil {
		return
	}

	// Broadcast vị trí bóng mới tới tất cả players trong room
	server.broadcastToRoom(targetRoom, protocol.MsgTypeBallUpdate, ballMsg)
}

func mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
