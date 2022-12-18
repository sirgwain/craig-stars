import { userNotFound, UserStatus, type User } from '$lib/types/User';
import { me } from './services/Context';

export async function authGuard(): Promise<User | undefined> {
	const response = await fetch(`/api/me`, {
		method: 'GET',
		headers: {
			accept: 'application/json'
		}
	});

	if (!response.ok) {
		// no user
		me.update(() => userNotFound);
	} else {
		// update the logged in user in the context
		const user = (await response.json()) as User;
		user.status = UserStatus.LoggedIn;
		me.update(() => user);

		return user;
	}
}
