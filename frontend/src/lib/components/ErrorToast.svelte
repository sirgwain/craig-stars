<script lang="ts">
	import { goto } from '$app/navigation';
	import { gameKey, getGameContext } from '$lib/services/GameContext';
	import { CSError, errors } from '$lib/services/Errors';
	import { FullGame } from '$lib/services/FullGame';
	import { hasContext } from 'svelte';
	import type { Readable } from 'svelte/store';
	import { fade } from 'svelte/transition';

	let game: Readable<FullGame> | undefined;
	let resetContext: (fg: FullGame) => void;
	$: hasContext(gameKey) && ({ game, resetContext } = getGameContext());

	function onFadeOut(err: CSError) {
		$errors = $errors.filter((e) => e !== err);
	}

	function getFadeOptions(err: CSError) {
		if (err.statusCode === 409) {
			return { delay: 0 };
		} else {
			return { delay: 3000 };
		}
	}
</script>

<div class="toast toast-top toast-center z-50 w-full md:max-w-2xl">
	{#each $errors as err}
		<div>
			<div
				class="alert alert-error"
				in:fade
				out:fade={getFadeOptions(err)}
				on:introend={() => err.statusCode != 409 && onFadeOut(err)}
			>
				<div>
					<span>{err.error}</span>
					{#if err.statusCode == 409}
						<!-- our game is out of date, offer refresh -->
						<button
							type="button"
							class="btn btn-outline"
							on:click|preventDefault={() => {
								$errors = [];
								if ($game && $game.id) {
									// reload the game
									resetContext(new FullGame());
									goto(`/games/${$game.id}`);
								} else {
									goto('/');
								}
							}}>Refresh</button
						>
					{/if}
				</div>
			</div>
		</div>
	{/each}
</div>
