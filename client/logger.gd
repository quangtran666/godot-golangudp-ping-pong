class_name Logger
extends RichTextLabel

func info(message: String, color: Color = Color.WHITE) -> void:
	append_text("[color=#%s]%s[/color]\n" % [color.to_html(false), message])
