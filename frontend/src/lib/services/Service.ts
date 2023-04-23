import type { ErrorResponse } from '$lib/types/ErrorResponse';

export abstract class Service {
	static async get<T>(url: string, body?: BodyInit): Promise<T> {
		const response = await fetch(url, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			},
			body: body
		});

		if (response.ok) {
			return (await response.json()) as T;
		} else {
			const errorResponse = (await response.json()) as ErrorResponse;
			console.error(errorResponse);
			throw new Error(errorResponse.error);
		}
	}
	static async update<T>(item: T, url: string): Promise<T> {
		const response = await fetch(url, {
			method: 'PUT',
			headers: {
				accept: 'application/json'
			},
			body: JSON.stringify(item)
		});

		if (response.ok) {
			return (await response.json()) as T;
		} else {
			const errorResponse = (await response.json()) as ErrorResponse;
			console.error(errorResponse);
			throw new Error(errorResponse.error);
		}
	}

	static async delete(url: string): Promise<void> {
		const response = await fetch(`${url}`, {
			method: 'DELETE',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			const errorResponse = (await response.json()) as ErrorResponse;
			console.error(errorResponse);
			throw new Error(errorResponse.error);
		}
		return Promise.resolve();
	}
}
