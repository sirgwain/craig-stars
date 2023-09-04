import { User, userNotFound, UserStatus, type SessionUser } from '$lib/types/User';
import { me } from './services/Stores';

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
		const sessionUser = (await response.json()) as SessionUser;
		const user = Object.assign(new User(), sessionUser)

		user.status = UserStatus.LoggedIn;
		me.update(() => user);

		return user;
	}
}
