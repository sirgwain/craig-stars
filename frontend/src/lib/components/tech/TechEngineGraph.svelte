<script lang="ts">
	
	import type { TechEngine } from '$lib/types/Tech';
	import type { ChartData,ChartOptions } from 'chart.js';
	import { onMount } from 'svelte';
	import Line from 'svelte-chartjs/src/Line.svelte';
	// import { EngineFuelUsageScale } from '$lib/EngineFuelUsageScale';

	export let engine: TechEngine;

	let dataLine: ChartData = {
		labels: [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
		datasets: [
			{
				label: 'Fuel Usage %',
				fill: true,
				backgroundColor: 'rgba(225, 204, 230, .3)',
				borderColor: 'rgb(0, 146, 214)',
				borderCapStyle: 'butt',
				borderDash: [],
				borderDashOffset: 0.0,
				borderJoinStyle: 'miter',
				pointBorderColor: 'rgb(194, 236, 255)',
				pointBackgroundColor: 'rgb(255, 255, 255)',
				pointBorderWidth: 10,
				pointHoverRadius: 5,
				pointHoverBackgroundColor: 'rgb(0, 0, 0)',
				pointHoverBorderColor: 'rgba(220, 220, 220,1)',
				pointHoverBorderWidth: 2,
				pointRadius: 1,
				pointHitRadius: 10,
				data: engine?.fuelUsage ? [...engine.fuelUsage] : []
			}
		]
	};

	const options: ChartOptions = {
		responsive: true,
		maintainAspectRatio: false,
		scales: {
			yAxis: {
				type: 'linear',
				min: 0,
				max: 1200,
				ticks: {
					callback: (value) => value + '%'
				}
			}
		},
		plugins: {
			tooltip: {
				displayColors: false,
				callbacks: {
					label: (context) => context.raw + '%'
				}
			}
		}
	};

	onMount(() => {});
</script>

<Line class="h-full w-full" data={dataLine} {options} />
