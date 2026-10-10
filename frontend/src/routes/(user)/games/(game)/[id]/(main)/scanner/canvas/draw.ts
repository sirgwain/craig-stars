/**
 * Draw the scanner onto a canvas. Everything is drawn in screen space: each object is projected
 * to a pixel-snapped anchor and its icon, name, fleet count and selection arrow are drawn at
 * fixed pixel offsets from that anchor, like the original Stars! scanner.
 */
import type { Universe } from '#lib/services/Universe.js';
import { StargateWarpSpeed } from '#lib/types/Consts.js';
import { MapObjectType, type Fleet, type Minefield, type Waypoint } from '#lib/types/cs-proto.js';
import {
	equal,
	type MapObjectLike,
	type MovingMapObject,
	type Position
} from '#lib/types/MapObject.js';
import type { CommandedPlayer } from '#lib/types/Player.js';
import type { PlayerSettings } from '#lib/types/PlayerSettings.js';
import { getDisplayColor, getStrokeColor } from '#lib/utils/colorUtils.js';
import {
	headingAngle,
	mineralBarMax,
	packetColor,
	type DestLine,
	type FleetDraw,
	type MinefieldDraw,
	type PlanetDraw,
	type ScannerCircle,
	type WaypointPath,
	type WormholeLink
} from './scene';
import type { ScannerView, ScreenPoint } from './view';
import { fixedScannerColors, type ScannerColors } from './colors';

const nameFont = 'bold 12px Arial, Helvetica, sans-serif';
const countFont = 'bold 10px Arial, Helvetica, sans-serif';
// gap in pixels between an icon and the text/arrow drawn above or below it
const gap = 2;
const fleetSize = 8;
// the selection arrow is a chevron this many pixels wide, and half as tall
const selectionArrowWidth = 10;
// year ticks are only drawn if they are at least this many pixels apart
const minTickSpacing = 6;
const tickLength = 3;
// waypoint paths stop at least this many pixels short of a waypoint
const waypointClearanceMin = 5;

export type ScannerFrame = {
	colors: ScannerColors;
	universe: Universe;
	player: CommandedPlayer;
	settings: PlayerSettings;
	scanners: ScannerCircle[];
	minefields: MinefieldDraw[];
	packetDests: DestLine[];
	routeDests: DestLine[];
	waypointPaths: WaypointPath[];
	planets: PlanetDraw[];
	planetsByNum: Map<number, PlanetDraw>;
	fleets: FleetDraw[];
	wormholeLinks: WormholeLink[];
	selectedMapObject: MapObjectLike | undefined;
	highlightedMapObject: MapObjectLike | undefined;
	commandedMapObject: MapObjectLike | undefined;
	selectedWaypoint: Waypoint | undefined;
	// the map object to show the locator X on, and how opaque it is
	locator: { mapObject: MapObjectLike; alpha: number } | undefined;
};

export function drawScanner(ctx: CanvasRenderingContext2D, view: ScannerView, frame: ScannerFrame) {
	const { colors } = frame;
	ctx.setTransform(view.dpr, 0, 0, view.dpr, 0, 0);
	ctx.fillStyle = colors.background;
	ctx.fillRect(0, 0, view.width, view.height);

	drawScanners(ctx, view, frame);
	if (frame.settings.showMinefields) {
		drawMinefields(ctx, view, frame);
	}
	drawDestLines(ctx, view, frame.packetDests, colors.packetDestLine);
	drawDestLines(ctx, view, frame.routeDests, colors.routeLine);
	drawWaypointPaths(ctx, view, frame);
	drawPlanets(ctx, view, frame);
	drawMineralPackets(ctx, view, frame);
	drawWormholes(ctx, view, frame);
	drawFleets(ctx, view, frame);
	drawMysteryTraders(ctx, view, frame);
	drawWarpLine(ctx, view, frame);
	drawWormholeLinks(ctx, view, frame);
	drawSalvages(ctx, view, frame);
	drawNames(ctx, view, frame);
	drawSelection(ctx, view, frame);
	drawLocator(ctx, view, frame);
}

function drawScanners(ctx: CanvasRenderingContext2D, view: ScannerView, frame: ScannerFrame) {
	const { colors } = frame;
	const scannerScale = frame.settings.scannerPercent / 100.0;

	// draw all the regular scanners as one path, then all the penetrating scanners on top
	for (const pen of [false, true]) {
		ctx.beginPath();
		for (const scanner of frame.scanners) {
			const range = pen ? scanner.scanRangePen : scanner.scanRange;
			if (range <= 0) {
				continue;
			}
			const p = view.anchor(scanner.position);
			const r = view.len(range * scannerScale);
			if (!view.visible(p, r)) {
				continue;
			}
			ctx.moveTo(p.x + r, p.y);
			ctx.arc(p.x, p.y, r, 0, Math.PI * 2);
		}
		ctx.fillStyle = pen ? colors.scannerPen : colors.scanner;
		ctx.fill();
	}
}

// minefields are filled with a dot pattern, one pattern per color per frame
const minefieldDots = [
	[4, 0],
	[2, 2],
	[6, 2],
	[0, 4],
	[2, 6],
	[6, 6]
];

