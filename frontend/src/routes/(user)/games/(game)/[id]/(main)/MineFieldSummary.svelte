<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { ownedBy } from '$lib/types/MapObject';
	import { MineFieldType, type MineField } from '$lib/types/MineField';

	const { player, universe } = getGameContext();

	export let mineField: MineField;
</script>

<div class="flex flex-row min-h-[11rem]">
	<div class="flex flex-col">
		<div class="avatar ">
			<div class="border-2 border-neutral mr-2 p-2 bg-black">
				<div
					class:standard-mine-field={mineField.mineFieldType === MineFieldType.Standard}
					class:heavy-mine-field={mineField.mineFieldType === MineFieldType.Heavy}
					class:speed-bump-mine-field={mineField.mineFieldType === MineFieldType.SpeedBump}
					class="mapobject-avatar bg-black"
				/>
			</div>
		</div>
		<div class="text-center">{$universe.getPlayerName(mineField.playerNum)}</div>
	</div>

	<div class="flex flex-col grow">
		<div class="flex flex-row">
			<div class="w-24">Location:</div>
			<div>
				({mineField.position.x}, {mineField.position.y})
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-24">Field Type:</div>
			<div>
				{mineField.mineFieldType}
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-24">Field Radius:</div>
			<div>
				{mineField.spec.radius.toFixed()} l.y. ({mineField.numMines} mines)
			</div>
		</div>
		{#if ownedBy(mineField, $player.num)}
			<div class="flex flex-row">
				<div class="w-24">Decay Rate:</div>
				<div>
					{mineField.spec.decayRate} / year
				</div>
			</div>
		{/if}
	</div>
</div>
