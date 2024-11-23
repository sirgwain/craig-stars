<script lang="ts">
	import type { Race } from '$lib/types/Race';
	import { loadWasm, type CS } from '$lib/wasm';
	import { User } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { onMount } from 'svelte';

	type Props = {
		race: Race;
		onPointsUpdated?: (points: number) => void;
	};

	let { race, onPointsUpdated }: Props = $props();
	let points = $state(0);

	let cs: CS | undefined = $state();

	onMount(async () => {
		cs = await loadWasm();
	});

	// update points from the server anytime things change
	const computeRacePoints = async (race: Race) => {
		if (cs) {
			points = cs.calculateRacePoints(race) ?? 0;
		}
	};

	$effect(() => {
		race && cs && computeRacePoints(race);
		onPointsUpdated?.(points);
	});
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
