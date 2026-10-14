<script module lang="ts">
	export type InnateScannerTooltipProps = {
		// the planet's penetrating scan range, from its starbase
		scanRangePen?: number;
	};
</script>

<script lang="ts">
	let { scanRangePen = 0 }: InnateScannerTooltipProps = $props();
</script>

<div class="flex flex-col sm:w-104 m-auto">
	<div>
		Your race cannot build planetary scanners. Your starbases have an innate scanning ability whose
		distance is equal to the square root of one tenth the starbase population.
	</div>
	{#if scanRangePen > 0}
		<div class="mt-2" data-type="innate-pen-scan">
			Your advanced starbase design also allows you to perform penetrating scans at half the range
			of normal scans.
		</div>
	{/if}
</div>
