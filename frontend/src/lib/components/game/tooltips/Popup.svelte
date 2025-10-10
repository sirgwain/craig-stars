<script lang="ts" module>
	import type { Component } from 'svelte';
	import { writable } from 'svelte/store';

	/** Common props all popups accept */
	export type PopupPropsBase = {
		onClose?: () => void;
	};

	/** Generic popup entry, tied to a specific component type */
	export type PopupComponentStore<T extends PopupPropsBase = PopupPropsBase> = {
		component: Component<T>;
		props?: Omit<T, 'onClose'>;
	};

	/** Non-generic erased base — the store can hold *any* popup component conforming to PopupPropsBase */
	export type PopupComponentBase = {
		component: Component<PopupPropsBase>;
		props?: Record<string, unknown>;
	};

	// ⬇ store uses the base type — no any, no unsafe cast
	export const popupComponent = writable<PopupComponentBase | undefined>();
	export const popupLocation = writable<{ x: number; y: number }>({ x: 0, y: 0 });

	/** Show a popup at screen coords (x, y) with a Svelte component + props (minus onClose). */
	export function showPopup<T extends PopupPropsBase>(
		x: number,
		y: number,
		component: Component<T>,
		props?: Omit<T, 'onClose'>
	) {
		popupLocation.set({ x, y });

		// Explicitly widen to base type, safely (no unknown or any)
		popupComponent.set({
			component: component as Component<PopupPropsBase>,
			props: props as Record<string, unknown>
		});
	}
</script>

<script lang="ts">
	import { clickOutside } from '$lib/clickOutside';
	import { createPopper, type Instance, type VirtualElement, type Placement } from '@popperjs/core';

	const minWidth = 250;
	const minHeight = 250;
	const placement: Placement = 'bottom-start';

	function hide() {
		$popupComponent = undefined;
		document.body.className = document.body.className
			.replaceAll('select-none', '')
			.replaceAll('touch-none', '');
	}

	let component: HTMLElement | undefined = $state();
	let popper: Instance | null = null;

	// Track dynamic size; ask Popper to recompute when content changes
	const resizeObserver = new ResizeObserver(() => {
		popper?.update();
	});

	// Body UX class toggling
	$effect(() => {
		if ($popupComponent) {
			document.body.className = document.body.className + ' select-none touch-none';
		}
	});

	// Observe size when element exists
	$effect(() => {
		if (component) {
			resizeObserver.disconnect();
			resizeObserver.observe(component);
		}
	});

	// --- Virtual reference (anchor at screen coords) ---
	let _rect = new DOMRect($popupLocation.x, $popupLocation.y, 0, 0);
	const virtualRef: VirtualElement = {
		getBoundingClientRect: () => _rect,
		contextElement: undefined
	};

	// Update anchor when coords change
	$effect(() => {
		_rect = new DOMRect($popupLocation.x, $popupLocation.y, 0, 0);
		popper?.update();
	});

	// Create / destroy Popper when popup appears or element changes
	$effect(() => {
		if (!$popupComponent || !component) {
			popper?.destroy();
			popper = null;
			return;
		}

		popper = createPopper(virtualRef, component, {
			placement,
			strategy: 'fixed', // stable on scroll; good for viewport-anchored popups
			modifiers: [
				{ name: 'offset', options: { offset: [0, 8] } },
				{ name: 'preventOverflow', options: { boundary: 'viewport', padding: 8 } },
				{
					name: 'flip',
					options: { fallbackPlacements: ['top-start', 'right-start', 'left-start'] }
				},
				{ name: 'computeStyles', options: { gpuAcceleration: true } }
			]
		});
	});
</script>

{#if $popupComponent}
	{@const SvelteComponent = $popupComponent.component}
	<div
		bind:this={component}
		use:clickOutside={hide}
		class="z-50 bg-base-200 rounded-md overflow-y-auto shadow-md border w-auto max-w-[min(95vw,48rem)]"
		style={`min-width:${minWidth}px; min-height:${minHeight}px;`}
		role="dialog"
	>
		<!-- Inject onClose for free; user props never had to include it -->
		<SvelteComponent {...$popupComponent.props} onClose={hide} />
	</div>
{/if}
