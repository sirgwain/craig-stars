tool
extends EditorPlugin

func _enter_tree():
	var texture = preload("res://addons/CraigStarsComponents/icon.svg")
	var button_texture = preload("res://addons/CraigStarsComponents/assets/icon_button.svg")
	add_custom_type("ProductionQueueItemsTable", "MarginContainer", preload("res://addons/CraigStarsComponents/src/ProductionQueueItemsTable.cs"), texture)
	add_custom_type("CSButton", "Button", preload("res://addons/CraigStarsComponents/src/CSButton.cs"), button_texture)
	
func _exit_tree():
	remove_custom_type("ProductionQueueItemsTable")
	remove_custom_type("CSButton")
	
