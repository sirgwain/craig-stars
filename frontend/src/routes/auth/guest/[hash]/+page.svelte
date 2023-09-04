<script lang="ts">
	import { page } from '$app/stores';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import { GameService } from '$lib/services/GameService';
	import { GameState } from '$lib/types/Game';
	import { onMount } from 'svelte';

	let hash = $page.params.hash;
	let loginError = '';

	onMount(async () => {
		const data = JSON.stringify({ hash });

		const response = await fetch(`/api/auth/guest/login?session=1`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				accept: 'application/json'
			},
			body: data
		});

		if (response.ok) {
			const resolvedResponse = (await response?.json()) as { attrs?: { game_id?: string } };
			const gameId = resolvedResponse?.attrs?.game_id;
			if (gameId) {
				const game = await GameService.loadGame(gameId);
				if (game.state == GameState.Setup) {
					document.location = `/join-game/${gameId}`;
				} else {
					document.location = `/games/${gameId}`;
				}
			} else {
				// no game_id, send them to / I guess...
				document.location = '/';
			}
		} else {
			const resolvedResponse = await response?.json();
			loginError = resolvedResponse.error;
			console.error(loginError);
		}
	});
</script>

{#if loginError == ''}
	<div class="flex items-center justify-center min-h-[100dvh] card">
		<div class="px-8 py-6 mt-4 bg-base-200 shadow-xl flex flex-col">
			<LoadingSpinner />
			<div class="mt-2">Signing in as guest</div>
		</div>
	</div>
{:else}
	<div class="flex items-center justify-center min-h-[100dvh] card">
		<div class="px-8 py-6 mt-4 bg-base-200 shadow-xl flex flex-col">
			<div class="mt-2 text-error">Guest Login Failed</div>
		</div>
	</div>
{/if}
