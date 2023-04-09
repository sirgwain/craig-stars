<script lang="ts">
	export let capacity: number;
</script>

<div class={`flex flex-col h-full w-full justify-between border border-slate-900 bg-white text-sm`}>
	<div class="text-center">Cargo</div>
	<div class="text-center font-bold">{capacity}kT</div>
</div>
