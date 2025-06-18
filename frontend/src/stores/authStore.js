import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import authApi from '../api/auth';

const useAuthStore = create(
    persist(
        (set, get) => ({
            user: null,
            loading: false,
            isAuthenticated: false,

            checkAuth: () => {
                const token = localStorage.getItem('access_token');
                if (token) {
                    try {
                        const payload = JSON.parse(atob(token.split('.')[1]));
                        if (payload.exp * 1000 > Date.now()) {
                            set({
                                user: {
                                    email: payload.email,
                                    name: payload.name,
                                },
                                isAuthenticated: true,
                                loading: false,
                            });
                            return;
                        }
                    } catch (e) {
                        console.error('Ошибка при разборе токена:', e);
                    }
                    localStorage.removeItem('access_token');
                }
                set({ loading: false });
            },

            signIn: async (email, password) => {
                set({ loading: true });
                try {
                    const response = await authApi.signIn(email, password);
                    localStorage.setItem('access_token', response.data.access_token);
                    localStorage.setItem('public_key', response.data.public_key);
                    set({
                        user: { email: response.data.email },
                        isAuthenticated: true,
                        loading: false,
                    });
                    return { success: true };
                } catch (error) {
                    set({ loading: false });
                    return { success: false, error: error.message };
                }
            },

            signUp: async (email, password, name, publicKey) => {
                set({ loading: true });
                try {
                    const response = await authApi.signUp(email, password, name, publicKey);
                    localStorage.setItem('access_token', response.data.access_token);
                    localStorage.setItem('public_key', response.data.public_key);
                    set({ loading: false });
                    return { success: true };
                } catch (error) {
                    set({ loading: false });
                    return { success: false, error: error.message };
                }
            },

            signOut: () => {
                localStorage.removeItem('access_token');
                set({
                    user: null,
                    isAuthenticated: false,
                });
            },
        }),
        {
            name: 'auth-storage',
            partialize: (state) => ({
                user: state.user,
                isAuthenticated: state.isAuthenticated,
            }),
        }
    )
);

export default useAuthStore;