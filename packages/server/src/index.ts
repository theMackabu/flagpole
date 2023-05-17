import { serve } from '@hono/node-server';
import { Hono } from 'hono';
import { log } from './logger';
import got from 'got';
import { jwt } from 'hono/jwt';

const app = new Hono();
const base = 'http://127.0.0.1:5000';

const urls = {
	login: `${base}/api/login`,
	create: `${base}/api/signup`,
};

const config = (body) => ({ json: body, throwHttpErrors: false });

app.use(
	'/api/*',
	jwt({
		secret: '2iuhygw3efdrtgyuewi2okedmf',
	})
);

app.post('/auth/login', async (c) => {
	const body = await c.req.json();
	const response = await got
		.post(urls.login, config(body))
		.json()
		.catch((err) => log.info(err));

	response.error != null && c.status(401);
	return c.json(response);
});

app.post('/api/user/create', async (c) => {
	const body = await c.req.json();
	const response = await got
		.post(urls.create, config(body))
		.json()
		.catch((err) => log.info(err));

	response.created == false && c.status(400);
	return c.json(response);
});

log.info('started on port 6000');
serve({
	fetch: app.fetch,
	port: 6000,
});
