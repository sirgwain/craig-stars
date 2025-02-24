<script lang="ts">
	import { type PlayerMessage } from '$lib/types/cs';

	type Props = {
		message: PlayerMessage;
	};

	let { message }: Props = $props();
</script>

{#if message.text}
	{message.text}
{:else}
	No message {message.type} - {message.type}
{/if}
