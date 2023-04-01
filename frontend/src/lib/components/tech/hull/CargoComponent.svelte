<script lang="ts">
	export let capacity: number;
	export let width: number | undefined = undefined;
	export let height: number | undefined = undefined;

	$: {
		if (!width) {
			width = capacity < 100 ? 64 : capacity < 1000 ? 128 : 192;
		}
		if (!height) {
			height = capacity < 500 ? 64 : 128;
		}
	}
</script>

<div
	class={`flex flex-col justify-between border border-slate-900 bg-white text-sm`}
	style={`width: ${width}px; height: ${height}px`}
>
	<div class="text-center">Cargo</div>
	<div class="text-center font-bold">{capacity}kT</div>
</div>
