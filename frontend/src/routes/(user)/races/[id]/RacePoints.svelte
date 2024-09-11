<script lang="ts">
	import { assets } from '$app/paths';
	import type { Race } from '$lib/types/Race';
	import { loadWasm, type CS } from '$lib/wasm';
	import { User } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { onMount } from 'svelte';

	export let race: Race;
	export let points: number;

	let cs: CS | undefined;

	onMount(async () => {
		cs = await loadWasm();
	});

	// update points from the server anytime things change
	const computeRacePoints = async (race: Race) => {
		if (cs) {
			points = cs.calculateRacePoints(race) ?? 0;
		}
	};

	$: race && cs && computeRacePoints(race);
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
				<div class="stat-desc pt-1" />
			</div>
		</div>
	</div>
</div>
