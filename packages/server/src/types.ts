interface Flags {
	bools: { [key: string]: boolean };
	strings: { [key: string]: boolean };
	unknownFn: ((arg: any) => void | boolean) | null;
	allBools?: boolean;
}

interface Options {
	unknown?: () => boolean;
	boolean?: boolean | string[];
	alias?: { [key: string]: string[] };
	string?: string[];
	default?: { [key: string]: any };
	stopEarly?: boolean;
	'--'?: boolean;
}

interface ParsedArgs {
	[key: string]: any;
	_: any[];
	'--'?: string[];
}

export type { Flags, Options, ParsedArgs };
