import http from 'k6/http';

export class API {

    private baseUrl: string;

    constructor(baseUrl: string) {
        this.baseUrl = baseUrl;
    }
    
    get(path: string) {
        return http.get(`${this.baseUrl}${path}`);
    }

    post(path: string, body: any = {}) {
        return http.post(`${this.baseUrl}${path}`, body);
    }

}