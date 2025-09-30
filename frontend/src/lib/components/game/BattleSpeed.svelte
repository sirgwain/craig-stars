<script lang="ts">
	type Props = {
		speed: number;
	};

	let { speed }: Props = $props();

	let movesPerRoundBase = $derived(Math.floor(speed / 4));
</script>

{#if speed % 4 === 2}
	{movesPerRoundBase ? movesPerRoundBase : ''}&frac12;
{:else if speed % 4 === 1}
	{movesPerRoundBase ? movesPerRoundBase : ''}&frac14;
{:else if speed % 4 === 3}
	{movesPerRoundBase ? movesPerRoundBase : ''}&frac34;
{:else}
	{movesPerRoundBase ? movesPerRoundBase : '-'}
{/if}
