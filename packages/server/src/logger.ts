const log = {
	info: (object: any) => console.log(JSON.stringify({ level: 'info', time: new Date(), message: object })),
	error: (object: any) => console.error(JSON.stringify({ level: 'error', time: new Date(), message: object })),
};

export { log };
