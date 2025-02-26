// The status of the user
// This is used by the login redirect process to determine if the user has been loaded

import { RoleAdmin, RoleGuest, RoleUser, type UserRole } from './cs';

export type UserStatus = (typeof UserStatuses)[keyof typeof UserStatuses];

export const UserStatuses = {
	Unknown: 'Unknown',
	LoggedIn: 'LoggedIn',
	NotFound: 'NotFound'
} as const;

export type SessionUser = {
	id?: number;
	username: string;
	password?: string; // guest users have a hash for the password
	role: UserRole;
	status: UserStatus;
	discordId?: string;
	discordAvatar?: string;
};

export class User implements SessionUser {
	id = 0;
	createdAt = '';
	updatedAt = '';
	username = '';
	password = '';
	role: UserRole = RoleGuest;
	status: UserStatus = UserStatuses.Unknown;
	discordId = '';
	discordAvatar = '';
	lastLogin = '';

	isGuest() {
		return this.role == RoleGuest;
	}

	isAdmin() {
		return this.role == RoleAdmin;
	}
}

export const emptyUser = Object.assign(new User(), {
	username: '',
	role: RoleUser,
	status: UserStatuses.Unknown
});

export const userNotFound = Object.assign(new User(), {
	username: '',
	role: RoleUser,
	status: UserStatuses.NotFound
});
