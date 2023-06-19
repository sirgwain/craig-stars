import type { Handle } from '@sveltejs/kit';

// disable SSR because Pixi doesn't like doing that
export const handle: Handle = async ({ event, resolve }) => {
	// console.log(event);
	// return resolve(event, { ssr: !event.routeId?.startsWith('games/[id]') });
	let ssr = false;
	// if (event.routeId?.startsWith('techs')) {
	// 	ssr = true;
	// }
	return resolve(event, { ssr });
};
