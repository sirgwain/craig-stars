<script lang="ts">
	import type { Race } from '#lib/types/cs-proto.js';
	import type { WasmClient } from '#lib/wasm.js';
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

	// once the page is scrolled the box shrinks into the corner, so it covers less of the form
	const compactAfter = 60;
	let scrollY = $state(0);
	let compact = $derived(scrollY > compactAfter);
</script>

<svelte:window bind:scrollY />

<!-- The box stays under the menu and only the box itself takes clicks, not the row it sits in.
     It shrinks with a transform so nothing on the page moves when it does. -->
<div class="sticky top-16 z-5 pointer-events-none" data-type="race-points">
	<div class="flex justify-end">
		<div
			class="stats stats-horizontal shadow-sm border border-base-200 bg-base-100 pointer-events-auto origin-top-right transition-transform duration-200"
			class:scale-50={compact}
			data-compact={compact}
		>
			<div class="stat place-items-center">
				<div class="stat-title">Points</div>
				<div class="stat-figure"><Icon class="size-8 " src={User} /></div>
				<div class="stat-value" class:text-error={points < 0} class:text-success={points >= 0}>
					{points}
				</div>
				<div class="stat-desc pt-1"></div>
			</div>
		</div>
	</div>
</div>
