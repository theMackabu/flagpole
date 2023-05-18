const log = {
	info: (object: any) => console.log(JSON.stringify({ time: new Date(), message: object })),
};

export { log };
