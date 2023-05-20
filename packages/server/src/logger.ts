import bunyan from 'bunyan';

const log = bunyan.createLogger({ name: 'server' });
const time = (start: number) => Date.now() - start;

const getPathFromURL = (url: string, strict: boolean = true): string => {
	const queryIndex = url.indexOf('?', 8);
	const result = url.substring(url.indexOf('/', 8), queryIndex === -1 ? url.length : queryIndex);

	if (strict === false && /.+\/$/.test(result)) {
		return result.slice(0, -1);
	}

	return result;
};

const logger = () => {
	return async (c, next) => {
		const { method } = c.req;
		log.info({ method, status: c.res.status }, '[HTTP REQUEST - START]');

		const start = Date.now();
		await next();

		log.info({ method, status: c.res.status, route: getPathFromURL(c.req.url), duration: time(start) }, '[HTTP REQUEST - END]');
	};
};

export { log, logger };
