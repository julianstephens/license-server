export namespace model {

	export class APIKey {
	    id?: string;
	    created_at?: number;
	    updated_at?: number;
	    deleted_at?: number;
	    user_id: string;
	    key: number[];
	    expires_at: number;
	    authentication_scopes: string;

	    static createFrom(source: any = {}) {
	        return new APIKey(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	        this.deleted_at = source["deleted_at"];
	        this.user_id = source["user_id"];
	        this.key = source["key"];
	        this.expires_at = source["expires_at"];
	        this.authentication_scopes = source["authentication_scopes"];
	    }
	}
	export class AuthRequest {
	    name: string;
	    email: string;
	    password: string;

	    static createFrom(source: any = {}) {
	        return new AuthRequest(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.email = source["email"];
	        this.password = source["password"];
	    }
	}
	export class DisplayAPIKey {
	    id?: string;
	    created_at?: number;
	    updated_at?: number;
	    deleted_at?: number;
	    user_id: string;
	    key: string;
	    expires_at: string;
	    authentication_scopes: string;

	    static createFrom(source: any = {}) {
	        return new DisplayAPIKey(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	        this.deleted_at = source["deleted_at"];
	        this.user_id = source["user_id"];
	        this.key = source["key"];
	        this.expires_at = source["expires_at"];
	        this.authentication_scopes = source["authentication_scopes"];
	    }
	}
	export class User {
	    id?: string;
	    created_at?: number;
	    updated_at?: number;
	    deleted_at?: number;
	    name?: string;
	    email?: string;

	    static createFrom(source: any = {}) {
	        return new User(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	        this.deleted_at = source["deleted_at"];
	        this.name = source["name"];
	        this.email = source["email"];
	    }
	}

}
