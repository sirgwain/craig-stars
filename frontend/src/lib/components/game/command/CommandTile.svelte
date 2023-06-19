<script lang="ts">
	export let title = '';
</script>

<div class="p-2 border-2 border-slate-700">
	<h2 class="text-center pb-1">{title}</h2>
	<slot />
</div>
