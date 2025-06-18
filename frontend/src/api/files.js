import client from './client';

export default {
    getMyFiles() {
        return client.get('/files', {
            headers: {
                'Content-Type': 'application/json',
            }
        });
    },
    getAvailableFiles() {
        return client.get('/available-files', {
            headers: {
                'Content-Type': 'application/json',
            }
        });
    },
    createFile(formData) {
        return client.post('/file', formData, {
            headers: {
                'Content-Type': 'multipart/form-data',
            },
        });
    },
    shareFile(fileUUID, shareData) {
        return client.post(`/files/${fileUUID}`, shareData);
    },
    deleteFile(fileUUID) {
        return client.delete(`/files/${fileUUID}`);
    },
    deleteFileAccess(fileUUID, data) {
        return client.post(`/file/${fileUUID}/access`, data);
    },
    downloadFile(fileUUID) {
        return client.post(`/download/files/${fileUUID}`, null,{
            responseType: 'arraybuffer',
            headers: {
                'Accept': 'application/octet-stream'
            }
        });
    },
    downloadCommonFile(fileUUID) {
        return client.post(`/download/common/files/${fileUUID}`, null,{
            responseType: 'arraybuffer',
            headers: {
                'Accept': 'application/octet-stream'
            }
        });
    },
};