function minefieldPattern(
	ctx: CanvasRenderingContext2D,
	view: ScannerView,
	color: string
): CanvasPattern | null {
	// Don't retain canvas bitmaps between frames: Chrome can discard them while a phone
	// sleeps. Restoring the main canvas doesn't restore the pixels of cached pattern tiles.
	const tile = document.createElement('canvas');
	tile.width = 8;
	tile.height = 8;
	const tileCtx = tile.getContext('2d');
	if (!tileCtx) {
		return null;
	}
	tileCtx.fillStyle = color;
	minefieldDots.forEach(([x, y]) => tileCtx.fillRect(x, y, 1, 1));
	const pattern = ctx.createPattern(tile, 'repeat');
	// anchor the pattern to the map so it doesn't swim while panning. The context is already
	// scaled by the dpr, so one tile pixel is one css pixel
	pattern?.setTransform(new DOMMatrix().translate(view.tx % 8, view.ty % 8));
	return pattern;
}

function drawMinefields(ctx: CanvasRenderingContext2D, view: ScannerView, frame: ScannerFrame) {
	const selected =
		frame.selectedMapObject?.mapObject?.type === MapObjectType.MINEFIELD
			? (frame.selectedMapObject as Minefield)
			: undefined;

	// batch minefields by color
	const byColor = new Map<string, MinefieldDraw[]>();
	let selectedDraw: MinefieldDraw | undefined;
	for (const mf of frame.minefields) {
		if (selected && equal(selected, mf.minefield)) {
			selectedDraw = mf;
			continue;
		}
		const list = byColor.get(mf.color) ?? [];
		list.push(mf);
		byColor.set(mf.color, list);
	}

	for (const [color, minefields] of byColor) {
		ctx.beginPath();
		for (const mf of minefields) {
			const p = view.anchor(mf.minefield.mapObject?.position);
			const r = view.len(mf.radius);
			if (view.visible(p, r)) {
				ctx.moveTo(p.x + r, p.y);
				ctx.arc(p.x, p.y, r, 0, Math.PI * 2);
			}
		}
		ctx.fillStyle = minefieldPattern(ctx, view, color) ?? color;
		ctx.fill();
	}

	// draw the selected minefield last, darkened, with a marker in the center
	if (selectedDraw) {
		const p = view.anchor(selectedDraw.minefield.mapObject?.position);
		const r = view.len(selectedDraw.radius);
		ctx.beginPath();
		ctx.arc(p.x, p.y, r, 0, Math.PI * 2);
		ctx.fillStyle =
			minefieldPattern(ctx, view, darken(selectedDraw.color, 0.6)) ?? selectedDraw.color;
		ctx.fill();

		const size = Math.max(view.len(2), 2);
		ctx.fillStyle = selectedDraw.color;
		ctx.fillRect(p.x - size / 2, p.y - size / 2, size, size);
	}
}

function darken(color: string, amount: number): string {
	let hex = color.replace('#', '');
	if (hex.length === 3) {
		hex = hex
			.split('')
			.map((c) => c + c)
			.join('');
	}
	const channel = (i: number) =>
		Math.round(parseInt(hex.substring(i, i + 2), 16) * amount)
			.toString(16)
			.padStart(2, '0');
	return `#${channel(0)}${channel(2)}${channel(4)}`;
}

// packet and route destination lines are dashed, one dash per year of travel, with an arrow
function drawDestLines(
	ctx: CanvasRenderingContext2D,
	view: ScannerView,
	lines: DestLine[],
	color: string
) {
	ctx.strokeStyle = color;
	ctx.fillStyle = color;
	ctx.lineCap = 'square';
	for (const line of lines) {
		const from = view.anchor(line.from);
		const to = view.anchor(line.to);
		const dx = to.x - from.x;
		const dy = to.y - from.y;
		const length = Math.hypot(dx, dy);
		if (length === 0) {
			continue;
		}

		const yearLength = view.len(line.yearDist);
		const gapLength = view.len(5);
		if (yearLength > gapLength) {
			ctx.setLineDash([yearLength - gapLength, gapLength]);
			ctx.lineDashOffset = view.len(line.yearDist / 2) - gapLength;
		}
		ctx.lineWidth = line.width;
		ctx.beginPath();
		ctx.moveTo(from.x, from.y);
		ctx.lineTo(to.x, to.y);
		ctx.stroke();
		ctx.setLineDash([]);
		ctx.lineDashOffset = 0;

		// arrow head just short of the target
		const ux = dx / length;
		const uy = dy / length;
		const tip = { x: to.x - ux * 4, y: to.y - uy * 4 };
		const arrowLength = 6 * line.width;
		const arrowWidth = 3 * line.width;
		ctx.beginPath();
		ctx.moveTo(tip.x, tip.y);
		ctx.lineTo(
			tip.x - ux * arrowLength - uy * arrowWidth,
			tip.y - uy * arrowLength + ux * arrowWidth
		);
		ctx.lineTo(
			tip.x - ux * arrowLength + uy * arrowWidth,
			tip.y - uy * arrowLength - ux * arrowWidth
		);
		ctx.closePath();
		ctx.fill();
	}
	ctx.lineCap = 'butt';
}

