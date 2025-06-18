import { useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import filesApi  from '../api/files';
import { toast } from 'react-toastify';
import { saveAs } from "file-saver";

export const FileDownloadRedirect = () => {
    const { id } = useParams();
    const navigate = useNavigate();

    useEffect(() => {
        const downloadFile = async () => {
            try {
                if (!id) {
                    throw new Error('File ID not specified');
                }

                const response = await filesApi.downloadCommonFile(id);

                const contentDisposition = response.headers.get('Content-Disposition');
                let fileName = 'downloaded_file.pdf';

                if (contentDisposition) {
                    let match = contentDisposition.match(/filename\*=(?:UTF-8'')?([^;]+)/i);
                    if (match && match[1]) {
                        fileName = decodeURIComponent(match[1].replace(/['"]/g, '').trim());
                    } else {
                        match = contentDisposition.match(/filename=([^;]+)/i);
                        if (match && match[1]) {
                            fileName = match[1].replace(/['"]/g, '').trim();
                            fileName = fixEncoding(fileName); // ← тут магия
                        }
                    }
                }

                const blob = new Blob([response.data]);

                saveAs(blob, fileName);

                navigate(-1);
                } catch (e) {
                    toast.error('Не удалось загрузить файл');
                    navigate('/storage/my');
                }
            };

            downloadFile();
        }, [id, navigate]);

        return <div>Downloading file...</div>;
    }

const fixEncoding = (str) => {
    try {
        const bytes = new Uint8Array([...str].map(ch => ch.charCodeAt(0)));
        const decoder = new TextDecoder('utf-8');
        return decoder.decode(bytes);
    } catch {
        return str;
    }
};