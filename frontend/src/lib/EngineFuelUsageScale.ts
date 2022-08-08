export {}
// import * as Chart from 'chart.js';

// // This totally doesn't work, but it's a start on how to create a custom scale
// // in chart.js
// export class EngineFuelUsageScale extends Chart.Scale {
// 	determineDataLimits() {
// 		const { min, max } = this.getMinMax(true);
// 		this.min = isFinite(min) ? Math.max(0, min) : 0;
// 		this.max = isFinite(max) ? Math.max(0, max) : 0;
// 	}

// 	buildTicks() {
// 		return [
// 			{ value: 0, label: 'Zero' },
// 			{ value: 1, label: '25' },
// 			{ value: 2, label: '50' },
// 			{ value: 3, label: '100' },
// 			{ value: 4, label: '200' },
// 			{ value: 5, label: '400' },
// 			{ value: 6, label: '800' }
// 		];
// 	}

// 	configure() {
// 		super.configure();

// 		this._startValue = 0;
// 		this._valueRange = 7; // 7 ticks
// 	}

// 	getPixelForValue(value) {
// 		value > 0 && console.log('getPixelForValue value:', value);
// 		if (value === undefined || value === 0) {
// 			value = this.min;
// 		}

// 		const scaledUsage = Math.max(0, Math.log(parseFloat(value) / 25) / Math.log(2) + 1);

// 		const pixel = this.getPixelForDecimal(scaledUsage);

// 		console.log(`getPixelForValue scaledUsage: ${scaledUsage} pixel: ${pixel}`);
// 		return pixel;
// 	}

// 	getValueForPixel(pixel) {
// 		const decimal = this.getDecimalForPixel(pixel);

// 		console.log('getValueForPixel decimal:', decimal);
// 		return 2 * decimal * this._valueRange;
// 	}
// }
// EngineFuelUsageScale.id = 'log2';
// // Log2Scale.defaults = {};

// Chart.Chart.register(EngineFuelUsageScale);