function drawWaypointPaths(ctx: CanvasRenderingContext2D, view: ScannerView, frame: ScannerFrame) {
	const { colors } = frame;
	ctx.lineCap = 'butt';
	for (const path of frame.waypointPaths) {
		ctx.strokeStyle = path.commanded ? colors.waypointLineCommanded : colors.waypointLine;
		const points = path.waypoints.map((wp) => view.anchor(wp.position));
		const exits = path.exits.map((exit) => (exit ? view.anchor(exit) : undefined));
		// like Stars!, leave a gap around each waypoint so the path doesn't cover what's there
		const clearances = path.waypoints.map((wp) => waypointClearance(frame, view, wp));
		const exitClearance = wormholeClearance(view);

		for (let i = 1; i < path.waypoints.length; i++) {
			const wp0 = path.waypoints[i - 1];
			const wp1 = path.waypoints[i];
			const width = frame.selectedWaypoint === wp0 ? 3 : path.commanded ? 2 : 1;

			// a leg after a wormhole starts where the fleet comes out
			const exit = exits[i - 1];
			const from = exit ?? points[i - 1];
			const segment = segmentBetween(
				from,
				points[i],
				(exit ? exitClearance : clearances[i - 1]) + width / 2,
				clearances[i] + width / 2
			);
			if (!segment) {
				continue;
			}
			ctx.lineWidth = width;
			ctx.beginPath();
			ctx.moveTo(segment.x0, segment.y0);
			ctx.lineTo(segment.x1, segment.y1);
			ctx.stroke();

			if (path.commanded && wp1.warpSpeed !== StargateWarpSpeed) {
				drawYearTicks(
					ctx,
					view,
					path.exits[i - 1] ?? wp0.position,
					wp1.position,
					from,
					points[i],
					wp1.warpSpeed,
					width,
					segment.start,
					segment.end,
					returnLegOverlaps(path, i)
				);
			}
		}

		// show the jump through each wormhole as a dashed line from the entrance to the exit
		ctx.lineWidth = path.commanded ? 2 : 1;
		ctx.setLineDash([4, 3]);
		points.forEach((entrance, i) => {
			const exit = exits[i];
			const segment = exit && segmentBetween(entrance, exit, exitClearance, exitClearance);
			if (segment) {
				ctx.beginPath();
				ctx.moveTo(segment.x0, segment.y0);
				ctx.lineTo(segment.x1, segment.y1);
				ctx.stroke();
			}
		});
		ctx.setLineDash([]);
	}
}

// Find portions of a leg that retrace an earlier leg in the opposite direction. Use world
// coordinates so overlap detection is independent of zoom and snapping, including wormhole exits.
function returnLegOverlaps(path: WaypointPath, legIndex: number) {
	const from = path.exits[legIndex - 1] ?? path.waypoints[legIndex - 1].position;
	const to = path.waypoints[legIndex].position;
	const x = Number(from?.x ?? 0);
	const y = Number(from?.y ?? 0);
	const dx = Number(to?.x ?? 0) - x;
	const dy = Number(to?.y ?? 0) - y;
	const lengthSquared = dx * dx + dy * dy;
	const overlaps: { start: number; end: number }[] = [];
	if (lengthSquared === 0) {
		return overlaps;
	}
	for (let i = 1; i < legIndex; i++) {
		const otherFrom = path.exits[i - 1] ?? path.waypoints[i - 1].position;
		const otherTo = path.waypoints[i].position;
		const ox = Number(otherFrom?.x ?? 0) - x;
		const oy = Number(otherFrom?.y ?? 0) - y;
		const odx = Number(otherTo?.x ?? 0) - Number(otherFrom?.x ?? 0);
		const ody = Number(otherTo?.y ?? 0) - Number(otherFrom?.y ?? 0);
		// A crossing or a parallel leg elsewhere is not a return along this route.
		if (dx * odx + dy * ody >= 0 || dx * ody !== dy * odx || dx * oy !== dy * ox) {
			continue;
		}
		const start = Math.max(0, ((ox + odx) * dx + (oy + ody) * dy) / lengthSquared);
		const end = Math.min(1, (ox * dx + oy * dy) / lengthSquared);
		if (start < end) overlaps.push({ start, end });
	}
	return overlaps;
}

// the part of a line from p0 to p1 that leaves a gap at each end, or undefined if nothing's left
function segmentBetween(p0: ScreenPoint, p1: ScreenPoint, startGap: number, endGap: number) {
	const length = Math.hypot(p1.x - p0.x, p1.y - p0.y);
	const start = startGap;
	const end = length - endGap;
	if (end <= start) {
		return undefined;
	}
	const ux = (p1.x - p0.x) / length;
	const uy = (p1.y - p0.y) / length;
	return {
		x0: p0.x + ux * start,
		y0: p0.y + uy * start,
		x1: p0.x + ux * end,
		y1: p0.y + uy * end,
		start,
		end
	};
}

// paths stop just outside wormhole icons
function wormholeClearance(view: ScannerView) {
	return Math.max(wormholeSize / 2, waypointClearanceMin) * view.iconScale;
}

// how far from a waypoint its path should stop, in pixels
function waypointClearance(frame: ScannerFrame, view: ScannerView, wp: Waypoint): number {
	const min = waypointClearanceMin * view.iconScale;
	const target = wp.mapObjectTarget;
	if (target?.targetType === MapObjectType.WORMHOLE) {
		return wormholeClearance(view);
	}
	if (target?.targetType === MapObjectType.PLANET) {
		const planet = frame.planetsByNum.get(target.targetNum);
		if (planet) {
			return Math.max(planetExtents(planet, view).bottom, min);
		}
	}
	return min;
}

/**
 * Draw a tick across the line for each year of travel.
 * Ticks are placed at exact multiples of the distance travelled per year from the start of
 * the leg, and skipped if they would be too close together to read. Only ticks on the drawn
 * part of the line (between start and end pixels along it) are shown, excluding portions
 * that retrace an earlier leg.
 */
