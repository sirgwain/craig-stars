<script lang="ts" module>
	import type { Component } from 'svelte';
	import { writable } from 'svelte/store';

	export type PopupProps = {
		onClose?: () => void;
	};

	type PopupComponentStore<T extends PopupProps = PopupProps> = {
		component: Component<T>;
		props?: T;
	};

	export const popupComponent = writable<PopupComponentStore | undefined>();
	export const popupLocation = writable<Vector>({ x: 0, y: 0 });

	export const showPopup = <T extends PopupProps = PopupProps>(
		x: number,
		y: number,
		component: Component<T>,
		props?: T
	) => {
		popupLocation.set({ x, y });
		popupComponent.set({
			// TODO: can't figure out a way around type assertion, but it at least
			// seems to work to force popup components to define an onClose function
			component: component as unknown as Component<PopupProps>,
			props
		});
	};
</script>

<script lang="ts">
	import { clickOutside } from '$lib/clickOutside';
	import type { Vector } from '$lib/types/cs';

	const minWidth = 250;
	const minHeight = 250;

	function hide() {
		$popupComponent = undefined;
		document.body.className = document.body.className
			.replaceAll('select-none', '')
			.replaceAll('touch-none', '');
	}

	let component: HTMLElement | undefined = $state();
	const resizeObserver = new ResizeObserver(() => {
		componentHeight = Math.max(component?.scrollHeight ?? 0, minHeight);
		componentWidth = Math.max(component?.scrollWidth ?? 0, minWidth);
	});

	// observe popup component height changes so we can react
	let componentHeight = $state(minHeight);
	let componentWidth = $state(minWidth);
	// when the popupComponent is set, register a pointerup listener to hide it
	$effect(() => {
		if ($popupComponent) {
			document.body.className = document.body.className + ' select-none touch-none';
		}
	});
	// TODO - JD - 2024-11-20 - This should be run once, I am afraid this could be multiple times
	// in an effect, update held-over
	// CORRECTION - moved resizeObserver, above, and did a disconnect and observe.
	// Also done in Tooltip
	$effect(() => {
		if (component) {
			resizeObserver.disconnect();
			resizeObserver.observe(component);
		}
	});
	let x = $derived(
		$popupLocation.x + componentWidth > window.innerWidth // we overshoot the window, move the popup left so it fits, or 0 if required
			? Math.max(0, $popupLocation.x - (componentWidth + $popupLocation.x - window.innerWidth) - 20)
			: $popupLocation.x
	);
	let y = $derived(
		window.scrollY + Math.min($popupLocation.y, window.innerHeight - componentHeight)
	);
</script>

{#if $popupComponent}
	{@const SvelteComponent = $popupComponent.component}
	<div
		bind:this={component}
		use:clickOutside={hide}
		class:block={!!$popupComponent}
		class:hidden={!$popupComponent}
		class={`absolute bg-base-200 w-[${minWidth}px] h-[${minHeight}px] rounded-md overflow-y-auto z-50`}
		style={`left: ${x}px; top: ${y}px;`}
	>
		<SvelteComponent {...$popupComponent.props} onClose={hide} />
	</div>
{/if}
