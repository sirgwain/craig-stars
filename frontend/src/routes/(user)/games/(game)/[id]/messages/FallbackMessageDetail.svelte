<script lang="ts">
	import { MessageType, type Message } from '$lib/types/Message';
	import { $enum as eu } from 'ts-enum-util';

	type Props = {
		message: Message;
	};

	let { message }: Props = $props();
</script>

{#if message.text}
	{message.text}
{:else}
	No message {message.type} - {eu(MessageType).getKeyOrDefault(message.type)}
{/if}
