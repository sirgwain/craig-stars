/**
 * clamp a value between a min and a max
 * @param value
 * @param min
 * @param max
 * @returns
 */
export const clamp = (value: number, min: number, max: number): number => {
	if (value < min) {
		return min;
	} else {
		if (value > max) {
			return max;
		}
	}
	return value;
};

/**
 * rollover a value if it goes over or under a min/max
 * This is used when cycling through a list of items
 * @param value the value to rollover
 * @param min the min, i.e. 0
 * @param max the max, i.e. planets.length
 * @returns max if value < 0, min if value > max, otherwise value
 */
export const rollover = (value: number, min: number, max: number): number => {
	if (value < 0) {
		return max;
	} else if (value > max) {
		return 0;
	}
	return value;
};

export const radiansToDegrees = (radians: number): number => radians * (180 / Math.PI);
