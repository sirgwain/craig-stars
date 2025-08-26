<script lang="ts">
	import type { Race } from '$lib/types/cs-proto';
	import type { WasmClient } from '$lib/wasm';
	import { User } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	type Props = {
		wasmClient: WasmClient;
		race: Race;
		onPointsUpdated?: (points: number) => void;
	};

	let { wasmClient, race, onPointsUpdated }: Props = $props();

	let points = $state(0);

	$effect(() => {
		wasmClient.calculateRacePoints({ race }).then((res) => (points = res.points));
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
					{points}
				</div>
				<div class="stat-desc pt-1"></div>
			</div>
		</div>
	</div>
</div>
