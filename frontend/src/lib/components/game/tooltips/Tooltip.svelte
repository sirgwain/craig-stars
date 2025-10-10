<script lang="ts">
	import { createPopper, type Instance, type Placement } from '@popperjs/core';
	import { tooltipComponent, tooltipLocation } from '$lib/services/Stores';

	// UI knobs
	// const minWidth = 380;
	// const minHeight = 380;
	const placement: Placement = 'top-start'; // you were placing above the cursor

	function onPointerUp() {
		window.removeEventListener('pointerup', onPointerUp);
		$tooltipComponent = undefined;
		document.body.className = document.body.className
			.replaceAll('select-none', '')
			.replaceAll('touch-none', '');
	}

	let component: HTMLElement | undefined = $state();
	let popper: Instance | null = null;

	// Track dynamic size so we can nudge Popper to recompute
	const resizeObserver = new ResizeObserver(() => {
		popper?.update();
	});

	// --- Virtual reference for pointer-based positioning ---
	let _rect = new DOMRect($tooltipLocation.x, $tooltipLocation.y, 0, 0);
	const virtualRef = {
		getBoundingClientRect: () => _rect,
		// optional but nice: width/height 0 makes it behave like a point
		contextElement: undefined as Element | undefined
	};

	// Whenever the store’s point moves, update the virtual rect and ask Popper to recompute
	$effect(() => {
		_rect = new DOMRect($tooltipLocation.x, $tooltipLocation.y, 0, 0);
		popper?.update();
	});

	// When tooltip is shown/hidden, manage body classes and global listener
	$effect(() => {
		if ($tooltipComponent) {
			document.body.className = document.body.className + ' select-none touch-none';
			window.addEventListener('pointerup', onPointerUp);
		}
	});

	// Start/stop observing size when the element exists
	$effect(() => {
		if (component) {
			resizeObserver.disconnect();
			resizeObserver.observe(component);
		}
	});

	// Create/destroy Popper instance when tooltip appears/disappears or when the element changes
	$effect(() => {
		if (!$tooltipComponent || !component) {
			popper?.destroy();
			popper = null;
			return;
		}

		// Important: use fixed strategy so it’s independent of any scrolling ancestors / portals
		popper = createPopper(virtualRef, component, {
			placement,
			strategy: 'fixed',
			modifiers: [
				{ name: 'offset', options: { offset: [8, 8] } }, // a little gap from the cursor
				{ name: 'preventOverflow', options: { boundary: 'viewport' } },
				{ name: 'flip', options: { fallbackPlacements: ['bottom-start', 'right', 'left'] } },
				// Use top/left instead of transforms if you prefer; transforms are fine for tooltips:
				{ name: 'computeStyles', options: { gpuAcceleration: true } }
			]
		});
	});

	// If you ever want to anchor to a real element instead of the cursor,
	// just pass that element to createPopper(...) instead of virtualRef.
</script>

<div
	bind:this={component}
	class:block={!!$tooltipComponent}
	class:hidden={!$tooltipComponent}
	class="bg-base-300 rounded-sm p-2 border-2 shadow-md z-[1000] text-base select-none w-full md:w-auto"
	role="tooltip"
>
	{#if $tooltipComponent}
		{@const SvelteComponent = $tooltipComponent.component}
		<SvelteComponent {...$tooltipComponent.props} />
	{/if}
</div>
