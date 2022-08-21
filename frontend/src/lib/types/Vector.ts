export interface Vector {
	x: number;
	y: number;
}

// compute the distance between two vectors
export const distance = (v1: Vector, v2: Vector): number =>
	Math.sqrt((v1.x - v2.x) * (v1.x - v2.x) + (v1.y - v2.y) * (v1.y - v2.y));
