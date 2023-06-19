<script lang="ts">
	import { tooltipComponent, tooltipLocation } from '$lib/services/Context';
	import { clamp } from 'lodash-es';

	// close this tooltip when the pointer is let up
	function onPointerUp(e: PointerEvent) {
		window.removeEventListener('pointerup', onPointerUp);
		$tooltipComponent = undefined;
		document.body.className = document.body.className
			.replaceAll('select-none', '')
			.replaceAll('touch-none', '');
	}

	// when the tooltipComponent is set, register a pointerup listener to hide it
	$: {
		if ($tooltipComponent) {
			document.body.className = document.body.className + ' select-none touch-none';
			window.addEventListener('pointerup', onPointerUp);
		}
	}

	$: x = $tooltipLocation.x + componentWidth > window.innerWidth ? 0 : $tooltipLocation.x;
	$: y = Math.max($tooltipLocation.y - componentHeight, 0);

	let component: HTMLElement | undefined;

	// observe tooltip component height changes so we can react
	let componentHeight = 0;
	let componentWidth = 0;
	$: component &&
		new ResizeObserver(() => {
			componentHeight = component?.scrollHeight ?? 0;
			componentWidth = component?.scrollWidth ?? 0;
		}).observe(component);
</script>

<div
	bind:this={component}
	class:block={!!$tooltipComponent}
	class:hidden={!$tooltipComponent}
	class="absolute bg-base-300 rounded-sm p-2 border-2 shadow-md z-50 text-base select-none w-full md:w-auto"
	style={`left: ${x}px; top: ${y}px;`}
>
	{#if $tooltipComponent}
		<svelte:component this={$tooltipComponent.component} {...$tooltipComponent.props} />
	{/if}
</div>
