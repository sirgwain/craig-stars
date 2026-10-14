<script lang="ts">
	import { getGameContext } from '#lib/services/GameContext.js';
	import { PlayerMailSchema, type PlayerMail } from '#lib/types/cs-proto.js';
	import { create } from '@bufbuild/protobuf';
	import { ConnectError } from '@connectrpc/connect';
	import hotkeys from 'hotkeys-js';
	import { onMount, untrack } from 'svelte';

	const { game, player, universe, readOnly, getPlayerMail, savePlayerMail, deletePlayerMail } =
		getGameContext();

	// the server rejects longer mail
	const maxLength = 1000;
	const everybody = 0;

	type Props = {
		// who to write to, 0 for everybody
		toPlayerNum: number;
		onClose: () => void;
	};

	let { toPlayerNum, onClose }: Props = $props();

	let to = $state(untrack(() => toPlayerNum));
	let text = $state('');
	// the undelivered mail being changed, a new mail is written otherwise
	let editing: PlayerMail | undefined = $state();
	let mail: PlayerMail[] = $state([]);
	let saving = $state(false);
	let error = $state('');
	let textarea: HTMLTextAreaElement | undefined = $state();

	// the newest mail is first
	let sortedMail = $derived(
		[...mail].sort((a, b) => Number((b.gameDbObject?.id ?? 0n) - (a.gameDbObject?.id ?? 0n)))
	);

	// a player we haven't met has no name yet
	function playerName(num: number): string {
		if (num === everybody) {
			return 'Everybody';
		}
		const intel = $universe.getPlayerIntel(num);
		return intel?.racePluralName || intel?.name || `Player ${num}`;
	}

	// mail is delivered by the turn after the year it was written in
	function isPending(m: PlayerMail): boolean {
		return m.fromPlayerNum === $player.num && m.year >= $game.year;
	}

	async function run(action: () => Promise<void>) {
		saving = true;
		error = '';
		try {
			await action();
		} catch (e) {
			error = e instanceof ConnectError ? e.rawMessage : 'Failed to save message';
		} finally {
			saving = false;
		}
	}

	function reset() {
		editing = undefined;
		text = '';
	}

	function save() {
		if (saving || $readOnly || !text.trim()) {
			return;
		}
		const sent = text;
		run(async () => {
			const saved = await savePlayerMail(
				create(PlayerMailSchema, { gameDbObject: editing?.gameDbObject, toPlayerNum: to, text })
			);
			const id = saved.gameDbObject?.id;
			mail = [...mail.filter((m) => m.gameDbObject?.id !== id), saved];
			editing = undefined;
			// keep anything typed while the mail was saving
			if (text === sent) {
				text = '';
			}
		});
	}

	function remove(m: PlayerMail) {
		if (saving || $readOnly) {
			return;
		}
		run(async () => {
			await deletePlayerMail(m);
			mail = mail.filter((other) => other.gameDbObject?.id !== m.gameDbObject?.id);
			if (editing?.gameDbObject?.id === m.gameDbObject?.id) {
				reset();
			}
		});
	}

	function edit(m: PlayerMail) {
		editing = m;
		to = m.toPlayerNum;
		text = m.text;
		textarea?.focus();
	}

	function onKeyDown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			e.preventDefault();
			onClose();
		} else if (e.key === 'Enter' && (e.metaKey || e.ctrlKey)) {
			e.preventDefault();
			save();
		}
	}

	onMount(() => {
		// keep the game's hotkeys from firing while this dialog is open
		const originalScope = hotkeys.getScope();
		const scope = 'mail';
		hotkeys('Esc', scope, onClose);
		hotkeys.setScope(scope);
		textarea?.focus();

		getPlayerMail()
			.then((loaded) => (mail = loaded))
			.catch(() => (error = 'Failed to load messages'));

		return () => {
			hotkeys.unbind('Esc', scope, onClose);
			hotkeys.deleteScope(scope);
			hotkeys.setScope(originalScope);
		};
	});
</script>

<div class="flex flex-col gap-2" data-type="mail" role="presentation" onkeydown={onKeyDown}>
	<h2 class="text-lg font-semibold text-secondary truncate">Messages</h2>

	{#if error}
		<div class="alert alert-error py-2" data-type="mail-error">{error}</div>
	{/if}

	{#if !$readOnly}
		<label class="flex flex-row items-center gap-2">
			<span class="font-semibold">To</span>
			<select class="select select-sm grow" bind:value={to} data-type="mail-to">
				<option value={everybody}>{playerName(everybody)}</option>
				{#each $universe.playerIntels as intel (intel.num)}
					<option value={intel.num}>{playerName(intel.num)}</option>
				{/each}
			</select>
		</label>

		<textarea
			bind:this={textarea}
			bind:value={text}
			maxlength={maxLength}
			rows="3"
			aria-label="Message"
			placeholder="Delivered when the next turn is generated"
			class="textarea textarea-bordered w-full text-base"></textarea>
	{/if}

	<div class="flex flex-row gap-2">
		{#if editing}
			<button type="button" class="btn btn-sm btn-outline" onclick={reset} data-type="mail-new"
				>New Message</button
			>
		{/if}
		<div class="grow"></div>
		<button type="button" class="btn btn-sm btn-outline" onclick={onClose} data-type="mail-close"
			>Close</button
		>
		{#if !$readOnly}
			<button
				type="button"
				class="btn btn-sm btn-primary"
				disabled={saving || !text.trim()}
				onclick={save}
				data-type="mail-save">{editing ? 'Save' : 'Send'}</button
			>
		{/if}
	</div>

	{#if sortedMail.length}
		<ul class="flex flex-col gap-2 max-h-64 overflow-y-auto" data-type="mail-list">
			{#each sortedMail as m (m.gameDbObject?.id)}
				{@const pending = isPending(m)}
				<li
					class="rounded-xs border border-base-300 p-2"
					class:border-accent={editing?.gameDbObject?.id === m.gameDbObject?.id}
					data-type="mail-item"
				>
					<div class="flex flex-row items-center gap-2 text-sm">
						<span class="grow font-semibold">
							{m.fromPlayerNum === $player.num ? 'You' : playerName(m.fromPlayerNum)} to
							{m.toPlayerNum !== $player.num
								? playerName(m.toPlayerNum)
								: m.fromPlayerNum === $player.num
									? 'yourself'
									: 'you'}, {m.year}
						</span>
						{#if pending}
							<span class="badge badge-sm badge-warning">Sends next turn</span>
							{#if !$readOnly}
								<button
									type="button"
									class="btn btn-xs btn-outline"
									disabled={saving}
									onclick={() => edit(m)}
									data-type="mail-edit">Edit</button
								>
								<button
									type="button"
									class="btn btn-xs btn-outline btn-error"
									disabled={saving}
									onclick={() => remove(m)}
									data-type="mail-delete">Delete</button
								>
							{/if}
						{/if}
					</div>
					<div class="whitespace-pre-wrap wrap-break-word select-text">{m.text}</div>
				</li>
			{/each}
		</ul>
	{/if}
</div>
