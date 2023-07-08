<!-- svelte-ignore a11y-no-noninteractive-tabindex -->
<script lang="ts">
	import { page } from '$app/stores';
	import DarkModeToggler from '$lib/components/DarkModeToggler.svelte';
	import UserAvatar from '$lib/components/UserAvatar.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import { me } from '$lib/services/Stores';
	import { GameState, type Game } from '$lib/types/Game';
	import { ArrowUpTray, Bars3 } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createEventDispatcher } from 'svelte';

	const { game, player, universe } = getGameContext();

	const dispatch = createEventDispatcher();
	const updateTitle = (game: Game) => (document.title = `${$game.name} - ${$game.year}`);

	// every turn/game update update the game
	$: updateTitle($game);
</script>

<!-- svelte-ignore a11y-label-has-associated-control -->

<div class="navbar bg-base-100 flex flex-row w-full">
	<div class="flex-1">
		<a class="btn btn-ghost text-xl text-primary" href={`/`}>cs</a>
		<div class="md:block">
			<a class="btn btn-ghost text-lg text-accent" href={`/games/${$game.id}`}
				>{$game.name} - {$game.year}</a
			>
		</div>
	</div>
	<div class="flex-initial">
		{#if $page.url.pathname === `/games/${$game.id}` && !$player.submittedTurn && $game.state === GameState.WaitingForPlayers}
			<button
				type="button"
				on:click={() => dispatch('submit-turn')}
				class="btn btn-primary"
				title="submit turn"
				><span class="hidden md:inline-block mr-1">Submit Turn</span><Icon
					src={ArrowUpTray}
					size="16"
				/></button
			>
		{/if}

		{#if !$player.submittedTurn}
			<div class="hidden md:inline-block">
				<div class="dropdown dropdown-end">
					<label for="reports" tabindex="0" class="btn btn-ghost w-40">Commands</label>
					<ul
						id="commands"
						tabindex="0"
						class=" menu menu-compact dropdown-content mt-3 p-2 shadow bg-base-300"
					>
						<li><a href={`/games/${$game.id}/research`}>Research</a></li>
						<li><a href={`/games/${$game.id}/designs`}>Ship Designs</a></li>
						<li><a href={`/games/${$game.id}/battle-plans`}>Battle Plans</a></li>
						<li><a href={`/games/${$game.id}/production-plans`}>Production Plans</a></li>
						<li><a href={`/games/${$game.id}/transport-plans`}>Transport Plans</a></li>
					</ul>
				</div>

				<div class="dropdown dropdown-end">
					<label for="reports" tabindex="0" class="btn btn-ghost">Reports</label>
					<ul
						id="reports"
						tabindex="0"
						class=" menu menu-compact dropdown-content mt-3 p-2 shadow bg-base-300"
					>
						<li><a href={`/games/${$game.id}/players`}>Players</a></li>
						<li><a href={`/games/${$game.id}/planets`}>Planets</a></li>
						<li><a href={`/games/${$game.id}/fleets`}>Fleets</a></li>
						<li><a href={`/games/${$game.id}/messages`}>Messages</a></li>
						<li><a href={`/games/${$game.id}/battles`}>Battles</a></li>
					</ul>
				</div>
			</div>
		{/if}

		<div class="dropdown dropdown-end">
			<label class="avatar w-8 h-8 place-content-center">
				<UserAvatar user={$me} />
			</label>
		</div>

		<div class="dropdown dropdown-end">
			<label for="menu" tabindex="0" class="btn btn-ghost">
				<div id="menu">
					<Icon src={Bars3} size="24" />
				</div>
			</label>
			<ul tabindex="0" class="menu menu-compact dropdown-content mt-3 p-2 shadow bg-base-300 w-52">
				<li>
					<DarkModeToggler />
				</li>
				{#if !$player.submittedTurn}
					<li class="md:hidden menu-title">
						<span>Commands</span>
					</li>
					<li class="md:hidden"><a href={`/games/${$game.id}/research`}>Research</a></li>
					<li class="md:hidden"><a href={`/games/${$game.id}/designs`}>Ship Designs</a></li>
					<li class="md:hidden"><a href={`/games/${$game.id}/battle-plans`}>Battle Plans</a></li>
					<li class="md:hidden">
						<a href={`/games/${$game.id}/production-plans`}>Production Plans</a>
					</li>
					<li class="md:hidden">
						<a href={`/games/${$game.id}/transport-plans`}>Transport Plans</a>
					</li>

					<li class="md:hidden menu-title">
						<span>Reports</span>
					</li>
					<li class="md:hidden"><a href={`/games/${$game.id}/players`}>Players</a></li>
					<li class="md:hidden"><a href={`/games/${$game.id}/planets`}>Planets</a></li>
					<li class="md:hidden"><a href={`/games/${$game.id}/fleets`}>Fleets</a></li>
					<li class="md:hidden"><a href={`/games/${$game.id}/messages`}>Messages</a></li>
					<li class="md:hidden"><a href={`/games/${$game.id}/battles`}>Battles</a></li>
				{/if}

				<li class="md:hidden menu-title">
					<span>Game</span>
				</li>
				<li>
					<a href={`/games/${$game.id}/race`} class="justify-between">Race</a>
				</li>
				<li>
					<a href={`/games/${$game.id}/techs`} class="justify-between">Techs</a>
				</li>
				<li><div class="divider" /></li>
				<li><a href="/auth/logout">Logout, {$me?.username}</a></li>
			</ul>
		</div>
	</div>
</div>
