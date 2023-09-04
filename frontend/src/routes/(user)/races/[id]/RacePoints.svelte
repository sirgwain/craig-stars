<script lang="ts">
	import HabChance from '$lib/components/game/race/HabChance.svelte';
	import { HabType } from '$lib/types/Hab';
	import type { Race } from '$lib/types/Race';
	import HabBar from './HabBar.svelte';
	import SpinnerNumberText from '../../../../lib/components/SpinnerNumberText.svelte';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { User } from '@steeze-ui/heroicons';
	import { Service } from '$lib/services/Service';

	export let race: Race;
	export let points: number;

	// update points from the server anytime things change
	const computeRacePoints = async () => {
		const body = JSON.stringify(race);
		const response = await fetch(`/api/races/points`, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			},
			body
		});

		if (!response.ok) {
			await Service.throwError(response);
		}

		const result = (await response.json()) as { points: number };
		points = result.points;
	};

	$: {
		if (race) {
			computeRacePoints();
		}
	}
</script>

<div class="sticky top-[4rem] z-10">
	<div class="flex justify-end">
		<div class="stats stats-horizontal shadow border border-base-200">
			<div class="stat place-items-center">
				<div class="stat-title">Points</div>
				<div class="stat-figure"><Icon class="w-8 h-8" src={User} /></div>
				<div class="stat-value" class:text-error={points < 0} class:text-success={points >= 0}>
					{points}
				</div>
				<div class="stat-desc pt-1" />
			</div>
		</div>
	</div>
</div>
