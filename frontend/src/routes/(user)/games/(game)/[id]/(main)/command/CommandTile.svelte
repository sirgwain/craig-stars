<script lang="ts">
	import { getCarouselContext } from '$lib/services/CarouselContext';
	import { ChevronDown, ChevronUp } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import type { Snippet } from 'svelte';
	import { readable } from 'svelte/store';

	type Props = {
		title?: string;
		children?: Snippet;
	};

	let { title = '', children }: Props = $props();

	// if we are in a CommandPaneCarousel, show the disclosure chevrons and hide/show the command pane on click
	let carouselContext = getCarouselContext();
	let showDisclosure = carouselContext != undefined;
	let open = carouselContext ? carouselContext.open : readable<boolean>(true);
</script>

<div class="w-screen md:w-[14rem] card bg-base-200 shadow rounded-sm border-2 border-base-300">
	<div class="card-body p-3 gap-0">
		<button
			class:cursor-default={!showDisclosure}
			class="w-full"
			onclick={carouselContext?.onDisclosureClicked}
		>
			<div class="flex flex-row items-center mb-1">
				<div class="flex-1 text-center text-lg font-semibold text-secondary">
					{title}
				</div>
				{#if showDisclosure}
					{#if $open}
						<Icon src={ChevronUp} size="16" class="hover:stroke-accent" />
					{:else}
						<Icon src={ChevronDown} size="16" class="hover:stroke-accent" />
					{/if}
				{/if}
			</div>
		</button>
		{@render children?.()}
	</div>
</div>
