extends Node2D

# References tới các node
@onready var network: Node = get_node("/root/Network")
@onready var player1_paddle = $Player1Paddle
@onready var player2_paddle = $Player2Paddle	
@onready var logger: Logger = $RichTextLabel

var local_player_id: String

func _ready() -> void:
	# Kết nối signal và join game
	network.game_started.connect(_on_game_started)
	network.join_game()
	
# Xử lý khi game bắt đầu
func _on_game_started(player1_id: String, player2_id: String) -> void:
	# Xác định player local
	local_player_id = network.player_id

	# Set up paddle controls dựa vào player ID
	if local_player_id == player1_id:
		player1_paddle.is_local_player = true
	elif local_player_id == player2_id:
		player2_paddle.is_local_player = true
		
	## Log thông tin game start
	logger.info("Player 1: " + str(player1_paddle.is_local_player))
	logger.info("Player 2: " + str(player2_paddle.is_local_player))
		
		
		
	
