export function fetchJWT(user: string): Promise<string> {
	return fetch('http://localhost:8080/auth/jwt', {
		body: JSON.stringify({ user }),
		method: 'POST'
	})
		.then(res => {
			const body = res.json();
			return body;
		})
		.then(body => {
			return body.token;
		})
		.catch(err => {
			throw err;
		});
}