function drawYearTicks(
	ctx: CanvasRenderingContext2D,
	view: ScannerView,
	from: Position | undefined,
	to: Position | undefined,
	p0: ScreenPoint,
	p1: ScreenPoint,
	warpSpeed: number,
	lineWidth: number,
	start: number,
	end: number,
	returnOverlaps: { start: number; end: number }[]
) {
	if (warpSpeed < 1 || warpSpeed > 10) {
		return;
	}
	const yearDist = warpSpeed * warpSpeed;
	const legDist = Math.hypot(
		Number(to?.x ?? 0) - Number(from?.x ?? 0),
		Number(to?.y ?? 0) - Number(from?.y ?? 0)
	);
	const dx = p1.x - p0.x;
	const dy = p1.y - p0.y;
	const screenDist = Math.hypot(dx, dy);
	if (legDist <= yearDist || (screenDist * yearDist) / legDist < minTickSpacing) {
		return;
	}

	// unit vector across the leg
	const nx = -dy / screenDist;
	const ny = dx / screenDist;
	const inner = lineWidth / 2;
	const outer = inner + tickLength;

	// Outline white ticks so they remain readable over space, scanner coverage, and paths.
	ctx.save();
	ctx.globalCompositeOperation = 'source-over';
	ctx.beginPath();
	for (let i = 1; i * yearDist < legDist; i++) {
		const t = (i * yearDist) / legDist;
		if (
			t * screenDist < start ||
			t * screenDist > end ||
			returnOverlaps.some((overlap) => t >= overlap.start && t <= overlap.end)
		) {
			continue;
		}
		const x = view.snap(p0.x + dx * t);
		const y = view.snap(p0.y + dy * t);
		for (const side of [-1, 1]) {
			ctx.moveTo(x + side * nx * inner, y + side * ny * inner);
			ctx.lineTo(x + side * nx * outer, y + side * ny * outer);
		}
	}
	ctx.strokeStyle = '#000000';
	ctx.lineWidth = 3;
	ctx.stroke();
	ctx.strokeStyle = fixedScannerColors.yearTick;
	ctx.lineWidth = 1;
	ctx.stroke();
	ctx.restore();
}

function drawPlanets(ctx: CanvasRenderingContext2D, view: ScannerView, frame: ScannerFrame) {
	const { colors } = frame;
	const s = view.iconScale;
	for (const draw of frame.planets) {
		const p = view.anchor(draw.planet.mapObject?.position);
		if (!view.visible(p, 40 * s)) {
			continue;
		}

		if (draw.value) {
			drawPlanetValue(ctx, p, s, draw.value);
		}
		if (draw.pop) {
			drawPopulation(ctx, view, p, draw.pop);
		}
		if (draw.minerals) {
			drawMineralGraph(ctx, view, p, draw.minerals, colors);
		}
		if (draw.normal) {
			drawNormalPlanet(ctx, p, s, draw.normal, colors);
		}
		if (draw.count) {
			// fleet counts go above the planet, like Stars!
			drawCount(ctx, p.x, p.y + planetExtents(draw, view).top - gap, draw.count, 'bottom');
		}
	}
}

function drawNormalPlanet(
	ctx: CanvasRenderingContext2D,
	p: ScreenPoint,
	s: number,
	normal: NonNullable<PlanetDraw['normal']>,
	colors: ScannerColors
) {
	if (normal.ring) {
		ctx.beginPath();
		ctx.arc(p.x, p.y, normal.ring.radius * s, 0, Math.PI * 2);
		ctx.lineWidth = normal.ring.width * s;
		ctx.strokeStyle = normal.ring.color;
		ctx.setLineDash(normal.ring.dashed ? [10 * s, 6 * s] : []);
		ctx.stroke();
		ctx.setLineDash([]);
	}

	ctx.beginPath();
	ctx.arc(p.x, p.y, normal.radius * s, 0, Math.PI * 2);
	ctx.fillStyle = normal.fill;
	ctx.fill();
	ctx.lineWidth = normal.strokeWidth * s;
	ctx.strokeStyle = colors.planet.outline;
	ctx.stroke();

	const w = normal.starbaseWidth * s;
	const xOffset = normal.starbaseXOffset * s;
	const yOffset = normal.starbaseYOffset * s;
	if (normal.starbase) {
		drawMarker(ctx, p.x + xOffset, p.y - yOffset, w, normal.starbase, s);
	}
	if (normal.stargate) {
		drawMarker(ctx, p.x - xOffset - w, p.y - yOffset, w, colors.stargate, s);
	}
	if (normal.massDriver) {
		drawMarker(ctx, p.x - w / 2, p.y - yOffset - w / 2, w, colors.massDriver, s);
	}
}

function drawMarker(
	ctx: CanvasRenderingContext2D,
	x: number,
	y: number,
	w: number,
	style: { fill: string; stroke: string },
	s: number
) {
	ctx.beginPath();
	ctx.roundRect(x, y, w, w, 0.5 * s);
	ctx.fillStyle = style.fill;
	ctx.fill();
	ctx.lineWidth = 0.5 * s;
	ctx.strokeStyle = style.stroke;
	ctx.stroke();
}

function drawPlanetValue(
	ctx: CanvasRenderingContext2D,
	p: ScreenPoint,
	s: number,
	value: NonNullable<PlanetDraw['value']>
) {
	// Stars! ellipses cover 2r+1 pixels, so add half a pixel to each radius
	for (const [radius, color] of [
		[value.outerRadius, value.outer],
		[value.innerRadius, value.inner]
	] as const) {
		ctx.beginPath();
		ctx.arc(p.x, p.y, (radius + 0.5) * s, 0, Math.PI * 2);
		ctx.fillStyle = color;
		ctx.fill();
	}

	if (value.flag) {
		// a 7x6 flag on a 21px pole, over a black backdrop so it shows up on scanner ranges
		ctx.fillStyle = '#000';
		ctx.fillRect(p.x - s, p.y - 20 * s, 9 * s, 8 * s);
		ctx.fillStyle = value.flag;
		ctx.fillRect(p.x, p.y - 19 * s, s, 21 * s);
		ctx.fillRect(p.x, p.y - 19 * s, 7 * s, 6 * s);
	}
}

