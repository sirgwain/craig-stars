<script lang="ts">
	import { tooltipComponent, tooltipLocation } from '$lib/services/Stores';

	const minWidth = 380;
	const minHeight = 380;

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

	$: x =
		$tooltipLocation.x + componentWidth > window.innerWidth // we overshoot the window, move the tooltip left so it fits, or 0 if required
			? Math.max(
					0,
					$tooltipLocation.x - (componentWidth + $tooltipLocation.x - window.innerWidth) - 20
			  )
			: $tooltipLocation.x;
	$: y = window.scrollY + Math.max($tooltipLocation.y - componentHeight, 0);

	let component: HTMLElement | undefined;

	// observe tooltip component height changes so we can react
	let componentHeight = minHeight;
	let componentWidth = minWidth;
	$: component &&
		new ResizeObserver(() => {
			componentHeight = Math.max(component?.scrollHeight ?? 0, minHeight);
			componentWidth = Math.max(component?.scrollWidth ?? 0, minWidth);
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
