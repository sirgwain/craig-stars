<script lang="ts">
	import { game, player } from '$lib/services/Context';
	import type { Message } from '$lib/types/Player';
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

<div
	class="card bg-base-200 shadow-xl w-[31rem] max-h-fit min-h-fit rounded-sm border-2 border-base-300"
>
	<div class="card-body p-4 text-base gap-0">
		<div class="flex flex-row items-center">
			<input type="checkbox" class="flex-none checkbox checkbox-xs" />
			<div class="flex-1 text-center text-xl font-semibold text-secondary">
				Year: {$game.year} Message {messageNum + 1} of {$player.messages.length}
			</div>
		</div>
		{#if message}
			<div class="flex flex-row">
				<div class="flex-1 max-h-20 overflow-y-auto">{message.text}</div>
				<div>
					<div class="flex flex-col gap-y-1 ml-1">
						<button
							on:click={previous}
							disabled={messageNum === 0}
							class="btn btn-outline btn-sm normal-case btn-secondary">Prev</button
						>
						<button
							on:click={gotoTarget}
							disabled={!message.targetId}
							class="btn btn-outline btn-sm normal-case btn-secondary">Goto</button
						>
						<button
							on:click={next}
							disabled={messageNum === $player.messages.length - 1}
							class="btn btn-outline btn-sm normal-case btn-secondary">Next</button
						>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>
