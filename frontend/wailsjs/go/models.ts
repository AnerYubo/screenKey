export namespace knock {
	
	export class KnockData {
	    id: number;
	    host: string;
	    targetPort: number;
	    knockPorts: number[];
	    remark: string;
	
	    static createFrom(source: any = {}) {
	        return new KnockData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.host = source["host"];
	        this.targetPort = source["targetPort"];
	        this.knockPorts = source["knockPorts"];
	        this.remark = source["remark"];
	    }
	}

}

export namespace totp {
	
	export class GenerateTOTPRequest {
	    issuer: string;
	    account: string;
	    remark: string;
	    category: string;
	
	    static createFrom(source: any = {}) {
	        return new GenerateTOTPRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.issuer = source["issuer"];
	        this.account = source["account"];
	        this.remark = source["remark"];
	        this.category = source["category"];
	    }
	}
	export class TOTPData {
	    id: number;
	    secret: string;
	    otpauth: string;
	    issuer: string;
	    account: string;
	    remark: string;
	    category: string;
	
	    static createFrom(source: any = {}) {
	        return new TOTPData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.secret = source["secret"];
	        this.otpauth = source["otpauth"];
	        this.issuer = source["issuer"];
	        this.account = source["account"];
	        this.remark = source["remark"];
	        this.category = source["category"];
	    }
	}

}

