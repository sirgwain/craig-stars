<script lang="ts">
</script>

<!-- ScannerMineFieldPattern -->
<!-- This is a pattern for the background of minefields -->
<svg>
	<defs>
		<pattern id="pattern-minefield" x="0" y="0" width="8" height="8" patternUnits="userSpaceOnUse">
			<rect x="4" y="0" width="1" height="1" fill="#FFFFFF" />
			<rect x="2" y="2" width="1" height="1" fill="#FFFFFF" />
			<rect x="6" y="2" width="1" height="1" fill="#FFFFFF" />
			<rect x="0" y="4" width="1" height="1" fill="#FFFFFF" />
			<rect x="2" y="6" width="1" height="1" fill="#FFFFFF" />
			<rect x="6" y="6" width="1" height="1" fill="#FFFFFF" />
		</pattern>

		<mask id="mask-minefield" x="0" y="0" width="1" height="1">
			<rect x="0" y="0" width="10000" height="10000" fill="url(#pattern-minefield)" />
		</mask>
	</defs>
</svg>
