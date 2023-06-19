export const clamp = (value: number, min: number, max: number) => {
	if (value < min) {
		return min;
	} else {
		if (value > max) {
			return max;
		}
	}
	return value;
};

export const radiansToDegrees = (radians: number): number => radians * (180/Math.PI)