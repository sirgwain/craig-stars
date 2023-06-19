<script lang="ts">
	export let title: string | undefined = '';
</script>

<div class="text-xl breadcrumbs">
	<ul>
		<li><a class="link-primary" href="/tests">Tests</a></li>
		{#if title}
			<li>{title}</li>
		{/if}
	</ul>
</div>
