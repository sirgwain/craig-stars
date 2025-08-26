<script lang="ts">
	import { TechLevelSchema, type Player } from '$lib/types/cs-proto';
	import type { TechLike } from '$lib/types/Tech';
	import { emptyTechLevel } from '$lib/types/TechLevel';
	import { equals } from '@bufbuild/protobuf';

	type Props = {
		tech: TechLike;
		player?: Player | undefined;
	};

	let { tech: techLike, player = undefined }: Props = $props();
</script>

<div class="flex flex-col">
	{#if techLike.tech}
		{@const tech = techLike.tech}
		<div>Tech Req:</div>
		{#if !tech.requirements?.techLevel || equals(TechLevelSchema, emptyTechLevel(), tech.requirements.techLevel)}
			<div>-- None --</div>
		{:else}
			{#if tech.requirements?.techLevel?.energy}
				<div
					class="flex flex-row justify-between gap-1"
					class:text-error={player?.techLevels &&
						player.techLevels.energy < tech.requirements.techLevel.energy}
				>
					<div>Energy:</div>
					<div>{tech.requirements.techLevel.energy}</div>
				</div>
			{/if}
			{#if tech.requirements?.techLevel?.weapons}
				<div
					class="flex flex-row justify-between gap-1"
					class:text-error={player?.techLevels &&
						player.techLevels.weapons < tech.requirements.techLevel.weapons}
				>
					<div>Weapons:</div>
					<div>{tech.requirements.techLevel.weapons}</div>
				</div>
			{/if}
			{#if tech.requirements?.techLevel?.propulsion}
				<div
					class="flex flex-row justify-between gap-1"
					class:text-error={player?.techLevels &&
						player.techLevels.propulsion < tech.requirements.techLevel.propulsion}
				>
					<div>Propulsion:</div>
					<div>{tech.requirements.techLevel.propulsion}</div>
				</div>
			{/if}
			{#if tech.requirements?.techLevel?.construction}
				<div
					class="flex flex-row justify-between gap-1"
					class:text-error={player?.techLevels &&
						player.techLevels.construction < tech.requirements.techLevel.construction}
				>
					<div>Construction:</div>
					<div>{tech.requirements.techLevel.construction}</div>
				</div>
			{/if}
			{#if tech.requirements?.techLevel?.electronics}
				<div
					class="flex flex-row justify-between gap-1"
					class:text-error={player?.techLevels &&
						player.techLevels.electronics < tech.requirements.techLevel.electronics}
				>
					<div>Electronics:</div>
					<div>{tech.requirements.techLevel.electronics}</div>
				</div>
			{/if}
			{#if tech.requirements?.techLevel?.biotechnology}
				<div
					class="flex flex-row justify-between gap-1"
					class:text-error={player?.techLevels &&
						player.techLevels.biotechnology < tech.requirements.techLevel.biotechnology}
				>
					<div>Biotechnology:</div>
					<div>{tech.requirements.techLevel.biotechnology}</div>
				</div>
			{/if}
		{/if}
	{/if}
</div>
