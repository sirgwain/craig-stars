<script lang="ts">
	import type { GameContext } from '$lib/types/GameContext';
	import { Viewport } from 'pixi-viewport';
	import { Application, BaseTexture, Container, Sprite, Texture } from 'pixi.js-legacy';
	import { getContext, onMount } from 'svelte';

	const { game, player } = getContext<GameContext>('game');

	let app: Application;
	let viewport: Viewport;
	let planetsContainer: Container;

	const unknownPlanetTexture = new Texture(BaseTexture.from('../../scanner/planet-unknown.png'));

	onMount(() => {
		setupApp();
	});

	function setupApp() {
		const width = el.clientWidth;
		const height = el.clientHeight;
		app = new Application({
			width: width,
			height: height,
			backgroundColor: 0x000000, // black hexadecimal
			resolution: window.devicePixelRatio || 1,
			antialias: true,
			autoDensity: true,
			forceCanvas: true
		});

		app.ticker.autoStart = false;
		app.ticker.maxFPS = 24;

		// create viewport
		viewport = new Viewport({
			screenWidth: window.innerWidth,
			screenHeight: window.innerHeight,
			worldWidth: game.area.x,
			worldHeight: game.area.y,
			divWheel: app.renderer.view, // Ensures that when using the scroll wheel it only takes affect when the mouse is over the canvas (prevents scrolling in overlaying divs from scrolling the canvas)
			stopPropagation: true,
			passiveWheel: true,
			interaction: app.renderer.plugins.interaction, // the interaction module is important for wheel() to work properly when renderer.view is placed or scaled
			disableOnContextMenu: true
		});
		// add the viewport to the stage
		app.stage.addChild(viewport);
		// Add a new map to the viewport
		// map = new Map(app, store, this);
		planetsContainer = new Container();

		player.planetIntels.forEach((planet) => {
			const sprite = new Sprite(unknownPlanetTexture);
			sprite.position.x = planet.position.x;
			sprite.position.y = planet.position.y;
			planetsContainer.addChild(sprite);
		});

		viewport.addChild(planetsContainer);
		app.start();
		app.render();

		mountView();
		setupViewport();
	}

	function handleResize() {
		app.renderer.resize(el.clientWidth, el.clientHeight);
	}

	function setupViewport() {
		// activate plugins
		viewport
			.drag()
			.pinch()
			.wheel({
				smooth: 5
			})
			.decelerate({ friction: 0.9 })
			.clamp({
				left: -viewport.worldWidth,
				top: -viewport.worldHeight,
				right: viewport.worldWidth * 2,
				bottom: viewport.worldHeight * 2,
				underflow: 'top-left'
			})
			.clampZoom({ minScale: 0.5, maxScale: 6 });
	}

	let el: HTMLElement;
	const mountView = () => {
		if (!el) return; // not ready
		const { firstChild } = el;
		if (firstChild) {
			if (app.view === firstChild) return; // didn't change
			if (el.childNodes.length > 1) {
				throw Error(`PixiView has ${el.childNodes.length} child nodes! Expected 0 or 1.`);
			}
			el.removeChild(firstChild);
		}
		el.appendChild(app.view);
	};
</script>

<svelte:window on:resize={handleResize} />

<div class="flex-1 h-full" bind:this={el} />
