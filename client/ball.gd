extends CharacterBody2D

const INITIAL_SPEED: float = 400.0 # Tốc độ ban đầu của ball
var speed: float           = INITIAL_SPEED
var direction: Vector2     = Vector2.ZERO
var network: Node
var game_started := false
var is_host := false
var lerp_speed := 0.3
var target_position := Vector2.ZERO
var last_sync_position := Vector2.ZERO
var last_sync_direction := Vector2.ZERO
var last_sync_speed := 0.0

func _ready() -> void:
	network = get_node("/root/Network")
	network.game_started.connect(_on_game_started)
	network.ball_updated.connect(_on_ball_updated)
	reset_ball()
	
func _on_game_started(player1_id: String, player2_id: String) -> void:
	game_started = true	
	# Player1 sẽ là host
	is_host = network.player_id == player1_id
	
func _on_ball_updated(x: float, y: float, dir_x: float, dir_y: float, spd: float) -> void:
	if not is_host:
		last_sync_position = Vector2(x, y)
		last_sync_direction = Vector2(dir_x, dir_y)
		last_sync_speed = spd
		position = last_sync_position

# Reset ball về vị trí giữa màn hình
func reset_ball() -> void:
	position = Vector2(get_viewport_rect().size.x / 2, get_viewport_rect().size.y / 2)
	direction = Vector2([-1, 1].pick_random(), randf_range(-0.8, 0.8)).normalized()
	speed = INITIAL_SPEED

func _physics_process(delta: float) -> void:
	if not game_started:
		return
	
	if is_host:
		var collison := move_and_collide(direction * speed * delta)
		if collison:
			direction = direction.bounce(collison.get_normal())
			speed *= 1.05
			
			if position.x <= 0 or position.x >= get_viewport_rect().size.x:
				reset_ball()
		
		# Gửi vị trí liên tục, không chỉ khi va chạm
		network.send_ball_position(position.x, position.y, direction, speed)
	else:
		position += last_sync_direction * last_sync_speed * delta
