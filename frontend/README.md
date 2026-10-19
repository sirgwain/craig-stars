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

Public guides live under `src/routes/docs` as mdsvex `+page.svx` files. Register
each page and its section anchors in `src/lib/docs/navigation.ts`.

Interactive examples go in `src/lib/docs/examples`, wrapped in `DocExample`. They
should reuse the production components with local state only — no API calls.
