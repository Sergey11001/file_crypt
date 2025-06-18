import { useState, useEffect } from 'react';
import FileItem from '../components/FileItem';
import filesApi from "../api/files";
import {formatBytes, formatDate} from "../utils/helpers";
import {decryptFile} from "../utils/crypto";
import {saveAs} from "file-saver";
import {toast} from "react-toastify";
import usePrivateKeyStore from "../stores/privateKeyStore";

export default function SharedFiles() {
    const [files, setFiles] = useState([]);
    const [loading, setLoading] = useState(true);
    const { privateKey } = usePrivateKeyStore();

    useEffect(() => {
        const fetchFiles = async () => {
            try {
                const response = await filesApi.getAvailableFiles();
                setFiles(response.data.files);
            } catch (error) {
                toast.error('Ошибка загрузки файлов:', error);
            } finally {
                setLoading(false);
            }
        };

        fetchFiles();
    }, []);

    const handleDownload = async (file, encryptedAesKey, privateKeyPem) => {
        try {
            const response = await filesApi.downloadFile(file.id);

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
                    privateKeyPem
                );
            }

            saveAs(blob, file.name);
        } catch (error) {
            toast.error("Ошибка при скачивании файла.");
            throw error;
        }
    }

    if (loading) return <div>Загрузка...</div>;

    return (
        <div className="files-list">
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
                        isOwner={false}
                        onDownload={() => handleDownload({
                            id: file.uuid,
                            name: file.name,
                        }, file.symmetric_key, privateKey)}
                    />
                ))
            ) : (
                <p>Вам пока не предоставили доступ к файлам</p>
            )}
        </div>
    );
}