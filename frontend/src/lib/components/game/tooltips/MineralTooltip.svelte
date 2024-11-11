<script lang="ts" module>
	export type MineralTooltipProps = {
		mineralType: string;
		surfaceAmount: number;
		concentration: number;
		miningRate: number;
		homeworld: boolean;
	};
</script>

<script lang="ts">
	type Props = {
		mineralType: string;
		surfaceAmount: number;
		concentration: number;
		miningRate: number;
		homeworld: boolean;
	};

	let { mineralType, surfaceAmount, concentration, miningRate, homeworld }: Props = $props();
</script>

<div class="flex flex-col sm:w-[26rem] m-auto">
	<div class="text-xl font-semibold text-center">{mineralType}</div>
	<div class="flex flex-row">
		<div class="font-semibold text-right w-40 mr-1">On Surface:</div>
		<div>{surfaceAmount}kT</div>
	</div>
	{#if homeworld}
		<div class="flex flex-row">
			<div class="font-semibold text-right w-40 mr-1">Mineral Concentration:</div>
			<div>{concentration}{concentration < 30 ? ' (30)' : ' (HW)'}</div>
		</div>
	{:else}
		<div class="flex flex-row">
			<div class="font-semibold text-right w-40 mr-1">Mineral Concentration:</div>
			<div>{concentration}</div>
		</div>
	{/if}
	<div class="flex flex-row">
		<div class="font-semibold text-right w-40 mr-1">Mining Rate:</div>
		<div>{miningRate}kT/year</div>
	</div>
</div>
