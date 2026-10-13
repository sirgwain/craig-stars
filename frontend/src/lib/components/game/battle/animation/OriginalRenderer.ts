// Draws the "Original" battle style (see original.ts) on the effects canvas: 1px beam lines and
// the original's 32x32 bang and torpedo icons, without smoothing.
//
// The icons are the original's own (Bang1-3, Torp1-4), loaded from /images/vcr/ like the other
// original images. Until they load, or if they aren't there, drawn stand-ins of the same shape
// and colors are used, so the style works without them.

import { VCR_ICONS, VCR_ICON_SIZE, type VcrFrame, type VcrIcon } from './original';

const ICON_PATH = '/images/vcr';
const RED = '#ff0000';
const DARK_RED = '#800000';
const YELLOW = '#ffff00';

// a burst: a filled core with eight spokes, bright red with dark red tips
function drawBang(ctx: CanvasRenderingContext2D, core: number, spoke: number) {
	const c = VCR_ICON_SIZE / 2;
	ctx.fillStyle = RED;
	ctx.beginPath();
	ctx.arc(c, c, core, 0, Math.PI * 2);
	ctx.fill();
	for (let i = 0; i < 8; i++) {
		const a = (i * Math.PI) / 4;
		const len = i % 2 === 0 ? spoke : spoke * 0.75;
		const w = i % 2 === 0 ? 2 : 1.5;
		ctx.strokeStyle = RED;
		ctx.lineWidth = w;
		ctx.beginPath();
		ctx.moveTo(c, c);
		ctx.lineTo(c + Math.cos(a) * len, c + Math.sin(a) * len);
		ctx.stroke();
		ctx.fillStyle = DARK_RED;
		ctx.fillRect(
			Math.round(c + Math.cos(a) * len) - 1,
			Math.round(c + Math.sin(a) * len) - 1,
			2,
			2
		);
	}
}

// a small spinning four-pointed star with a yellow spark
function drawTorpedo(ctx: CanvasRenderingContext2D, frame: number) {
	const c = VCR_ICON_SIZE / 2;
	ctx.save();
	ctx.translate(c, c);
	ctx.rotate((frame * Math.PI) / 8);
	ctx.fillStyle = RED;
	ctx.beginPath();
	ctx.arc(0, 0, 3.5, 0, Math.PI * 2);
	ctx.fill();
	ctx.strokeStyle = RED;
	ctx.lineWidth = 1.5;
	ctx.beginPath();
	ctx.moveTo(-6, 0);
	ctx.lineTo(6, 0);
	ctx.moveTo(0, -6);
	ctx.lineTo(0, 6);
	ctx.stroke();
	ctx.fillStyle = YELLOW;
	ctx.fillRect(-1.5, -1.5, 2, 2);
	ctx.restore();
}

function drawFallback(icon: VcrIcon): HTMLCanvasElement | undefined {
	if (typeof document === 'undefined') return undefined;
	const canvas = document.createElement('canvas');
	canvas.width = canvas.height = VCR_ICON_SIZE;
	const ctx = canvas.getContext('2d');
	if (!ctx) return undefined;
	switch (icon) {
		case 'bang1':
			drawBang(ctx, 3.5, 8);
			break;
		case 'bang2':
			drawBang(ctx, 5.5, 11);
			break;
		case 'bang3':
			drawBang(ctx, 8, 15.5);
			break;
		default:
			drawTorpedo(ctx, VCR_ICONS.indexOf(icon) - 3);
	}
	return canvas;
}

export class OriginalRenderer {
	private images = new Map<VcrIcon, HTMLImageElement>();
	private fallbacks = new Map<VcrIcon, HTMLCanvasElement | undefined>();
	private loading = false;

	// called when an icon has loaded, to redraw
	constructor(private onLoad: () => void = () => {}) {}

	// start loading the original icons (once)
	load() {
		if (this.loading || typeof Image === 'undefined') return;
		this.loading = true;
		for (const icon of VCR_ICONS) {
			const img = new Image();
			img.onload = () => {
				this.images.set(icon, img);
				this.onLoad();
			};
			img.src = `${ICON_PATH}/${icon}.png`;
		}
	}

	// the original icon if it loaded, else the drawn stand-in
	private icon(icon: VcrIcon): CanvasImageSource | undefined {
		const img = this.images.get(icon);
		if (img) return img;
		if (!this.fallbacks.has(icon)) {
			this.fallbacks.set(icon, drawFallback(icon));
		}
		return this.fallbacks.get(icon);
	}

	render(ctx: CanvasRenderingContext2D, frame: VcrFrame) {
		ctx.save();
		ctx.globalCompositeOperation = 'source-over';
		ctx.imageSmoothingEnabled = false;
		// 1px lines on the pixel grid
		ctx.lineWidth = 1;
		ctx.lineCap = 'butt';
		for (const line of frame.lines) {
			ctx.strokeStyle = line.color;
			ctx.beginPath();
			ctx.moveTo(Math.floor(line.from.x) + 0.5, Math.floor(line.from.y) + 0.5);
			ctx.lineTo(Math.floor(line.to.x) + 0.5, Math.floor(line.to.y) + 0.5);
			ctx.stroke();
		}
		const half = VCR_ICON_SIZE / 2;
		const stamps = frame.torpedo ? [...frame.stamps, frame.torpedo] : frame.stamps;
		for (const stamp of stamps) {
			const image = this.icon(stamp.icon);
			if (image) {
				ctx.drawImage(
					image,
					Math.floor(stamp.at.x) - half,
					Math.floor(stamp.at.y) - half,
					VCR_ICON_SIZE,
					VCR_ICON_SIZE
				);
			}
		}
		ctx.restore();
	}
}
