<script lang="ts">
	import Login from './auth/Login.svelte';
	import Discord from './icons/Discord.svelte';
	import GitHub from './icons/GitHub.svelte';
</script>

<div class="hero min-h-screen">
	<div class="hero-content">
		<div class="max-w-2xl bg-base-200 p-5 border border-base-300">
			<div class="justify-end flex flex-row gap-1">
				<a
					href="https://discord.gg/Ctdx7h6UZS"
					class="btn btn-outline btn-sm btn-secondary rounded-md"
					><Discord class="fill-base-content w-5 h-5" /></a
				>
				<a
					href="https://github.com/sirgwain/craig-stars"
					class="btn btn-outline btn-sm btn-secondary rounded-md"
					><GitHub class="fill-base-content w-5 h-5" /></a
				>
			</div>
			<h1 class="text-center text-5xl font-bold text-accent">craig-stars</h1>
			<div class="flex flex-col md:flex-row py-6">
				<div>
					A clone of <a href="https://wiki.starsautohost.org/" class="cs-link">Stars!</a>, an interstellar, turn based strategy game and a member of the 4<span class="text-success"
						>X</span
					>
					genre (e<span class="text-success">X</span>plore, e<span class="text-success">X</span
					>pand, e<span class="text-success">X</span>ploit, e<span class="text-success">X</span
					>terminate).
				</div>
				<div class="md:grow w-full">
					<div class="text-center py-2 sm:py-0">
						<Login />
					</div>
				</div>
			</div>
			<div class="border border-primary">
				<img src="scanner.png" alt="screenshot" />
			</div>
		</div>
	</div>
</div>
