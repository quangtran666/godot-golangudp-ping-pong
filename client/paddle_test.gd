extends CharacterBody2D

const SPEED = 400.0 # Tốc độ di chuyển của paddle

func _ready() -> void:
	pass


func _physics_process(delta: float) -> void:
	var direction := Input.get_axis("move_up", "move_down")

	velocity.y += direction * SPEED * delta
	move_and_slide()
