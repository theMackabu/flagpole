import { serve } from '@hono/node-server';
import { Hono } from 'hono';
import { log } from './logger';
import got from 'got';
import { jwt } from 'hono/jwt';
import { parse } from './cli';
import { logger } from 'hono/logger';

const app = new Hono();
const args = parse(process.argv.slice(2));

const urls = {
	login: `http://${args.auth}/_/login`,
	create: `http://${args.auth}/_/signup`,
	refresh: `http://${args.auth}/_/refresh`,
};

const config = (body) => ({ json: body, throwHttpErrors: false });

app.use('*', logger());
app.use('/api/*', jwt({ secret: args.secret }));

app.post('/auth/login', async (c) => {
	const body = await c.req.json();
	const response = await got
		.post(urls.login, config(body))
		.json()
		.catch((err) => log.info(err));

	response.error != null && c.status(401);
	return c.json(response);
});

app.post('/api/user/refresh', async (c) => {
	const body = await c.req.json();
	const token = c.req.header('Authorization');

	const response = await got
		.post(urls.refresh, {
			...config(body),
			headers: {
				Authorization: token,
			},
		})
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

log.info(`started on port ${args.port}`);
serve({
	fetch: app.fetch,
	port: args.port,
});
