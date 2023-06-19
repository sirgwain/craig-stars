<script lang="ts">
	import type { FullGame } from '$lib/services/FullGame';
	import { ownedBy } from '$lib/types/MapObject';
	import { MineFieldType, type MineField } from '$lib/types/MineField';
	import type { Player } from '$lib/types/Player';

	export let game: FullGame;
	export let mineField: MineField;
	export let player: Player;
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
		<div class="text-center">{game.getPlayerName(mineField.playerNum)}</div>
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
		{#if ownedBy(mineField, player.num)}
			<div class="flex flex-row">
				<div class="w-24">Decay Rate:</div>
				<div>
					{mineField.spec.decayRate} / year
				</div>
			</div>
		{/if}
	</div>
</div>
