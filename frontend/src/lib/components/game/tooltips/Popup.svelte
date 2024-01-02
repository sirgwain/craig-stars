<script lang="ts" context="module">
	export type PopupEvent = {
		close: { event?: Event };
	};
</script>

<script lang="ts">
	import { clickOutside } from '$lib/clickOutside';
	import { popupComponent, popupLocation } from '$lib/services/Stores';

	const minWidth = 250;
	const minHeight = 250;

	function hide(event: MouseEvent) {
		$popupComponent = undefined;
		document.body.className = document.body.className
			.replaceAll('select-none', '')
			.replaceAll('touch-none', '');
	}

	// when the popupComponent is set, register a pointerup listener to hide it
	$: {
		if ($popupComponent) {
			document.body.className = document.body.className + ' select-none touch-none';
		}
	}

	$: x =
		$popupLocation.x + componentWidth > window.innerWidth // we overshoot the window, move the popup left so it fits, or 0 if required
			? Math.max(0, $popupLocation.x - (componentWidth + $popupLocation.x - window.innerWidth) - 20)
			: $popupLocation.x;
	$: y = window.scrollY + Math.min($popupLocation.y, window.innerHeight - componentHeight);

	let component: HTMLElement | undefined;

	// observe popup component height changes so we can react
	let componentHeight = minHeight;
	let componentWidth = minWidth;
	$: component &&
		new ResizeObserver(() => {
			componentHeight = Math.max(component?.scrollHeight ?? 0, minHeight);
			componentWidth = Math.max(component?.scrollWidth ?? 0, minWidth);
		}).observe(component);
</script>

{#if $popupComponent}
	<div
		bind:this={component}
		use:clickOutside={hide}
		class:block={!!$popupComponent}
		class:hidden={!$popupComponent}
		class={`absolute bg-base-200 w-[${minWidth}px] h-[${minHeight}px] rounded-md overflow-y-auto z-50`}
		style={`left: ${x}px; top: ${y}px;`}
	>
		<svelte:component this={$popupComponent.component} {...$popupComponent.props} on:close={hide} />
	</div>
{/if}
