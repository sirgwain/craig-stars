<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { getScannerTarget } from '$lib/types/Battle';
	import type { MapObject } from '$lib/types/MapObject';
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
		Eye,
		MagnifyingGlassMinus,
		MagnifyingGlassPlus
	} from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import hotkeys from 'hotkeys-js';
	import { onMount } from 'svelte';
	import MessageDetail from '../messages/MessageDetail.svelte';

	const {
		game,
		player,
		universe,
		settings,
		messageNum,
		selectMapObject,
		zoomToMapObject,
		gotoTarget
	} = getGameContext();

	export let showMessages = false;
	export let messages: Message[];
	let showFilteredMessages = false;
	let viewBattle = false;

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
			if (
				i >= 0 &&
				messages.length < i &&
				(showFilteredMessages || $settings.isMessageVisible(messages[i].type))
			) {
				return i;
			}
		}
		return num;
	}

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

	const previous = (event: Event) => {
		event.preventDefault();
		$messageNum = getPreviousVisibleMessageNum($messageNum, showFilteredMessages, messages);
		viewBattle = false;
		showMessages = true;
	};

	const next = (event: Event) => {
		event.preventDefault();
		$messageNum = getNextVisibleMessageNum($messageNum, showFilteredMessages, messages, $settings);
		viewBattle = false;
		showMessages = true;
	};

	function goto() {
		if (message) {
			// battles go to the location on first click, and go to vattle view on second clic
			if (message.battleNum) {
				if (viewBattle) {
					viewBattle = false;
					gotoTarget(message, $game.id, $player.num, $universe);
				} else {
					viewBattle = true;
					const battle = $universe.getBattle(message.battleNum);
					if (battle) {
						const target = getScannerTarget(battle, $universe);
						if (target) {
							selectMapObject(target);
							zoomToMapObject(target);
						} else {
							zoomToMapObject({ position: battle.position } as MapObject);
						}
					}
				}
			} else {
				gotoTarget(message, $game.id, $player.num, $universe);
			}
		}
	}

	onMount(() => {
		hotkeys('up', 'root', previous);
		hotkeys('down', 'root', next);
		hotkeys('enter', 'root', goto);

		return () => {
			hotkeys.unbind('up', 'root', previous);
			hotkeys.unbind('down', 'root', next);
			hotkeys.unbind('enter', 'root', goto);
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
									on:click={goto}
									disabled={!isMessageGotoable(message)}
									class="btn btn-outline btn-sm normal-case btn-secondary"
									title="goto"
									>{#if viewBattle}
										<Icon src={Eye} size="16" class="hover:stroke-accent inline" />
									{:else}
										<Icon
											src={ArrowTopRightOnSquare}
											size="16"
											class="hover:stroke-accent inline"
										/>
									{/if}
								</button>
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
