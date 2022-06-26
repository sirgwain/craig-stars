<script lang="ts">
	export let title = '';
</script>

<div class="card bg-base-200 shadow-xl mb-1">
	<div class="card-body p-4 text-sm">
		<h2 class="card-title">{title}</h2>
		<slot />
	</div>
</div>