/**
 * Stars! zoom levels are pixels per light year (100% is 1 pixel per light year). Like Stars!,
 * the population view draws smaller circles below 200% and the mineral views draw smaller
 * graphs below 100%.
 */
function smallPopulation(view: ScannerView) {
	return view.len(1) < 2;
}

function populationRadius(view: ScannerView, pop: NonNullable<PlanetDraw['pop']>) {
	return smallPopulation(view) ? (pop.radius + 1) >> 1 : pop.radius;
}

function drawPopulation(
	ctx: CanvasRenderingContext2D,
	view: ScannerView,
	p: ScreenPoint,
	pop: NonNullable<PlanetDraw['pop']>
) {
	// Stars! ellipses cover 2r+1 pixels with a 1 pixel outline
	const r = populationRadius(view, pop);
	ctx.beginPath();
	ctx.arc(p.x, p.y, r + 0.5, 0, Math.PI * 2);
	ctx.fillStyle = pop.fill;
	ctx.fill();
	ctx.beginPath();
	ctx.arc(p.x, p.y, r, 0, Math.PI * 2);
	ctx.lineWidth = 1;
	ctx.strokeStyle = pop.ring;
	ctx.stroke();
}

// the mineral graph layout from vrgScanPO in Stars!, relative to the planet, in pixels
const mineralGraphs = {
	large: { left: 7, bottom: 12, axis: 19, barWidth: 4, barSpacing: 6, barScale: 1 },
	small: { left: 3, bottom: 10, axis: 11, barWidth: 2, barSpacing: 3, barScale: 0.5 }
};

function mineralGraph(view: ScannerView) {
	return view.len(1) < 1 ? mineralGraphs.small : mineralGraphs.large;
}

// a small graph of ironium, boranium and germanium bars up and to the left of the planet
function drawMineralGraph(
	ctx: CanvasRenderingContext2D,
	view: ScannerView,
	p: ScreenPoint,
	minerals: NonNullable<PlanetDraw['minerals']>,
	colors: ScannerColors
) {
	const graph = mineralGraph(view);
	const left = p.x - graph.left;
	const bottom = p.y - graph.bottom;

	// axes
	ctx.fillStyle = colors.mineralAxis;
	ctx.fillRect(left - 2, bottom, graph.axis, 1);
	ctx.fillRect(left - 2, bottom - graph.axis + 1, 1, graph.axis);

	const barColors = [colors.ironium, colors.boranium, colors.germanium];
	minerals.forEach((level, i) => {
		const height = Math.trunc(Math.min(level, mineralBarMax) * graph.barScale);
		if (height > 0) {
			ctx.fillStyle = barColors[i];
			ctx.fillRect(left + i * graph.barSpacing, bottom - height, graph.barWidth, height);
		}
	});
}

/**
 * How far a planet's drawing extends above and below its center, in pixels, so names, fleet
 * counts, the selection arrow and waypoint paths can be placed around it.
 */
function planetExtents(draw: PlanetDraw, view: ScannerView) {
	let top = draw.top * view.iconScale;
	let bottom = draw.bottom * view.iconScale;
	if (draw.pop) {
		const r = populationRadius(view, draw.pop) + 0.5;
		top = Math.min(top, -r);
		bottom = Math.max(bottom, r);
	}
	if (draw.minerals) {
		const graph = mineralGraph(view);
		top = Math.min(top, -(graph.bottom + graph.axis - 1));
	}
	return { top, bottom };
}

function drawCount(
	ctx: CanvasRenderingContext2D,
	x: number,
	y: number,
	count: { value: number; color: string },
	baseline: CanvasTextBaseline
) {
	ctx.font = countFont;
	ctx.textAlign = 'center';
	ctx.textBaseline = baseline;
	ctx.fillStyle = count.color;
	ctx.fillText(`${count.value}`, x, y);
}

function drawMineralPackets(ctx: CanvasRenderingContext2D, view: ScannerView, frame: ScannerFrame) {
	const s = view.iconScale;
	const size = 10 * s;
	ctx.lineWidth = 2 * s;
	for (const packet of frame.universe.mineralPackets) {
		const p = view.anchor(packet.mapObject?.position);
		if (!view.visible(p, size)) {
			continue;
		}
		ctx.strokeStyle = packetColor(
			packet.mapObject?.playerNum,
			frame.player,
			frame.universe,
			frame.settings,
			frame.colors
		);
		ctx.strokeRect(p.x - size / 2, p.y - size / 2, size, size);
	}
}

// wormhole icon from game-icons.net, drawn in a 512x512 box
let wormholeOuter: Path2D | undefined;
let wormholeInner: Path2D | undefined;
const wormholeSize = 12;

function drawWormholes(ctx: CanvasRenderingContext2D, view: ScannerView, frame: ScannerFrame) {
	const { colors } = frame;
	if (frame.universe.wormholes.length === 0) {
		return;
	}
	wormholeOuter ??= new Path2D(wormholeOuterPath);
	wormholeInner ??= new Path2D(wormholeInnerPath);

	const size = wormholeSize * view.iconScale;
	const scale = size / 512;
	for (const wormhole of frame.universe.wormholes) {
		const p = view.anchor(wormhole.mapObject?.position);
		if (!view.visible(p, size)) {
			continue;
		}
		ctx.save();
		ctx.translate(p.x - size / 2, p.y - size / 2);
		ctx.scale(scale, scale);
		ctx.fillStyle = colors.wormhole;
		ctx.fill(wormholeOuter);
		ctx.fillStyle = colors.background;
		ctx.fill(wormholeInner);
		ctx.restore();
	}
}

