<script lang="ts">
	import fuzzy from 'fuzzy';
	import { tick, createEventDispatcher, afterUpdate, onMount } from 'svelte';

	type Result = {
		original: string;
		index: number;
		score: number;
		string: string;
	};

	export let id = 'typeahead-' + Math.random().toString(36);
	export let wrapperClass = '';
	export let value = '';

	export let data: string[] = [];

	/** Set to `false` to prevent the first result from being selected */
	export let autoselect = true;

	export let results: Result[] = [];

	/** Set to `true` to re-focus the input after selecting a result */
	export let focusAfterSelect = true;

	/**
	 * Specify the maximum number of results to return
	 * @type {number}
	 */
	export let limit = Infinity;

	const dispatch = createEventDispatcher();

	let comboboxRef: HTMLDivElement | null = null;
	let inputRef: HTMLInputElement | null = null;
	let hideDropdown = true;
	let selectedIndex = -1;
	let prevResults = '';
	let originalValue = value;

	afterUpdate(() => {
		if (prevResults !== resultsId && autoselect) {
			selectedIndex = 0;
		}

		if (prevResults !== resultsId) {
			hideDropdown = results.length === 0;
		}

		prevResults = resultsId;
	});

	async function select() {
		const result = results[selectedIndex];
		const selectedValue = result.original;

		// update the value
		value = selectedValue;
		originalValue = value;

		dispatch('select', {
			selectedIndex,
			selected: selectedValue,
			original: result.original,
			originalIndex: result.index
		});

		await tick();

		if (focusAfterSelect && inputRef) inputRef.focus();
		close();
	}

	function change(direction: -1 | 1) {
		let index =
			direction === 1 && selectedIndex === results.length - 1 ? 0 : selectedIndex + direction;
		if (index < 0) index = results.length - 1;

		selectedIndex = index;
	}

	const open = () => (hideDropdown = false);
	const close = () => (hideDropdown = true);

	$: results = fuzzy
		.filter(value, data)
		.filter(({ score }) => score > 0)
		.slice(0, limit);
	$: resultsId = results.join('');
	$: showResults = !hideDropdown && results.length > 0;

	const onWindowClick = (event: Event) => {
		if (!comboboxRef?.contains(event.target as Node) && !event.defaultPrevented) {
			close();
		}
	};
</script>

<svelte:window on:click={onWindowClick} />

<div class="dropdown dropdown-top {wrapperClass}" id="{id}-typeahead" bind:this={comboboxRef}>
	<input
		type="text"
		{id}
		{...$$restProps}
		bind:value
		on:change={() => {
			open();
		}}
		on:focus={() => {
			if (results.length > 1) {
				open();
			}
		}}
		on:keydown={(e) => {
			if (results.length === 0) return;

			switch (e.key) {
				case 'Enter':
					e.preventDefault();
					select();
					break;
				case 'ArrowDown':
					e.preventDefault();
					change(1);
					break;
				case 'ArrowUp':
					e.preventDefault();
					change(-1);
					break;
				case 'Escape':
					e.preventDefault();
					e.currentTarget.value = originalValue;
					e.currentTarget.focus();
					close();
					break;
			}
		}}
	/>
	{#if showResults}
		<ul
			class="dropdown-content menu p-2 shadow bg-base-200 border border-base-300 my-2 z-10"
			id="{id}-listbox"
		>
			{#each results as result, index}
				<li
					id="{id}-result-{index}"
					class="cursor-pointer"
					class:text-accent-focus={selectedIndex === index}
					on:mousedown|preventDefault={() => {
						selectedIndex = index;
						select();
					}}
					on:mouseenter|preventDefault={() => {
						selectedIndex = index;
					}}
				>
					<slot {result} {index} {value}>
						{@html result.string}
					</slot>
				</li>
			{/each}
		</ul>
	{/if}
</div>
