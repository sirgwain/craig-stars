// Draws the "Original" battle style (see original.ts) on the effects canvas: 1px beam lines and
// the 32x32 bang and torpedo icons, without smoothing.

import bang1 from '#lib/components/icons/BattleBang1.svg';
import bang2 from '#lib/components/icons/BattleBang2.svg';
import bang3 from '#lib/components/icons/BattleBang3.svg';
import torpedo from '#lib/components/icons/BattleTorpedo.svg';
import torpedoLight from '#lib/components/icons/BattleTorpedoLight.svg';
import { VCR_ICON_SIZE, type VcrFrame, type VcrIcon } from './original';

// the torpedo is one star, turned for each of its frames, with a highlight that stays put
const images = { bang1, bang2, bang3, torpedo, torpedoLight };
type ImageName = keyof typeof images;
const TORPEDO_FRAMES: VcrIcon[] = ['torp1', 'torp2', 'torp3', 'torp4'];
// a sixteenth of a turn per frame, so the four pointed star loops over its four frames
const TORPEDO_FRAME_TURN = Math.PI / 8;

export class OriginalRenderer {
	private images = new Map<ImageName, HTMLImageElement>();
	private loading = false;

	// called when an icon has loaded, to redraw
	constructor(private onLoad: () => void = () => {}) {}

	// start loading the icons (once)
	load() {
		if (this.loading || typeof Image === 'undefined') return;
		this.loading = true;
		for (const name of Object.keys(images) as ImageName[]) {
			const img = new Image();
			img.onload = () => {
				this.images.set(name, img);
				this.onLoad();
			};
			img.src = images[name];
		}
	}

	private drawIcon(ctx: CanvasRenderingContext2D, name: ImageName, turn = 0) {
		const image = this.images.get(name);
		if (!image) return;
		const half = VCR_ICON_SIZE / 2;
		ctx.save();
		ctx.rotate(turn);
		ctx.drawImage(image, -half, -half, VCR_ICON_SIZE, VCR_ICON_SIZE);
		ctx.restore();
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
		const stamps = frame.torpedo ? [...frame.stamps, frame.torpedo] : frame.stamps;
		for (const stamp of stamps) {
			ctx.save();
			ctx.translate(Math.floor(stamp.at.x), Math.floor(stamp.at.y));
			const torpedoFrame = TORPEDO_FRAMES.indexOf(stamp.icon);
			if (torpedoFrame === -1) {
				this.drawIcon(ctx, stamp.icon as ImageName);
			} else {
				this.drawIcon(ctx, 'torpedo', torpedoFrame * TORPEDO_FRAME_TURN);
				this.drawIcon(ctx, 'torpedoLight');
			}
			ctx.restore();
		}
		ctx.restore();
	}
}
