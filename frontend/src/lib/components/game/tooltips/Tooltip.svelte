<script lang="ts">
	import { tooltipComponent, tooltipLocation } from '$lib/services/Context';
	
	// close this tooltip when the pointer is let up
	function onPointerUp(e: PointerEvent) {
		window.removeEventListener('pointerup', onPointerUp);
		$tooltipComponent = undefined;
	}

	// when the tooltipComponent is set, register a pointerup listener to hide it
	$: $tooltipComponent && window.addEventListener('pointerup', onPointerUp);

	$: x = Math.min($tooltipLocation.x, window.innerWidth - 400);
	$: y = Math.max($tooltipLocation.y - componentHeight, 0);

	let component: HTMLElement | undefined;

	// observe tooltip component height changes so we can react
	let componentHeight = 0;
	$: component &&
		new ResizeObserver(() => (componentHeight = component?.clientHeight ?? 0)).observe(component);
</script>

<div
	bind:this={component}
	class:block={!!$tooltipComponent}
	class:hidden={!$tooltipComponent}
	class="absolute bg-base-300 rounded-sm p-2 shadow-md z-50 text-base"
	style={`left: ${x}px; top: ${y}px;`}
>
	{#if $tooltipComponent}
		<svelte:component this={$tooltipComponent.component} {...$tooltipComponent.props} />
	{/if}
</div>
