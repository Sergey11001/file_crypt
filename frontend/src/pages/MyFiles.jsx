import { useState, useEffect } from 'react';
import filesApi from '../api/files';
import FileItem from '../components/FileItem';
import { useNavigate } from 'react-router-dom';
import { formatBytes, formatDate } from '../utils/helpers';
import useAuthStore from "../stores/authStore"
import { saveAs } from "file-saver";
import {decryptFile, decryptWithPrivateKey, encryptWithPublicKey} from "../utils/crypto";
import Spinner from "../components/Spinner";
import { toast } from 'react-toastify';
import usePrivateKeyStore from "../stores/privateKeyStore";

export default function MyFiles() {
    const { isAuthenticated } = useAuthStore();
    const navigate = useNavigate();
    const { privateKey } = usePrivateKeyStore();

    useEffect(() => {
        if (!isAuthenticated) {
            navigate("/signin");
        }
    }, [isAuthenticated, navigate]);

    const [files, setFiles] = useState([]);
    const [loading, setLoading] = useState(true);
    const [downloadingFile, setDownloadingFile] = useState(null);

    const fetchFiles = async () => {
        try {
            const response = await filesApi.getMyFiles();
            setFiles(response.data.files);
        } catch (err) {
            toast.warn("Ошибка загрузки файла");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchFiles();
    }, []);

    const handleDownload = async (file, encryptedAesKey, privateKey) => {
        if (!privateKey) {
            toast.warn('Пожалуйста, сначала загрузите свой закрытый ключ, чтобы расшифровать файлы.');
            return;
        }

        setDownloadingFile(file.id);
        try {
            let response
            if (!encryptedAesKey) {
                response = await filesApi.downloadCommonFile(file.id);
            }else {
                response = await filesApi.downloadFile(file.id);
            }

            if (!response.data || !(response.data instanceof ArrayBuffer)) {
                throw new Error(`Invalid data type: ${response.data?.constructor?.name}`);
            }

            let blob
            if (!encryptedAesKey) {
                blob = new Blob([response.data]);
            } else {
                blob = await decryptFile(
                    response.data,
                    encryptedAesKey,
                    privateKey
                );
            }

            saveAs(blob, file.name);
        } catch (error) {
            toast.error('Не удалось загрузить файл. Попробуйте еще раз.');
        } finally {
            setDownloadingFile(null);
        }
    }

    const handleDelete = async (fileUuid) => {
        try {
            await filesApi.deleteFile(fileUuid);
            await fetchFiles();
            toast.success('Файл успешно удален.');
        } catch (err) {
            toast.error('Не удалось удалить файл.');
        }
    };

    const handleShare = async (fileUuid, encryptedAesKey, recipientPublicKey, recipientUuid) => {
        const symmetricKey = await decryptWithPrivateKey(encryptedAesKey, privateKey);

        encryptedAesKey = await encryptWithPublicKey(symmetricKey, recipientPublicKey);

        await filesApi.shareFile(fileUuid, {
            recipient_uuid: recipientUuid,
            symmetric_key: encryptedAesKey,
        });
    };

    if (loading) return (
        <div className="loading-container">
            <Spinner size="large" color="primary" />
            <p>Loading files...</p>
        </div>
    );

    return (
        <div className="files-list">
            {!privateKey && (
                <div className="key-warning">
                    <span>⚠️ Вам необходимо загрузить свой закрытый ключ для расшифровки файлов.</span>
                </div>
            )}

            {files.length > 0 ? (
                files.map((file) => (
                    <FileItem
                        key={file.uuid}
                        file={{
                            id: file.uuid,
                            name: file.name,
                            size: formatBytes(file.size),
                            date: formatDate(file.created_at),
                            encryptedKey: file.symmetric_key
                        }}
                        onDownload={() => handleDownload({
                            id: file.uuid,
                            name: file.name,
                        }, file.symmetric_key, privateKey)}
                        onShare={handleShare}
                        onDelete={handleDelete}
                        isOwner={true}
                        isDownloading={downloadingFile === file.uuid}
                        isDisabled={!privateKey}
                    />
                ))
            ) : (
                <p className="empty-message">У вас нет никаких файлов.</p>
            )}
        </div>
    );
}