import type { User } from '$lib/types/User';

export async function authGuard(): Promise<User | undefined> {
	const response = await fetch(`/api/me`, {
		method: 'GET',
		headers: {
			accept: 'application/json'
		}
	});

	if (!response.ok) {
		document.location = '/login';
	} else {
		return (await response.json()) as User;
	}
}
