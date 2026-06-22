<script module lang="ts">
	import { popupComponent, popupLocation } from './Popup';
</script>

<script lang="ts">
	import { clickOutside } from '$lib/clickOutside';
	import { computePosition, flip, offset, shift, type VirtualElement } from '@floating-ui/dom';

	let component: HTMLElement | undefined = $state();

	// position the popup
	async function updatePosition() {
		if (!component) return;

		const virtualElement: VirtualElement = {
			getBoundingClientRect: () => new DOMRect($popupLocation.x, $popupLocation.y, 0, 0)
		};

		const { x, y } = await computePosition(virtualElement, component, {
			placement: 'bottom-start',
			strategy: 'fixed',
			middleware: [
				offset(5),
				flip({ fallbackPlacements: ['top-start', 'right-start', 'left-start'] }),
				shift()
			]
		});

		component.style.position = 'fixed';
		component.style.left = `${x}px`;
		component.style.top = `${y}px`;
	}

	// Update anchor when coords change (autoUpdate can't detect virtualRef changes)
	$effect(() => {
		if (!$popupComponent) {
			document.body.className = document.body.className
				.replaceAll('select-none', '')
				.replaceAll('touch-none', '');

			return;
		}

		document.body.className = document.body.className + ' select-none touch-none';
		updatePosition();
	});
</script>

{#if $popupComponent}
	{@const SvelteComponent = $popupComponent.component}
	<div
		use:clickOutside={() => popupComponent.set(undefined)}
		bind:this={component}
		class="fixed z-50 bg-base-200 rounded-md overflow-y-auto shadow-md border w-auto max-w-[min(95vw,48rem)]"
		role="dialog"
	>
		<!-- Inject onClose for free; user props never had to include it -->
		<SvelteComponent {...$popupComponent.props} onClose={() => popupComponent.set(undefined)} />
	</div>
{/if}
