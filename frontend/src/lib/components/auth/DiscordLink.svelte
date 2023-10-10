<script lang="ts">
	import Discord from '../icons/Discord.svelte';

	export let from = '/';
</script>

<a class="btn bg-[#7289da]" href={`/api/auth/discord/login?from=${from}`}
	>Login&nbsp;
	<Discord class="fill-base-content w-4 h-4" />
</a>
