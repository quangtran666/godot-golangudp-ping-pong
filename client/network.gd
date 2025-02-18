extends Node

# Signals để thông báo các sự kiện network
signal game_started(player1_id: String, player2_id: String)
signal player_input_received(player_id: String, position: float)
signal ball_updated(x: float, y: float, dir_x: float, dir_y: float, spd: float)
signal game_ended(winner_id: String)
# Biến để quản lý kết nối UDP
var socket         := PacketPeerUDP.new()
var server_address := "2602:fbaf:848:1::10"
#var server_address := "127.0.0.1"
var server_port    := 50358
#var server_port    := 8080
var player_id: String


func _ready() -> void:
	# Tạo random player ID và kết nối tới server
	player_id = str(randi())
	socket.connect_to_host(server_address, server_port)


func _process(delta: float) -> void:
	# Kiểm tra và xử lý các packet nhận được
	while socket.get_available_packet_count() > 0:
		var packet = socket.get_packet()
		var json   = JSON.parse_string(packet.get_string_from_utf8())
		handle_message(json)


# Gửi request tham gia game		
func join_game() -> void:
	send_message("join_game", {
		"player_id": player_id,
	})


# Gửi vị trí paddle	
func send_paddle_position(position: float) -> void:
	send_message("player_input", {
		"player_id": player_id,
		"position": position,
	})
	
func send_ball_position(x: float, y: float, direction: Vector2, speed: float) -> void:
	send_message("ball_update", {
		"x": x,
		"y": y,
		"direction_x": direction.x,
		"direction_y": direction.y,
		"speed": speed
	})
	
# Hàm helper để gửi message
func send_message(type: String, payload: Dictionary) -> void:
	var message: Dictionary = {
		"type": type,
		"payload": payload,
	}
	
	var json: String = JSON.stringify(message)
	socket.put_packet(json.to_utf8_buffer())

# Xử lý các message nhận được từ server
func handle_message(message: Dictionary) -> void:
	match message.type:
		"game_start":
			game_started.emit(message.payload.player1_id, message.payload.player2_id)
		"player_input":
			player_input_received.emit(message.payload.player_id, message.payload.position)
		"ball_update":
			ball_updated.emit(
				message.payload.x,
				message.payload.y,
				message.payload.direction_x,
				message.payload.direction_y,
				message.payload.speed
			)
		"game_end":
			game_ended.emit(message.payload.winner_id)
		