// draw the fleet triangle centered on p, pointing along the angle
function trianglePath(ctx: CanvasRenderingContext2D, p: ScreenPoint, s: number, angle: number) {
	const cos = Math.cos(angle);
	const sin = Math.sin(angle);
	const half = fleetSize / 2;
	ctx.beginPath();
	for (const [i, [x, y]] of [
		[0, 0],
		[0, fleetSize],
		[fleetSize, fleetSize]
	].entries()) {
		const lx = (x - half) * s;
		const ly = (y - half) * s;
		const sx = p.x + lx * cos - ly * sin;
		const sy = p.y + lx * sin + ly * cos;
		if (i === 0) {
			ctx.moveTo(sx, sy);
		} else {
			ctx.lineTo(sx, sy);
		}
	}
	ctx.closePath();
}

function drawFleets(ctx: CanvasRenderingContext2D, view: ScannerView, frame: ScannerFrame) {
	const { colors } = frame;
	const s = view.iconScale;
	ctx.lineJoin = 'round';
	for (const fleet of frame.fleets) {
		const p = view.anchor(fleet.fleet.mapObject?.position);
		if (!view.visible(p, 20)) {
			continue;
		}
		const color = fleet.commanded ? colors.commandedFleet : fleet.color;
		trianglePath(ctx, p, s, fleet.angle);
		ctx.lineWidth = s;
		ctx.strokeStyle = getStrokeColor(color);
		ctx.stroke();
		ctx.fillStyle = color;
		ctx.fill();

		if (fleet.count) {
			// token counts go just above the fleet, like Stars!
			drawCount(ctx, p.x, p.y - fleetExtent * s - gap, fleet.count, 'bottom');
		}
	}
}

function drawMysteryTraders(ctx: CanvasRenderingContext2D, view: ScannerView, frame: ScannerFrame) {
	const { colors } = frame;
	const s = view.iconScale;
	ctx.fillStyle = colors.mysteryTrader;
	for (const mt of frame.universe.mysteryTraders) {
		const p = view.anchor(mt.mapObject?.position);
		if (!view.visible(p, 20)) {
			continue;
		}
		trianglePath(ctx, p, s, headingAngle(mt.heading));
		ctx.fill();
	}
}

// show where a foreign fleet, packet or mystery trader is going for the next 5 years
function drawWarpLine(ctx: CanvasRenderingContext2D, view: ScannerView, frame: ScannerFrame) {
	const { colors } = frame;
	const mo = frame.selectedMapObject;
	const type = mo?.mapObject?.type;
	if (
		!mo?.mapObject ||
		!(
			type === MapObjectType.MINERAL_PACKET ||
			type === MapObjectType.MYSTERY_TRADER ||
			(type === MapObjectType.FLEET && mo.mapObject.playerNum != frame.player.num)
		)
	) {
		return;
	}

	const moving = mo as MovingMapObject;
	if (!moving.warpSpeed || !moving.heading) {
		return;
	}

	let color = colors.warpLine;
	if (mo.mapObject.playerNum) {
		color = getDisplayColor(mo.mapObject.playerNum, frame.player, frame.universe, frame.settings);
	} else if (type === MapObjectType.MYSTERY_TRADER) {
		color = colors.mysteryTrader;
	}

	const yearDist = moving.warpSpeed * moving.warpSpeed;
	const heading = moving.heading;
	const origin = mo.mapObject.position;
	const points = [-5, -4, -3, -2, -1, 1, 2, 3, 4, 5].map((year) =>
		view.anchor({
			x: Number(origin?.x ?? 0) + heading.x * Math.ceil(yearDist * year),
			y: Number(origin?.y ?? 0) + heading.y * Math.ceil(yearDist * year)
		})
	);

	ctx.strokeStyle = color;
	ctx.lineWidth = 1;
	ctx.beginPath();
	points.forEach((p, i) => (i === 0 ? ctx.moveTo(p.x, p.y) : ctx.lineTo(p.x, p.y)));
	ctx.stroke();

	// chevrons pointing in the direction of travel at each year
	const length = Math.hypot(heading.x, heading.y) || 1;
	const ux = heading.x / length;
	const uy = heading.y / length;
	const size = 4;
	ctx.lineWidth = 2;
	ctx.beginPath();
	for (const p of points.slice(1, -1)) {
		ctx.moveTo(p.x - ux * size - uy * size, p.y - uy * size + ux * size);
		ctx.lineTo(p.x, p.y);
		ctx.lineTo(p.x - ux * size + uy * size, p.y - uy * size - ux * size);
	}
	ctx.stroke();
}

function drawWormholeLinks(ctx: CanvasRenderingContext2D, view: ScannerView, frame: ScannerFrame) {
	const { colors } = frame;
	ctx.strokeStyle = colors.wormholeLink;
	ctx.lineWidth = 1;
	ctx.lineCap = 'round';
	ctx.beginPath();
	for (const link of frame.wormholeLinks) {
		const from = view.anchor(link.from);
		const to = view.anchor(link.to);
		ctx.moveTo(from.x, from.y);
		ctx.lineTo(to.x, to.y);
	}
	ctx.stroke();
	ctx.lineCap = 'butt';
}

