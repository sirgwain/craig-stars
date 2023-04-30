<script lang="ts">
	import {
		commandMapObject,
		getMapObject,
		selectMapObject,
		zoomToMapObject
	} from '$lib/services/Context';
	import type { Fleet } from '$lib/types/Fleet';
	import type { Game } from '$lib/types/Game';
	import { MapObjectType, None, ownedBy, type MapObject } from '$lib/types/MapObject';
	import { MessageTargetType, type Message, type Player } from '$lib/types/Player';
	import { ArrowTopRightOnSquare, ArrowLongLeft, ArrowLongRight } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	export let game: Game;
	export let player: Player;

	let messageNum = 0;
	let message: Message | undefined;

	$: player.messages?.length && (message = player.messages[messageNum]);

	const previous = () => {
		messageNum--;
	};
	const next = () => {
		messageNum++;
	};
	const gotoTarget = () => {
		if (message) {
			const targetType = message.targetType ?? MessageTargetType.None;
			let moType = MapObjectType.None;

			if (message.targetNum) {
				switch (targetType) {
					case MessageTargetType.Planet:
						moType = MapObjectType.Planet;
						break;
					case MessageTargetType.Fleet:
						moType = MapObjectType.Fleet;
						break;
					case MessageTargetType.Wormhole:
						moType = MapObjectType.Wormhole;
						break;
					case MessageTargetType.MineField:
						moType = MapObjectType.MineField;
						break;
					case MessageTargetType.MysteryTrader:
						moType = MapObjectType.MysteryTrader;
						break;
					case MessageTargetType.Battle:
						// TODO: handle battles
						break;
				}

				if (moType != MapObjectType.None) {
					const target = getMapObject(moType, message.targetNum, message.targetPlayerNum);
					if (target) {
						if (target.type == MapObjectType.Fleet) {
							const orbitingPlanetNum = (target as Fleet).orbitingPlanetNum;
							if (orbitingPlanetNum && orbitingPlanetNum != None) {
								const orbiting = player.planetIntels[orbitingPlanetNum - 1];
								selectMapObject(orbiting);
							}
						} else {
							selectMapObject(target);
						}
						if (ownedBy(target, player.num)) {
							commandMapObject(target);
						}

						// zoom on goto
						zoomToMapObject(target);
					}
				}
			}
		}
	};
</script>

<div class="card bg-base-200 shadow-xl rounded-sm border-2 border-base-300">
	<div class="card-body p-4 gap-0">
		<div class="flex flex-row items-center">
			<input type="checkbox" class="flex-initial checkbox checkbox-xs" />
			<div class="flex-1 text-center text-lg font-semibold text-secondary">
				Year: {game.year} Message {messageNum + 1} of {player?.messages?.length}
			</div>
		</div>
		{#if message}
			<div class="flex flex-row">
				<div class="flex-1 h-12 overflow-y-auto">{message.text}</div>
				<div>
					<div class="flex flex-col gap-y-1 ml-1">
						<div class="flex flex-row btn-group">
							<div class="tooltip" data-tip="previous">
								<button
									on:click={previous}
									disabled={messageNum === 0}
									class="btn btn-outline btn-sm normal-case btn-secondary"
									title="previous"
									><Icon src={ArrowLongLeft} size="16" class="hover:stroke-accent inline" /></button
								>
							</div>
							<div class="tooltip" data-tip="goto">
								<button
									on:click={gotoTarget}
									disabled={!message.targetNum}
									class="btn btn-outline btn-sm normal-case btn-secondary"
									title="goto"
									><Icon
										src={ArrowTopRightOnSquare}
										size="16"
										class="hover:stroke-accent inline"
									/></button
								>
							</div>
							<div class="tooltip" data-tip="next">
								<button
									on:click={next}
									disabled={player.messages && messageNum === player.messages.length - 1}
									class="btn btn-outline btn-sm normal-case btn-secondary"
									title="next"
									><Icon
										src={ArrowLongRight}
										size="16"
										class="hover:stroke-accent inline"
									/></button
								>
							</div>
						</div>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>