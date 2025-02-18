extends CharacterBody2D

const SPEED = 1200.0 # Tốc độ di chuyển của paddle
const FRICTION = 0.2
var is_local_player := false # Có phải là player local không
var network: Node # Reference tới network singleton
var last_sent_position: float = 300.0
@onready var logger: Logger = $"../RichTextLabel"

func _ready() -> void:
	# Lấy reference tới network singleton và kết nối signal
	network = get_node("/root/Network")
	network.player_input_received.connect(_on_player_input_received)

func _physics_process(delta: float) -> void:
	if is_local_player:
		var direction := Input.get_axis("move_up", "move_down")
		if direction != 0:
			velocity.y = direction * SPEED
		else:
			velocity.y	= 0
			
		move_and_slide()
		clamp_position()
		
		# Gửi vị trí nếu thay đổi đáng kể
		if abs(position.y - last_sent_position) > 1.0:
			network.send_paddle_position(position.y)
			last_sent_position = position.y
		
		
func clamp_position() -> void:
	var viewport_size := get_viewport_rect().size
	position.y = clamp(position.y, 80, viewport_size.y - 80)

# Xử lý khi nhận được vị trí paddle từ player khác
func _on_player_input_received(player_id: String, pos: float) -> void:
	if !is_local_player:
		if abs(pos - position.y) > 1.0:
			position.y = pos
