import axios from 'axios';

const instance = axios.create({
  baseURL: 'http://localhost:8080/api',
  //timeout: 5000, // 5 seconds
  headers: {
    'Content-Type': 'application/json',
  },
});

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

export default instance;
