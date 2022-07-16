<script lang="ts">
	import { game, player } from '$lib/services/Context';
	import type { Message } from '$lib/types/Player';
	import { ExternalLink, ArrowNarrowLeft, ArrowNarrowRight } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	let messageNum = 0;
	let message: Message | undefined;

	$: $player && (message = $player.messages[messageNum]);

	const previous = () => {
		messageNum--;
	};
	const next = () => {
		messageNum++;
	};
	const gotoTarget = () => {};
</script>

{#if $game}
	<div class="card bg-base-200 shadow-xl rounded-sm border-2 border-base-300">
		<div class="card-body p-4 gap-0">
			<div class="flex flex-row items-center">
				<input type="checkbox" class="flex-initial checkbox checkbox-xs" />
				<div class="flex-1 text-center text-lg font-semibold text-secondary">
					Year: {$game.year} Message {messageNum + 1} of {$player.messages.length}
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
										><Icon
											src={ArrowNarrowLeft}
											size="16"
											class="hover:stroke-accent inline"
										/></button
									>
								</div>
								<div class="tooltip" data-tip="goto">
									<button
										on:click={gotoTarget}
										disabled={!message.targetId}
										class="btn btn-outline btn-sm normal-case btn-secondary"
										title="goto"
										><Icon
											src={ExternalLink}
											size="16"
											class="hover:stroke-accent inline"
										/></button
									>
								</div>
								<div class="tooltip" data-tip="next">
									<button
										on:click={next}
										disabled={messageNum === $player.messages.length - 1}
										class="btn btn-outline btn-sm normal-case btn-secondary"
										title="next"
										><Icon
											src={ArrowNarrowRight}
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
{/if}
