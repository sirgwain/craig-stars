<script lang="ts" context="module">
	export type MinesTooltipProps = {
		planetName: string;
		mines: number;
		maxMines: number;
		maxPossibleMines: number;
		canBuildMines: boolean;
	};
</script>

<script lang="ts">
	export let planetName: string;
	export let mines: number;
	export let maxMines: number;
	export let maxPossibleMines: number;
	export let canBuildMines: boolean;
</script>

<div class="flex flex-col sm:w-[26rem] m-auto">
	<div class="text-xl font-semibold text-center">Mines</div>
	<div>
		{#if canBuildMines}
			Your race is incapable of building mines. However, your colonists have an innate mining
			ability equal to one tenth the square root of the starbase population.
		{:else}
			You have <span class="font-semibold">{mines} Mines</span> on <span class="font-semibold">{planetName}</span>.
			You may build up to <span class="font-semibold">{maxPossibleMines} Mines</span>; however, your colonists are
			currently capable of operating only <span class="font-semibold">{maxMines}</span> of them.
		{/if}
	</div>
</div>
