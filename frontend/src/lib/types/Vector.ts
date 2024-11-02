export interface Vector {
	x: number;
	y: number;
}

export const emptyVector = { x: 0, y: 0 };

export const equal = (v1: Vector, v2: Vector) => v1.x === v2.x && v1.y === v2.y;

// compute the distance between two vectors
export const distance = (v1: Vector, v2?: Vector): number =>
	v2 ? Math.sqrt((v1.x - v2.x) * (v1.x - v2.x) + (v1.y - v2.y) * (v1.y - v2.y)) : 0;

export const perpendicular = (v: Vector): Vector => {
	return { x: v.y, y: -v.x };
};

// dot product of two vectors
export const dot = (v1: Vector, v2: Vector): number => v1.x * v2.x + v1.y * v2.y;
export const determinant = (v1: Vector, v2: Vector): number => v1.x * v2.y - v1.y * v2.x;

export const lengthSquared = (v: Vector): number => v.x * v.x + v.y * v.y;

export const normalized = (from: Vector): Vector => {
	const v = { x: from.x, y: from.y };
	const lengthsq = lengthSquared(v);

	if (lengthsq === 0) {
		v.x = 0;
		v.y = 0;
	} else {
		const length = Math.sqrt(lengthsq);
		v.x /= length;
		v.y /= length;
	}
	return v;
};

export const subtract = (from: Vector, to: Vector): Vector => {
	return { x: from.x - to.x, y: from.y - to.y };
};
