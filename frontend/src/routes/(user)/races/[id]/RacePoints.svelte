<script lang="ts">
	import type { Race } from '$lib/types/cs';
	import { loadWasm, type CS } from '$lib/wasm';
	import { User } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { onMount } from 'svelte';

	type Props = {
		race: Race;
		onPointsUpdated?: (points: number) => void;
	};

	let { race, onPointsUpdated }: Props = $props();

	let cs: CS | undefined = $state();
	let points = $derived(cs ? (cs.calculateRacePoints(race) ?? 0) : 0);

	onMount(async () => {
		cs = await loadWasm();
	});

	$effect(() => onPointsUpdated?.(points));
</script>

<div class="sticky top-[4rem] z-10">
	<div class="flex justify-end">
		<div class="stats stats-horizontal shadow border border-base-200">
			<div class="stat place-items-center">
				<div class="stat-title">Points</div>
				<div class="stat-figure"><Icon class="w-8 h-8" src={User} /></div>
				<div class="stat-value" class:text-error={points < 0} class:text-success={points >= 0}>
					{#if cs}
						{points}
					{/if}
				</div>
				<div class="stat-desc pt-1"></div>
			</div>
		</div>
	</div>
</div>
