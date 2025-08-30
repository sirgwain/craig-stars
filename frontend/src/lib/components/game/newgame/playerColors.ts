const colors = [
	'#0000FF', // Player 1 - Blue
	'#FF0000', // Bright Red
	'#00FFFF', // Cyan
	'#00FF00', // Lime
	'#FF00FF', // Magenta
	'#FFFF00', // Yellow
	'#FF7F00', // Bright Orange
	'#00FF7F', // Spring Green
	'#7FFF00', // Chartreuse
	'#FF1493', // Deep Pink
	'#FFD700', // Gold
	'#40E0D0', // Turquoise
	'#9ACD32', // Yellow Green
	'#FF69B4', // Hot Pink
	'#1E90FF', // Dodger Blue
	'#FF4500', // Orange Red
	'#ADFF2F', // Green Yellow
	'#BA55D3', // Medium Orchid
	'#87CEFA', // Light Sky Blue
	'#00CED1', // Dark Turquoise
	'#FFA500', // Orange
	'#32CD32', // Bright Green
	'#DC143C', // Crimson
	'#8A2BE2', // Blue Violet
	'#7FFFD4' // Aquamarine
];

export const getColor = (index: number) =>
	index < colors.length ? colors[index] : '#' + Math.floor(Math.random() * 16777215).toString(16);

export const getFirstAvailableColor = (usedColors: Set<string>) =>
	colors.find((c) => !usedColors.has(c)) ?? '#' + Math.floor(Math.random() * 16777215).toString(16);
