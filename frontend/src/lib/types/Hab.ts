export interface Hab {
	grav?: number;
	temp?: number;
	rad?: number;
}

export function getGravString(grav: number): string {
	let result = 0;
	const tmp = Math.abs(grav - 50);
	if (tmp <= 25) result = (tmp + 25) * 4;
	else result = tmp * 24 - 400;
	if (grav < 50) result = 10000 / result;

	const value = result / 100 + (result % 100) / 100.0;
	return `${value.toFixed(2)}g`;
}

export function getTempString(temp: number): string {
	return `${(temp - 50) * 4}°C`;
}

export function getRadString(rad: number): string {
	return `${rad}mR`;
}