function drawSalvages(ctx: CanvasRenderingContext2D, view: ScannerView, frame: ScannerFrame) {
	const { colors } = frame;
	const s = view.iconScale;
	// a 10x10 square rotated 45º is a diamond
	const r = (10 * s) / Math.SQRT2;
	ctx.strokeStyle = colors.salvage;
	ctx.lineWidth = s;
	for (const salvage of frame.universe.salvages) {
		const p = view.anchor(salvage.mapObject?.position);
		if (!view.visible(p, r)) {
			continue;
		}
		ctx.beginPath();
		ctx.moveTo(p.x, p.y - r);
		ctx.lineTo(p.x + r, p.y);
		ctx.lineTo(p.x, p.y + r);
		ctx.lineTo(p.x - r, p.y);
		ctx.closePath();
		ctx.stroke();
	}
}

/**
 * Show planet names below each planet. If names are turned off, only show the name of the
 * planet under the mouse or the selected planet.
 */
function drawNames(ctx: CanvasRenderingContext2D, view: ScannerView, frame: ScannerFrame) {
	const s = view.iconScale;
	const selectedPlanet = selectedPlanetDraw(frame);
	ctx.font = nameFont;
	ctx.textAlign = 'center';
	ctx.textBaseline = 'top';

	const drawName = (draw: PlanetDraw) => {
		const p = view.anchor(draw.planet.mapObject?.position);
		if (!view.visible(p, 60)) {
			return;
		}
		// make room for the selection arrow under the selected planet
		const arrow = draw === selectedPlanet ? selectionArrowHeight(s) + gap : 0;
		ctx.fillStyle = draw.nameColor;
		ctx.fillText(
			draw.planet.mapObject?.name ?? '',
			p.x,
			p.y + planetExtents(draw, view).bottom + gap + arrow
		);
	};

	if (frame.settings.showNames) {
		frame.planets.forEach(drawName);
		return;
	}

	for (const mo of [frame.highlightedMapObject, frame.selectedMapObject]) {
		if (mo?.mapObject?.type === MapObjectType.PLANET) {
			const draw = frame.planetsByNum.get(mo.mapObject.num);
			if (draw) {
				drawName(draw);
			}
		}
	}
}

// the planet the selected map object is at, if any
function selectedPlanetDraw(frame: ScannerFrame): PlanetDraw | undefined {
	const mo = frame.selectedMapObject;
	if (mo?.mapObject?.type === MapObjectType.PLANET) {
		return frame.planetsByNum.get(mo.mapObject.num);
	}
	if (mo?.mapObject?.type === MapObjectType.FLEET && (mo as Fleet).orbitingPlanetNum) {
		return frame.planetsByNum.get((mo as Fleet).orbitingPlanetNum);
	}
	return undefined;
}

// how far below its anchor each type of icon extends, in icon units
const fleetExtent = 6;
// how far below its anchor the selected map object extends, in pixels
function selectedIconBottom(frame: ScannerFrame, view: ScannerView, mo: MapObjectLike): number {
	const planet = selectedPlanetDraw(frame);
	if (planet) {
		return planetExtents(planet, view).bottom;
	}
	return iconBottom(mo) * view.iconScale;
}

function iconBottom(mo: MapObjectLike): number {
	switch (mo.mapObject?.type) {
		case MapObjectType.FLEET:
		case MapObjectType.MYSTERY_TRADER:
			return fleetExtent;
		case MapObjectType.MINERAL_PACKET:
			return 6;
		case MapObjectType.WORMHOLE:
			return wormholeSize / 2;
		case MapObjectType.SALVAGE:
			return 10 / Math.SQRT2;
		default:
			return 1;
	}
}

function selectionArrowHeight(s: number) {
	return (selectionArrowWidth * Math.max(s, 0.75)) / 2;
}

// draw a chevron under the selected map object
function drawSelection(ctx: CanvasRenderingContext2D, view: ScannerView, frame: ScannerFrame) {
	const mo = frame.selectedMapObject;
	if (!mo?.mapObject) {
		return;
	}
	const s = view.iconScale;
	const p = view.anchor(mo.mapObject.position);
	const height = selectionArrowHeight(s);
	const top = p.y + selectedIconBottom(frame, view, mo) + gap;

	ctx.beginPath();
	ctx.moveTo(p.x - height, top + height);
	ctx.lineTo(p.x, top);
	ctx.lineTo(p.x + height, top + height);
	ctx.lineWidth = 1.5;
	ctx.lineCap = 'round';
	ctx.lineJoin = 'round';
	ctx.strokeStyle = frame.colors.selection;
	ctx.stroke();
	ctx.lineCap = 'butt';
}

// draw a big X through a map object so the player can find it
function drawLocator(ctx: CanvasRenderingContext2D, view: ScannerView, frame: ScannerFrame) {
	const { colors } = frame;
	if (!frame.locator?.mapObject.mapObject) {
		return;
	}
	const p = view.anchor(frame.locator.mapObject.mapObject.position);
	const size = Math.max(view.width, view.height) * 2;
	ctx.globalAlpha = frame.locator.alpha;
	ctx.strokeStyle = colors.locator;
	ctx.lineWidth = 2;
	ctx.lineCap = 'round';
	ctx.beginPath();
	ctx.moveTo(p.x - size, p.y - size);
	ctx.lineTo(p.x + size, p.y + size);
	ctx.moveTo(p.x + size, p.y - size);
	ctx.lineTo(p.x - size, p.y + size);
	ctx.stroke();
	ctx.globalAlpha = 1;
	ctx.lineCap = 'butt';
}

