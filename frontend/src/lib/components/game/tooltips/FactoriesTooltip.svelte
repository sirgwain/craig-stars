<script lang="ts" context="module">
	export type FactoriesTooltipProps = {
		planetName: string;
		factories: number;
		maxFactories: number;
		maxPossibleFactories: number;
		canBuildFactories: boolean;
	};
</script>

<script lang="ts">
	export let planetName: string;
	export let factories: number;
	export let maxFactories: number;
	export let maxPossibleFactories: number;
	export let canBuildFactories: boolean;
</script>

<div class="flex flex-col sm:w-[26rem] m-auto">
	<div class="text-xl font-semibold text-center">Factories</div>
	<div>
		{#if canBuildFactories}
			Your race is incapable of building factories.
		{:else}
			You have <span class="font-semibold">{factories} Factories</span> on
			<span class="font-semibold">{planetName}</span>. You may build up to
			<span class="font-semibold">{maxPossibleFactories} Factories</span>; however, your colonists
			are currently capable of operating only <span class="font-semibold">{maxFactories}</span> of them
		{/if}
	</div>
</div>
