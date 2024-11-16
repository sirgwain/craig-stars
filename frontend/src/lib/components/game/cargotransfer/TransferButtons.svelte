<script lang="ts">
	import { ArrowLongLeft, ArrowLongRight } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createEventDispatcher } from 'svelte';
	import type { HTMLAttributes, HTMLButtonAttributes } from 'svelte/elements';
	type Props = {
		ontransfertosource: () => void;
		ontransfertodest: () => void;
	} & HTMLAttributes<HTMLDivElement>;

	let { ontransfertosource, ontransfertodest, ...props }: Props = $props();
</script>

<div class="flex justify-center {props.class}">
	<button
		onclick={ontransfertosource}
		class="btn btn-outline btn-xs normal-case btn-secondary inline-block p-1"
		><Icon src={ArrowLongLeft} size="16" class="hover:stroke-accent inline" />
	</button>
	<button
		onclick={ontransfertodest}
		class="btn btn-outline btn-xs normal-case btn-secondary inline-block p-1"
		><Icon src={ArrowLongRight} size="16" class="hover:stroke-accent inline" /></button
	>
</div>
