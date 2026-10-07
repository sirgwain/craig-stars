import type { Position } from '$lib/types/MapObject';

export type ScreenPoint = { x: number; y: number };

/**
 * ScannerView converts world (light year) coordinates into screen (css pixel) coordinates.
 *
 * Like the original Stars! scanner, every object is projected to a single pixel-snapped
 * anchor, and everything drawn for that object (icon, name, fleet count, selection arrow)
 * is positioned at fixed pixel offsets from that anchor so it all lines up at any zoom.
 */
export class ScannerView {
	constructor(
		// css pixels per light year before zooming
		readonly sx: number,
		readonly sy: number,
		// d3 zoom transform
		readonly k: number,
		readonly tx: number,
		readonly ty: number,
		// viewport size in css pixels
		readonly width: number,
		readonly height: number,
		// device pixel ratio, used for snapping
		readonly dpr: number,
		// icons are scaled down when zoomed out past minObjectZoom
		readonly iconScale: number
	) {}

	x(wx: number): number {
		return wx * this.sx * this.k + this.tx;
	}

	y(wy: number): number {
		return wy * this.sy * this.k + this.ty;
	}

	// convert a distance in light years to css pixels
	len(ly: number): number {
		return ly * this.sx * this.k;
	}

	// project a world position to a device-pixel-snapped screen point
	anchor(position: Position | undefined): ScreenPoint {
		return {
			x: this.snap(this.x(Number(position?.x ?? 0))),
			y: this.snap(this.y(Number(position?.y ?? 0)))
		};
	}

	snap(v: number): number {
		return Math.round(v * this.dpr) / this.dpr;
	}

	// true if a point (with some radius/margin around it) is in the viewport
	visible(p: ScreenPoint, margin: number): boolean {
		return (
			p.x > -margin && p.y > -margin && p.x < this.width + margin && p.y < this.height + margin
		);
	}
}
