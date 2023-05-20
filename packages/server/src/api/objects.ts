import { parseArgs } from '../parse';
import { version } from '../../package.json';

const args = parseArgs(process.argv.slice(2));

const notFound = {
	error: {
		code: 404,
		message: 'not found',
	},
};

const urls = {
	login: `http://${args.auth}/_/login`,
	create: `http://${args.auth}/_/signup`,
	refresh: `http://${args.auth}/_/refresh`,
	database: `http://${args.db_url}/${args.database}`,
	list: `http://${args.db_url}/list`,
};

const agent = {
	headers: {
		'user-agent': `flagpole_server/v${version}`,
	},
};

const config = (body, plain) => {
	return plain
		? {
				body: body,
				throwHttpErrors: false,
				...agent,
		  }
		: {
				json: body,
				throwHttpErrors: false,
				...agent,
		  };
};

export { notFound, urls, agent, config, args };
