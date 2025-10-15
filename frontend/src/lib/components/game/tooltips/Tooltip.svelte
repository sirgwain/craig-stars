<script lang="ts">
	import { tooltipComponent, tooltipLocation } from '$lib/services/Stores';
	import {
		computePosition,
		detectOverflow,
		flip,
		offset,
		shift,
		type Middleware,
		type VirtualElement
	} from '@floating-ui/dom';

	let component: HTMLElement | undefined = $state();

	// Snap to viewport top-left when any side overflows the viewport
	const snapTopLeft: Middleware = {
		name: 'snapTopLeft',
		async fn(state) {
			const overflow = await detectOverflow(state, {
				rootBoundary: 'viewport',
				padding: 0
			});

			if (overflow.top > 0 || overflow.right > 0 || overflow.bottom > 0 || overflow.left > 0) {
				return { x: 0, y: 0 };
			}
			return {};
		}
	};

	// update the position of the tooltip to optimal
	async function updatePosition() {
		if (!component?.classList.contains('block')) return;

		const virtualElement: VirtualElement = {
			getBoundingClientRect: () => new DOMRect($tooltipLocation.x, $tooltipLocation.y, 0, 0)
		};

		const { x, y } = await computePosition(virtualElement, component, {
			placement: 'right',
			strategy: 'fixed',
			middleware: [
				offset(10),
				shift(),
				flip({
					// Ensure we flip to the perpendicular axis if it doesn't fit
					// on narrow viewports.
					crossAxis: 'alignment',
					fallbackAxisSideDirection: 'start'
				}),
				snapTopLeft // fallback to top left if we overflow
			]
		});

		// Floating UI defaults to transforms, but fixed + left/top is fine too.
		component.style.position = 'fixed';
		component.style.left = `${x}px`;
		component.style.top = `${y}px`;
	}

	// Close on pointer up (matches your current UX)
	function onPointerUp() {
		window.removeEventListener('pointerup', onPointerUp);
		$tooltipComponent = undefined;
		document.body.className = document.body.className
			.replaceAll('select-none', '')
			.replaceAll('touch-none', '');
	}

	// When tooltip opens/closes, manage body classes + listeners
	$effect(() => {
		if ($tooltipComponent) {
			document.body.className = document.body.className + ' select-none touch-none';
			window.addEventListener('pointerup', onPointerUp);
			updatePosition();
		} else {
			window.removeEventListener('pointerup', onPointerUp);
		}
	});
</script>

<div
	bind:this={component}
	class:block={!!$tooltipComponent}
	class:hidden={!$tooltipComponent}
	class="fixed bg-base-300 rounded-sm p-2 border-2 shadow-md z-[1000] text-base select-none w-full md:w-max top-0 left-0 max-h-full"
	role="tooltip"
>
	{#if $tooltipComponent}
		{@const SvelteComponent = $tooltipComponent.component}
		<SvelteComponent {...$tooltipComponent.props} />
	{/if}
</div>
