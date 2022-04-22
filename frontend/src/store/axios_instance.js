import axios from "axios"
import {router} from "../main"

const instance = axios.create({
    baseURL: 'https://backend.nilswilhelm.net/',
    // baseURL: 'http://localhost:9000/',
    timeout: 10000,
    params: {} // do not remove this, its added to add params later in the config
});

// Add a request interceptor
instance.interceptors.request.use(function (config) {
    const token = localStorage.getItem("token")
    const userId = localStorage.getItem("userId")
    config.headers.Authorization = token;
    config.headers.userId = userId;

    return config;
});

instance.interceptors.response.use(res => {
    console.log("Hallo");

    // Important: response interceptors **must** return the response.
    return res;
});

instance.interceptors.response.use(
    res => res,
    err => {
        if (err.response.status === 401) {
            console.log("401")
            router.push("/login")
        }
        return  err;
    }
);

export default instance