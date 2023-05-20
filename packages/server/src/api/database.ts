import got from 'got';
import { log } from '../logger';
import { Cache } from '../cache';
import { throttle } from '../helpers';
import { urls, agent } from './objects';

const cache = new Cache();

const fetch = throttle(async (data) => {
	const { address, builder } = data;
	const response = await got
		.get(address, { ...agent, throwHttpErrors: false })
		.text()
		.catch((err) => log.error(err));

	cache.set(builder, response);
}, 10 * 1000);

const cacheHandler = (id, key) => {
	const address = urls.database + `/${id}.flags.${key}`;
	const builder = id + '.flags.' + key;

	fetch({ address, builder });
	return cache.get(builder);
};

export { fetch, cache, cacheHandler };
