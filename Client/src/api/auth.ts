import axiosInstance from './axiosInstance';

export function login(email: string, password: string) {
  return axiosInstance.post<{ token: string }>('/auth/login', { email, password });
}

export function signup(formData: FormData) {
  return axiosInstance.postForm('/auth/signup', formData);
}
