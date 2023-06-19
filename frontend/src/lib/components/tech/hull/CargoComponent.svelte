<script lang="ts">
	export let capacity: number;
</script>

<div
	class={`flex flex-col h-full w-full justify-between border border-slate-900 bg-white text-sm text-black`}
>
	<div class="flex-grow whitespace-pre-wrap text-center">Cargo</div>
	<span class="h-[1rem] mt-auto text-center font-bold">{capacity}kT</span>
</div>
