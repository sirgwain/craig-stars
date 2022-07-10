<script lang="ts">
	import type { Player } from '$lib/types/Player';

	import type { Tech,TechRequirements } from '$lib/types/Tech';

	export let tech: Tech;
	export let player: Player | undefined = undefined;

	const noRequirements: TechRequirements = {};
</script>

<div class="flex flex-col">
	<div>Tech Req:</div>
	{#if JSON.stringify(tech.requirements) == JSON.stringify(noRequirements)}
		<div>-- None --</div>
	{:else}
		{#if tech.requirements.energy}
			<div
				class="flex flex-row justify-between gap-1"
				class:text-error={player?.techLevels &&
					(player.techLevels.energy ?? 0) < (tech.requirements.energy ?? 0)}
			>
				<div>Energy:</div>
				<div>{tech.requirements.energy}</div>
			</div>
		{/if}
		{#if tech.requirements.weapons}
			<div
				class="flex flex-row justify-between gap-1"
				class:text-error={player?.techLevels &&
					(player.techLevels.weapons ?? 0) < (tech.requirements.weapons ?? 0)}
			>
				<div>Weapons:</div>
				<div>{tech.requirements.weapons}</div>
			</div>
		{/if}
		{#if tech.requirements.propulsion}
			<div
				class="flex flex-row justify-between gap-1"
				class:text-error={player?.techLevels &&
					(player.techLevels.propulsion ?? 0) < (tech.requirements.propulsion ?? 0)}
			>
				<div>Propulsion:</div>
				<div>{tech.requirements.propulsion}</div>
			</div>
		{/if}
		{#if tech.requirements.construction}
			<div
				class="flex flex-row justify-between gap-1"
				class:text-error={player?.techLevels &&
					(player.techLevels.construction ?? 0) < (tech.requirements.construction ?? 0)}
			>
				<div>Construction:</div>
				<div>{tech.requirements.construction}</div>
			</div>
		{/if}
		{#if tech.requirements.electronics}
			<div
				class="flex flex-row justify-between gap-1"
				class:text-error={player?.techLevels &&
					(player.techLevels.electronics ?? 0) < (tech.requirements.electronics ?? 0)}
			>
				<div>Electronics:</div>
				<div>{tech.requirements.electronics}</div>
			</div>
		{/if}
		{#if tech.requirements.biotechnology}
			<div
				class="flex flex-row justify-between gap-1"
				class:text-error={player?.techLevels &&
					(player.techLevels.biotechnology ?? 0) < (tech.requirements.biotechnology ?? 0)}
			>
				<div>Biotechnology:</div>
				<div>{tech.requirements.biotechnology}</div>
			</div>
		{/if}
	{/if}
</div>
