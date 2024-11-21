<script lang="ts" module>
	export type PopupEvent = {
		close?: { event?: Event };
	};
</script>

<script lang="ts">
	import { clickOutside } from '$lib/clickOutside';
	import { popupComponent, popupLocation } from '$lib/services/Stores';

	const minWidth = 250;
	const minHeight = 250;

	function hide(_event: MouseEvent) {
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
		<SvelteComponent {...$popupComponent.props} on:close={hide} />
	</div>
{/if}
