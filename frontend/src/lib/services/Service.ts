export abstract class Service {
	protected async update<T>(item: T, url: string): Promise<T> {
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
			console.error(response);
		}
		return Promise.resolve(item);
	}

}
