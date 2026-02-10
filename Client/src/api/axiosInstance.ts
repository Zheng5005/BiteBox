import axios from 'axios';

const instance = axios.create({
  baseURL: 'http://localhost:8080/api',
  //timeout: 5000, // 5 seconds
  headers: {
    'Content-Type': 'application/json',
  },
});

let on401Callback: (() => void) | null = null;

export const setOn401 = (callback: () => void) => {
  on401Callback = callback;
};

// Add a request interceptor to include the token
instance.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Add a response interceptor to handle 401 errors
instance.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      if (on401Callback) {
        on401Callback();
      }
    }
    return Promise.reject(error);
  }
);

export default instance;
