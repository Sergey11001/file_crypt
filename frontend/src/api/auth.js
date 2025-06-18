import client from './client';

export default {
    async signIn(email, password) {
        return client.post('/sign-in', { email, password });
    },

    async signUp(email, password, name, publicKey) {
        return client.post('/sign-up', { email, password, name, public_key: publicKey });
    },
};