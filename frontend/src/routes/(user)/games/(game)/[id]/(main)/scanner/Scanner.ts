export type ViewportCoords = {
	x: number;
	y: number;
	visible: boolean;
};

// return the viewport coords and whether they are in the viewport
export function getViewportCoords(
	x: number,
	y: number,
	width: number,
	height: number,
	padding = 0
): ViewportCoords {
	return {
		x,
		y,
		visible: x > -padding && y > -padding && x < width + padding && y < height + padding
	};
}
