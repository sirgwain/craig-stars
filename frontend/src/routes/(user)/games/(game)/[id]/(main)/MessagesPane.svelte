<script lang="ts">
	import { getGameContext } from '#lib/services/GameContext.js';
	import { getScannerTarget } from '#lib/types/Battle.js';
	import type { PlayerMessage } from '#lib/types/cs-proto.js';
	import {
		MapObjectSchema,
		MapObjectType,
		PlayerMessageTargetType,
		PlayerMessageType
	} from '#lib/types/cs-proto.js';
	import { getMapObjectTarget, getNextVisibleMessageNum } from '#lib/types/Message.js';
	import { create } from '@bufbuild/protobuf';
	import {
		ArrowLongLeft,
		ArrowLongRight,
		ArrowTopRightOnSquare,
		ArrowUturnLeft,
		Eye,
		MagnifyingGlassMinus,
		MagnifyingGlassPlus,
		PencilSquare
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
		gotoTarget,
		openMail
	} = getGameContext();

	type Props = {
		showMessages?: boolean;
		messages: PlayerMessage[];
	};

	let { showMessages = $bindable(false), messages }: Props = $props();
	let showFilteredMessages = $state(false);
	let viewBattle = $state(false);

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
		messages: PlayerMessage[]
	): number {
		for (let i = num - 1; i >= 0; i--) {
			if (
				i >= 0 &&
				messages.length > i &&
				(showFilteredMessages || $settings.isMessageVisible(messages[i].type))
			) {
				return i;
			}
		}
		return num;
	}

	function isMessageGotoable(message: PlayerMessage | undefined): boolean {
		if (!message) {
			return false;
		}

		// battle messages are goto-able
		if (message.battleNum) {
			return true;
		}

		if (
			message.type === PlayerMessageType.PLAYER_MAIL ||
			message.type === PlayerMessageType.PLAYER_GAIN_TECH_LEVEL ||
			message.type === PlayerMessageType.BATTLE_REPORTS ||
			(message.type === PlayerMessageType.PLAYER_TECH_GAINED && message.spec?.techGained) ||
			(message.type === PlayerMessageType.MYSTERY_TRADER_MET_WITH_REWARD &&
				message.spec?.mysteryTrader?.tech)
		) {
			return true;
		}

		if (message.target?.targetType === PlayerMessageTargetType.UNSPECIFIED) {
			return false;
		}

		// the target may be gone, i.e. a fleet that was merged into another fleet
		const target = getMapObjectTarget(message);
		if (target.targetType !== MapObjectType.UNSPECIFIED && !$universe.getMapObject(target)) {
			// a fleet destroyed by mines still goes to its salvage or the minefield
			const spec = message.spec;
			if (spec?.mapObjectTarget?.targetType === MapObjectType.MINEFIELD) {
				const hitPosition = spec.minefieldDamage?.position;
				return !!(
					(hitPosition && $universe.getSalvageAtPosition(hitPosition)) ||
					$universe.getMinefield(
						spec.mapObjectTarget.targetPlayerNum,
						spec.mapObjectTarget.targetNum
					)
				);
			}
			return false;
		}

		return true;
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

	const first = (event: Event) => {
		event.preventDefault();
		const num = getNextVisibleMessageNum(-1, showFilteredMessages, messages, $settings);
		if (num >= 0) {
			$messageNum = num;
			viewBattle = false;
			showMessages = true;
		}
	};

	const last = (event: Event) => {
		event.preventDefault();
		const num = getPreviousVisibleMessageNum(messages.length, showFilteredMessages, messages);
		if (num < messages.length) {
			$messageNum = num;
			viewBattle = false;
			showMessages = true;
		}
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
							zoomToMapObject({
								mapObject: create(MapObjectSchema, { position: battle.position })
							});
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
		hotkeys('home', 'root', first);
		hotkeys('end', 'root', last);

		return () => {
			hotkeys.unbind('up', 'root', previous);
			hotkeys.unbind('down', 'root', next);
			hotkeys.unbind('enter', 'root', goto);
			hotkeys.unbind('home', 'root', first);
			hotkeys.unbind('end', 'root', last);
		};
	});
	let message = $derived(messages.length ? messages[$messageNum] : undefined);
	let nextVisibleMessageNum = $derived(
		getNextVisibleMessageNum($messageNum, showFilteredMessages, messages, $settings)
	);
	let previousVisibleMessageNum = $derived(
		getPreviousVisibleMessageNum($messageNum, showFilteredMessages, messages)
	);
	let visible = $derived((message && $settings.isMessageVisible(message.type)) ?? false);
	// mail from another player is replied to instead of gone to
	let isMail = $derived(message?.type === PlayerMessageType.PLAYER_MAIL);
</script>

<div class:hidden={!showMessages} class:block={showMessages}>
	<div class="card bg-base-200 shadow-sm rounded-xs border-2 border-base-300">
		<div class="card-body p-1 gap-0">
			<div class="flex flex-row items-center mb-1">
				<div class="tooltip tooltip-right" data-tip="Filter these types of messages">
					<input
						type="checkbox"
						class="flex-initial checkbox checkbox-xs"
						checked={visible}
						onclick={() => message && onFilterMessageType(message.type)}
					/>
				</div>

				<div class="flex-1 text-center text-lg font-semibold text-secondary">
					Year: {$game.year} Message {$messageNum + 1} of {messages.length}
				</div>
				<div class="tooltip tooltip-left" data-tip="Write a message to other players">
					<button
						type="button"
						class="mr-1 align-middle"
						onclick={() => openMail()}
						title="write message"
						data-type="write-mail"
						><Icon src={PencilSquare} size="24" class="hover:stroke-accent" /></button
					>
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
						{#if !visible || message == null}
							<div class="absolute w-full text-center">
								<span class="text-[1.5rem] text-warning -rotate-12">FILTERED</span>
							</div>
						{/if}
						{#if message}
							<span class="select-text">
								<MessageDetail {message} />
							</span>
						{/if}
					</div>
				</div>
				<div>
					<div class="flex flex-col gap-y-1 ml-1">
						<div class="flex flex-row">
							<div class="tooltip" data-tip="previous">
								<button
									onclick={previous}
									disabled={$messageNum === previousVisibleMessageNum}
									class="btn btn-outline btn-sm normal-case btn-secondary"
									title="previous"
									><Icon src={ArrowLongLeft} size="16" class="hover:stroke-accent inline" /></button
								>
							</div>
							<div class="tooltip" data-tip={isMail ? 'reply' : 'goto'}>
								<button
									onclick={goto}
									disabled={!isMessageGotoable(message)}
									class="btn btn-outline btn-sm normal-case btn-secondary"
									title={isMail ? 'reply' : 'goto'}
									>{#if isMail}
										<Icon src={ArrowUturnLeft} size="16" class="hover:stroke-accent inline" />
									{:else if viewBattle}
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
									onclick={next}
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
