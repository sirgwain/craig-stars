<script lang="ts">
	import type { SVGAttributes } from 'svelte/elements';

	type Props = SVGAttributes<SVGElement>;

	let { ...rest }: Props = $props();
</script>

<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" {...rest}>
	<!-- Short metal horns give the mine its distinctive silhouette at toolbar size. -->
	<g stroke="#686557" stroke-width="2" stroke-linecap="square">
		<path d="M12 2.5v4M12 17.5v4M2.5 12h4M17.5 12h4" />
		<path d="m5 5 3 3m8 8 3 3M5 19l3-3M16 8l3-3" />
	</g>
	<path
		d="M10.5 2h3M10.5 22h3M2 10.5v3M22 10.5v3M4 6l2-2M18 20l2-2M4 18l2 2M18 4l2 2"
		stroke="#928b71"
		stroke-width="1.5"
		stroke-linecap="square"
	/>
	<circle cx="12" cy="12" r="6.5" fill="#9b8c43" stroke="#514f43" stroke-width="1.5" />
	<path
		d="M7.5 11A4.7 4.7 0 0 1 11 7.5"
		stroke="#e1c45a"
		stroke-width="1.5"
		stroke-linecap="round"
	/>
	<path
		d="M16.5 13a4.7 4.7 0 0 1-3.5 3.5"
		stroke="#716634"
		stroke-width="1.5"
		stroke-linecap="round"
	/>
	<g fill="#f5d44a" stroke="#635b35" stroke-width="0.5">
		<circle cx="12" cy="9" r="0.9" />
		<circle cx="9" cy="12" r="0.9" />
		<circle cx="15" cy="12" r="0.9" />
		<circle cx="12" cy="15" r="0.9" />
	</g>
	<circle cx="12" cy="12" r="1.4" fill="#625e40" />
	<circle cx="11.6" cy="11.6" r="0.5" fill="#f5d44a" />
</svg>
