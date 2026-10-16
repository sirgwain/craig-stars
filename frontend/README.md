# frontend

Powered by [`SvelteKit`](https://kit.svelte.dev).

## Developing

Once you've created a project and installed dependencies with `npm install` start a development server:

```bash
npm run dev
```

## Building

To build the static site for hosting

```bash
npm run build
```

You can preview the production build with `npm run preview`.

## Player documentation

Public guides live under `src/routes/docs`. Write prose in `+page.svx` files;
mdsvex compiles these at build time and `rehype-slug` adds heading anchors.
Register pages and their section links in `src/lib/docs/navigation.ts`.
The shared docs layout supplies navigation, breadcrumbs, and previous/next links.
Keep `/docs/mcp` at its existing URL.

Put interactive examples in `src/lib/docs/examples` and wrap them in `DocExample`.
Examples should own fresh local models and reuse production controls, defaults,
and calculations. Keep account loading, API mutations, and navigation in live
page wrappers rather than adding a documentation mode to game components.

`GameCreationExample` demonstrates this with `GameCreationForm` and the same
settings factories used by the single-player and host-game pages. The form
remounts its child editors when its settings model is replaced, including Reset.

The docs Playwright tests cover anonymous use, local preview/reset, navigation,
and phone layouts. They attach screenshots to the existing HTML test report,
which CI uploads as the `playwright-report` artifact; screenshots are not source
files and should not be committed.
