// The status of the user
// This is used by the login redirect process to determine if the user has been loaded

import { UserRole } from './cs-proto';

export type UserStatus = (typeof UserStatuses)[keyof typeof UserStatuses];

export const UserStatuses = {
	Unknown: 'Unknown',
	LoggedIn: 'LoggedIn',
	NotFound: 'NotFound'
} as const;

export class UserSession {
	id = BigInt(0);
	createdAt = '';
	updatedAt = '';
	username = '';
	password = '';
	role: UserRole = UserRole.GUEST;
	status: UserStatus = UserStatuses.Unknown;
	discordId = '';
	discordAvatar = '';
	lastLogin = '';

	isGuest() {
		return this.role == UserRole.GUEST;
	}

	isAdmin() {
		return this.role == UserRole.ADMIN;
	}
}

export const emptyUser = Object.assign(new UserSession(), {
	username: '',
	role: UserRole.USER,
	status: UserStatuses.Unknown
});

export const userNotFound = Object.assign(new UserSession(), {
	username: '',
	role: UserRole.USER,
	status: UserStatuses.NotFound
});
