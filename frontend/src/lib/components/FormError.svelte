<script lang="ts">
	import { ExclamationTriangle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { fade } from 'svelte/transition';

	export let error = '';
</script>

{#if error !== ''}
	<div
		class="alert alert-error shadow-lg w-1/2 mx-auto"
		in:fade
		out:fade={{ delay: 5000 }}
		on:introend={(e) => (error = '')}
	>
		<div>
			<Icon src={ExclamationTriangle} size="24" class="hover:stroke-accent" />
			<span>{error}</span>
		</div>
	</div>
{/if}
