const colors = [
	'#0000FF',
	'#C33232',
	'#1F8BA7',
	'#43A43E',
	'#8D29CB',
	'#B88628',
	'#FF4500',
	'#FF8C00',
	'#008000',
	'#00FA9A',
	'#7FFFD4',
	'#8A2BE2',
	'#FF1493',
	'#D2691E',
	'#F0FFF0'
];
export const getColor = (index: number) =>
	index < colors.length ? colors[index] : '#' + Math.floor(Math.random() * 16777215).toString(16);
