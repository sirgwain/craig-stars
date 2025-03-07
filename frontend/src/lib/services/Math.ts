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
 * @param max the max, i.e. planets.length
 * @returns max if value < 0, min if value > max, otherwise value
 */
export const rollover = (value: number, max: number): number => {
	if (value < 0) {
		return max;
	} else if (value > max) {
		return 0;
	}
	return value;
};

export const radiansToDegrees = (radians: number): number => radians * (180 / Math.PI);

/**
 * Round a value to a multiple of 100 using the specified rounding function
 * and return the result as an integer.
 * @param value the value beind rounded
 * @param roundFunc the function to round the result with; defaults to math.round
 * @returns value rounded to the nearest multiple of 100
 */
export const roundTo100 = (value: number, roundFunc: (a: number) => number = Math.round): number => {
	return roundFunc(value / 100) * 100;
}
