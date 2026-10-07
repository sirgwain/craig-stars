<!--
  @component
  Draws the whole scanner on a single canvas. Scene data is rebuilt when the universe or
  settings change, and the canvas is redrawn at most once per animation frame.
 -->
<script lang="ts">
	import { getGameContext } from '#lib/services/GameContext.js';
	import type { MapObjectLike } from '#lib/types/MapObject.js';
	import type { ZoomTransform } from 'd3-zoom';
	import { onDestroy } from 'svelte';
	import { readScannerColors, type ScannerColors } from './canvas/colors';
	import { drawScanner, type ScannerFrame } from './canvas/draw';
	import {
		buildFleets,
		buildMinefields,
		buildPacketDests,
		buildPlanets,
		buildRouteDests,
		buildScanners,
		buildWaypointPaths,
		buildWormholeLinks
	} from './canvas/scene';
	import { ScannerView } from './canvas/view';

	const {
		player,
		universe,
		settings,
		commandedFleet,
		commandedPlanet,
		commandedMapObject,
		selectedMapObject,
		highlightedMapObject,
		selectedWaypoint
	} = getGameContext();

	type Props = {
		transform: ZoomTransform | undefined;
		// css pixels per light year, before zooming
		pixelsPerLightYear: { x: number; y: number };
		// icons shrink when zoomed out below this
		minObjectZoom: number;
		// show a big X on this map object
		locate: MapObjectLike | undefined;
	};

	let { transform, pixelsPerLightYear, minObjectZoom, locate }: Props = $props();

	let container: HTMLDivElement | undefined = $state();
	let canvas: HTMLCanvasElement | undefined = $state();
	let colors: ScannerColors | undefined = $state();
	let width = $state(0);
	let height = $state(0);
	let dpr = $state(typeof window !== 'undefined' ? window.devicePixelRatio || 1 : 1);

	// scene data, rebuilt only when what it depends on changes
	let scanners = $derived(buildScanners($universe, $player, $settings));
	let minefields = $derived(buildMinefields($universe, $player, $settings));
	let packetDests = $derived(buildPacketDests($universe, $commandedPlanet));
	let routeDests = $derived(buildRouteDests($universe, $commandedPlanet));
	let waypointPaths = $derived(buildWaypointPaths($universe, $player, $settings, $commandedFleet));
	let planets = $derived(
		colors
			? buildPlanets($universe, $player, $settings, $commandedMapObject, $commandedPlanet, colors)
			: []
	);
	let planetsByNum = $derived(new Map(planets.map((p) => [p.planet.mapObject?.num ?? 0, p])));
	let fleets = $derived(
		colors ? buildFleets($universe, $player, $settings, $commandedFleet, colors) : []
	);
	let wormholeLinks = $derived(buildWormholeLinks($universe));

	let view = $derived(
		new ScannerView(
			pixelsPerLightYear.x,
			pixelsPerLightYear.y,
			transform?.k ?? 1,
			transform?.x ?? 0,
			transform?.y ?? 0,
			width,
			height,
			dpr,
			Math.min(1, (transform?.k ?? 1) / minObjectZoom)
		)
	);

	// fade the locator X out after it's hidden
	const locatorFadeMs = 400;
	let locator: ScannerFrame['locator'] = undefined;
	let locatorFadeStart = 0;

	let frame: ScannerFrame | undefined = $derived(
		colors
			? {
					colors,
					universe: $universe,
					player: $player,
					settings: $settings,
					scanners,
					minefields,
					packetDests,
					routeDests,
					waypointPaths,
					planets,
					planetsByNum,
					fleets,
					wormholeLinks,
					selectedMapObject: $selectedMapObject,
					highlightedMapObject: $highlightedMapObject,
					commandedMapObject: $commandedMapObject,
					selectedWaypoint: $selectedWaypoint,
					locator: undefined
				}
			: undefined
	);

	let animationFrame = 0;
	function requestDraw() {
		if (!animationFrame) {
			animationFrame = requestAnimationFrame(draw);
		}
	}

	function draw(time: number) {
		animationFrame = 0;
		const ctx = canvas?.getContext('2d');
		if (!canvas || !ctx || !width || !height || !frame) {
			return;
		}

		// resize the canvas only when we draw, so the browser never stretches a stale image to
		// a new size
		const deviceWidth = Math.round(width * dpr);
		const deviceHeight = Math.round(height * dpr);
		if (canvas.width !== deviceWidth || canvas.height !== deviceHeight) {
			canvas.width = deviceWidth;
			canvas.height = deviceHeight;
			canvas.style.width = `${width}px`;
			canvas.style.height = `${height}px`;
		}

		if (locate) {
			locator = { mapObject: locate, alpha: 1 };
			locatorFadeStart = 0;
		} else if (locator) {
			locatorFadeStart ||= time;
			locator.alpha = 1 - (time - locatorFadeStart) / locatorFadeMs;
			if (locator.alpha <= 0) {
				locator = undefined;
			} else {
				requestDraw();
			}
		}

		drawScanner(ctx, view, { ...frame, locator });

		// expose the view so e2e tests can click on map objects
		canvas.dataset.view = `${view.sx},${view.sy},${view.k},${view.tx},${view.ty}`;
	}

	// CSS colors are cached between frames and refreshed when the app theme changes.
	$effect(() => {
		if (!canvas) {
			return;
		}
		const element = canvas;
		const updateColors = () => (colors = readScannerColors(element));
		updateColors();
	});

	// redraw whenever anything we draw changes
	$effect(() => {
		void [view, frame, locate];
		requestDraw();
	});

	// Canvas pixels can be discarded while a phone sleeps, even when scene data hasn't
	// changed. Redraw after restoration and when returning to the page.
	$effect(() => {
		if (!canvas) {
			return;
		}
		const element = canvas;
		const visible = () => {
			if (document.visibilityState === 'visible') {
				requestDraw();
			}
		};
		element.addEventListener('contextrestored', requestDraw);
		document.addEventListener('visibilitychange', visible);
		window.addEventListener('pageshow', requestDraw);
		return () => {
			element.removeEventListener('contextrestored', requestDraw);
			document.removeEventListener('visibilitychange', visible);
			window.removeEventListener('pageshow', requestDraw);
		};
	});

	// When the scanner is resized (like when the summary below it changes size on a phone),
	// redraw right away. Resize observers run before paint, so the resized scanner shows up in the
	// same frame instead of flashing a stretched image for a frame.
	$effect(() => {
		if (!container) {
			return;
		}
		const observer = new ResizeObserver(([entry]) => {
			width = entry.contentRect.width;
			height = entry.contentRect.height;
			if (animationFrame) {
				cancelAnimationFrame(animationFrame);
			}
			draw(performance.now());
		});
		observer.observe(container);
		return () => observer.disconnect();
	});

	// the device pixel ratio changes when moving between displays or zooming the browser
	$effect(() => {
		const media = matchMedia(`(resolution: ${dpr}dppx)`);
		const update = () => (dpr = window.devicePixelRatio || 1);
		media.addEventListener('change', update);
		return () => media.removeEventListener('change', update);
	});

	onDestroy(() => {
		if (animationFrame) {
			cancelAnimationFrame(animationFrame);
		}
	});
</script>

<div class="absolute inset-0 overflow-hidden" bind:this={container}>
	<canvas bind:this={canvas} data-type="scanner-canvas" class="block"></canvas>
</div>
