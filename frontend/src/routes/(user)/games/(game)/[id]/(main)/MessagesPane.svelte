<script lang="ts">
	import { goto } from '$app/navigation';
	import { getGameContext } from '$lib/services/Contexts';
	import { commandMapObject, selectMapObject, zoomToMapObject } from '$lib/services/Stores';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType, None, ownedBy } from '$lib/types/MapObject';
	import {
		MessageTargetType,
		MessageType,
		getNextVisibleMessageNum,
		type Message
	} from '$lib/types/Message';
	import {
		ArrowLongLeft,
		ArrowLongRight,
		ArrowTopRightOnSquare,
		MagnifyingGlassMinus,
		MagnifyingGlassPlus
	} from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import hotkeys from 'hotkeys-js';
	import { kebabCase } from 'lodash-es';
	import { onMount } from 'svelte/internal';
	import MessageDetail from '../messages/MessageDetail.svelte';

	const { game, player, universe, settings, messageNum } = getGameContext();

	export let showMessages = false;
	export let messages: Message[];
	let showFilteredMessages = false;

	$: message = messages.length ? messages[$messageNum] : undefined;
	$: nextVisibleMessageNum = getNextVisibleMessageNum(
		$messageNum,
		showFilteredMessages,
		messages,
		$settings
	);
	$: previousVisibleMessageNum = getPreviousVisibleMessageNum(
		$messageNum,
		showFilteredMessages,
		messages
	);
	$: visible = (message && $settings.isMessageVisible(message.type)) ?? false;

	function onFilterMessageType(type: number) {
		if ($settings.isMessageVisible(type)) {
			$settings.filterMessageType(type);
		} else {
			$settings.showMessageType(type);
		}
		$settings = $settings;
		visible = (message && $settings.isMessageVisible(message.type)) ?? false;
	}

	function getPreviousVisibleMessageNum(
		num: number,
		showFilteredMessages: boolean,
		messages: Message[]
	): number {
		for (let i = num - 1; i >= 0; i--) {
			if (showFilteredMessages || $settings.isMessageVisible(messages[i].type)) {
				return i;
			}
		}
		return num;
	}

	const previous = () => {
		$messageNum = getPreviousVisibleMessageNum($messageNum, showFilteredMessages, messages);
	};
	const next = () => {
		$messageNum = getNextVisibleMessageNum($messageNum, showFilteredMessages, messages, $settings);
	};

	function isMessageGotoable(message: Message | undefined): boolean {
		if (!message) {
			return false;
		}

		if (message.targetType !== MessageTargetType.None) {
			return true;
		}

		if (message.type === MessageType.GainTechLevel) {
			return true;
		}

		return false;
	}

	const gotoTarget = () => {
		if (message) {
			const targetType = message.targetType ?? MessageTargetType.None;
			let moType = MapObjectType.None;

			if (message.battleNum) {
				goto(`/games/${$game.id}/battles/${message.battleNum}`);
				return;
			}

			if (message.type === MessageType.GainTechLevel) {
				goto(`/games/${$game.id}/research`);
			}

			if (message.type === MessageType.TechGained && message.spec.techGained) {
				goto(`/games/${$game.id}/techs/${kebabCase(message.spec.techGained)}`);
			}

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
					case MessageTargetType.MineralPacket:
						moType = MapObjectType.MineralPacket;
						break;
					case MessageTargetType.Battle:
						break;
				}

				if (moType != MapObjectType.None) {
					const target = $universe.getMapObject(message);
					if (target) {
						if (target.type == MapObjectType.Fleet) {
							const orbitingPlanetNum = (target as Fleet).orbitingPlanetNum;
							if (orbitingPlanetNum && orbitingPlanetNum != None) {
								const orbiting = $universe.getPlanet(orbitingPlanetNum);
								if (orbiting) {
									selectMapObject(orbiting);
								}
							}
						} else {
							selectMapObject(target);
						}
						if (ownedBy(target, $player.num)) {
							commandMapObject(target);
						}

						// zoom on goto
						zoomToMapObject(target);
					}
				}
			}
		}
	};

	onMount(() => {
		hotkeys('up', () => {
			showMessages = true;
			previous();
		});
		hotkeys('down', () => {
			showMessages = true;
			next();
		});
		hotkeys('enter', () => {
			showMessages = true;
			gotoTarget();
		});

		return () => {
			hotkeys.unbind('up');
			hotkeys.unbind('down');
			hotkeys.unbind('enter');
		};
	});
</script>

<div class:hidden={!showMessages} class:block={showMessages}>
	<div class="card bg-base-200 shadow rounded-sm border-2 border-base-300">
		<div class="card-body p-1 gap-0">
			<div class="flex flex-row items-center mb-1">
				<div class="tooltip tooltip-right" data-tip="Filter these types of messages">
					<input
						type="checkbox"
						class="flex-initial checkbox checkbox-xs"
						checked={visible}
						on:click={() => message && onFilterMessageType(message.type)}
					/>
				</div>

				<div class="flex-1 text-center text-lg font-semibold text-secondary">
					Year: {$game.year} Message {$messageNum + 1} of {messages.length}
				</div>
				<div
					class="tooltip tooltip-left"
					data-tip={showFilteredMessages ? 'Hide filtered messages' : 'Show all messages'}
				>
					<label class="swap">
						<!-- this hidden checkbox controls the state -->
						<input type="checkbox" bind:checked={showFilteredMessages} />

						<!-- filter messages -->
						<Icon src={MagnifyingGlassMinus} size="24" class="swap-off" />

						<!-- show filtered messages -->
						<Icon src={MagnifyingGlassPlus} size="24" class="swap-on" />
					</label>
				</div>
			</div>
			<div class="flex flex-row">
				<div class="mt-1 h-12 grow overflow-y-auto">
					<div class="relative">
						{#if !visible || message == undefined}
							<div class="absolute w-full text-center">
								<span class="text-[1.5rem] text-warning -rotate-12">FILTERED</span>
							</div>
						{/if}
						{#if message}
							<MessageDetail {message} />
						{/if}
					</div>
				</div>
				<div>
					<div class="flex flex-col gap-y-1 ml-1">
						<div class="flex flex-row btn-group">
							<div class="tooltip" data-tip="previous">
								<button
									on:click={previous}
									disabled={$messageNum === previousVisibleMessageNum}
									class="btn btn-outline btn-sm normal-case btn-secondary"
									title="previous"
									><Icon src={ArrowLongLeft} size="16" class="hover:stroke-accent inline" /></button
								>
							</div>
							<div class="tooltip" data-tip="goto">
								<button
									on:click={gotoTarget}
									disabled={!isMessageGotoable(message)}
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
									disabled={$messageNum === nextVisibleMessageNum}
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
		</div>
	</div>
</div>
