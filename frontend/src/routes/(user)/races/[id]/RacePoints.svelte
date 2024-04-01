<script lang="ts">
	import { Service } from '$lib/services/Service';
	import type { Race } from '$lib/types/Race';
	import { User } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

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