const wormholeOuterPath =
	'M 307.819885 30.001923 C 185.877243 30.001923 67.211006 131.409271 42.770405 256.501953 C 34.083199 300.968384 38.480839 342.436584 53.223545 377.425781 C 48.210655 371.217102 43.497414 364.741028 39.141521 358.012024 C 18.859095 326.67804 8.898401 294.657959 8.440338 264.720398 C -5.976942 299.852051 -1.48597 344.637268 26.22492 387.44278 C 60.241844 439.991028 126.388168 476.536499 197.991928 481.9263 C 204.940216 482.628845 212.049393 483 219.31575 483 C 341.258392 483 459.927094 381.592651 484.365265 256.5 C 495.761536 198.167572 484.647675 144.990387 457.184875 104.839111 C 483.741425 126.029144 502.082306 149.207672 511.976654 172.183655 C 512.74176 142.71701 494.705383 108.453644 457.684692 79.343597 C 416.277374 46.785736 362.846252 30.009674 316.465363 30.190887 C 313.606476 30.076538 310.730377 30 307.819885 30 Z M 300.94281 65.233368 C 350.338593 65.233368 391.38855 84.945648 418.469452 117.109131 C 438.013794 158.669403 433.832306 222.070923 419.672913 260.38974 C 399.999603 313.632782 369.81665 353.451691 321.711578 382.592712 C 270.391479 413.683441 200.81398 406.227631 179.412796 376.974213 C 177.736526 374.684357 176.358643 372.291809 175.262009 369.81781 C 189.659637 383.442566 222.158676 384.299194 251.176208 371.077576 C 282.466858 356.822998 297.845734 331.611389 285.522247 314.768494 C 273.197571 297.922668 237.842087 295.825684 206.550201 310.081238 C 194.478485 315.581512 184.7892 322.71463 178.133179 330.367157 C 178.755798 329.112244 179.414032 327.857361 180.129974 326.607269 C 157.181458 350.339111 149.966675 377.119568 164.352036 396.78241 C 211.879913 441.283722 284.786652 421.49588 337.194794 390.555359 C 383.799194 363.042297 426.007202 319.282288 446.928223 262.653442 C 448.48172 258.447784 449.895203 254.245056 451.185852 250.048096 C 450.834625 252.197449 450.461304 254.3526 450.038879 256.519409 C 429.399078 362.164368 329.184235 447.804382 226.20018 447.804382 C 182.427658 447.804382 145.205536 432.323975 118.451294 406.416626 C 85.917854 359.51886 98.817268 268.732208 124.192406 226.998474 C 157.036545 172.982147 205.406891 108.029205 289.148682 113.810516 C 345.67572 117.71283 366.636017 170.632233 337.199707 226.432556 C 326.626221 246.477173 311.174927 264.52948 292.732117 279.064148 C 323.322754 262.978027 350.106506 238.238403 365.889374 208.324066 C 395.325653 152.523712 374.412018 98.413391 317.835846 95.135132 C 212.546753 89.034027 148.809845 151.872528 99.805855 219.642487 C 88.495529 238.926422 79.933578 258.743347 74.045097 278.470123 C 74.651741 271.265289 75.656288 263.944183 77.106613 256.519409 C 97.742722 150.87442 197.957535 65.233368 300.94162 65.233368 Z';

const wormholeInnerPath =
	'm 185.97505,441.63265 c -33.79077,-8.3011 -65.17665,-29.43744 -73.93112,-49.78777 -9.25322,-21.50968 -11.70586,-38.04494 -10.49572,-70.76013 1.75896,-47.55179 10.14725,-73.07955 37.78282,-114.98306 16.61208,-25.1887 51.79059,-61.34891 70.20182,-72.16089 42.09293,-24.71904 93.39447,-25.85478 119.75511,-2.65122 22.44745,19.75904 25.90601,55.63103 8.67361,89.96223 -8.34782,16.6309 -23.06348,36.33571 -35.45772,47.47919 -14.63676,13.15971 -8.88632,12.86271 10.17806,-0.52568 78.38243,-55.04583 91.18675,-151.02042 22.80049,-170.900848 -18.56364,-5.396599 -64.73097,-3.085436 -90.19614,4.515268 -38.93881,11.62225 -70.95916,32.06381 -104.36497,66.62583 -25.55133,26.43569 -41.948807,48.50014 -52.889771,71.16848 -10.437743,21.62573 -11.392848,18.12208 -2.38734,-8.75756 11.941674,-35.64352 28.221991,-62.34354 55.144741,-90.43844 80.58664,-84.095152 195.98494,-98.32558 266.94343,-32.91833 10.64131,9.80882 16.62034,24.22178 20.65408,49.78842 14.08686,89.2855 -37.65569,189.72121 -118.32372,229.67411 -24.04789,11.91033 -53.78702,16.90337 -79.00222,13.26405 -17.58109,-2.53748 -39.96863,-12.00443 -46.81155,-19.79505 -4.13868,-4.71185 -4.13707,-4.74168 0.15783,-2.93677 2.38644,1.00289 12.37846,2.24427 22.20448,2.75863 22.21593,1.16294 41.94606,-4.2612 59.10495,-16.24892 29.45219,-20.5762 33.14471,-47.91416 8.01297,-59.32477 -12.36213,-5.61282 -37.8216,-5.59612 -55.91982,0.0367 -49.22461,15.32043 -77.51002,65.89056 -52.4003,93.68401 10.20483,11.29554 30.43008,20.78859 51.00362,23.93938 50.97705,7.80702 120.11367,-20.04727 170.62928,-68.74444 21.16049,-20.39877 34.03655,-36.82299 46.28033,-59.03353 12.97764,-23.54172 9.40453,-10.26326 -4.85272,18.03383 -33.56509,66.61832 -95.43515,115.83499 -163.15777,129.78951 -22.31631,4.59836 -58.97125,4.25083 -79.33674,-0.75221 z';
