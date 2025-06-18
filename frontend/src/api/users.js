import client from './client';

export default {
    getUsers() {
        return client.get('/users', {
            headers: {
                'Content-Type': 'application/json',
            }
        });
    },
    getSharedUsers(fileUUID){
        return client.get(`/users/available/${fileUUID}`)
    },
    getUsersForShare(fileUUID) {
        return client.get(`/users/for-share/${fileUUID}`);
    },
    updateUserKeys(data) {
        return client.post(`/users/update-keys`, data)
    }
};