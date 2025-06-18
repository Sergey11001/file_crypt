import {useCallback, useState} from 'react';
import {useDropzone} from 'react-dropzone';
import {useNavigate} from 'react-router-dom';
import {encryptFile, encryptWithPublicKey, generateAesKey,} from '../utils/crypto';
import filesApi from '../api/files';
import UserSelector from "../components/UserSelector";
import {toast} from "react-toastify";

export default function UploadFile() {
    const [file, setFile] = useState(null);
    const [accessType, setAccessType] = useState('private');
    const [recipientEmail, setRecipientEmail] = useState('');
    const [recipientPublicKey, setRecipientPublicKey] = useState(null);
    const [recipientUuid, setRecipientUuid] = useState(null);
    const [isLoading, setIsLoading] = useState(false);
    const [progress, setProgress] = useState(0);
    const navigate = useNavigate();

    const onDrop = useCallback((acceptedFiles) => {
        if (acceptedFiles.length) {
            const selectedFile = acceptedFiles[0];

            if (selectedFile.size > 50 * 1024 * 1024) {
                toast.warn('Размер файла превышает лимит в 50 МБ')
                return;
            }

            setFile(selectedFile);
            setProgress(0);
        }
    }, []);

    const { getRootProps, getInputProps, isDragActive } = useDropzone({
        onDrop,
        maxFiles: 1,
        accept: {
            'application/pdf': ['.pdf'],
            'image/*': ['.png', '.jpg', '.jpeg'],
            'text/plain': ['.txt']
        },
        disabled: isLoading
    });

    const handleSubmit = async (e) => {
        e.preventDefault();

        setIsLoading(true);
        setProgress(0);

        try {
            let symmetricKey = null;
            let encryptedAesKey = null;

            if (accessType === 'private' || accessType === 'specific') {
                setProgress(5);
                symmetricKey = await generateAesKey();
            }

            if (accessType === 'private' || accessType === 'specific') {
                setProgress(15);
                const publicKeyPem = localStorage.getItem('public_key');
                if (!publicKeyPem) throw new Error('Public key not found in localStorage');
                encryptedAesKey = await encryptWithPublicKey(symmetricKey, publicKeyPem);
            }

            setProgress(30);
            const fileToUpload = accessType === 'private' || accessType === 'specific'
                ? await encryptFile(file, symmetricKey)
                : file;

            setProgress(50);
            const formData = new FormData();
            formData.append('file', fileToUpload, file.name);
            formData.append('name', file.name);

            if (accessType === 'private' || accessType === 'specific') {
                formData.append('symmetric_key', encryptedAesKey);
            }

            setProgress(75);
            const response = await filesApi.createFile(formData, (progressEvent) => {
                const percentCompleted = Math.round(
                    (progressEvent.loaded * 100) / progressEvent.total
                );
                setProgress(75 + percentCompleted * 0.25);
            });

            if (accessType === 'specific') {
                encryptedAesKey = await encryptWithPublicKey(symmetricKey, recipientPublicKey);
                await filesApi.shareFile(response.data.file.uuid, {
                    recipient_uuid: recipientUuid,
                    symmetric_key: encryptedAesKey,
                });
            }

            setProgress(100);
            navigate('/storage/my');
        } catch (err) {
            toast.error('Ошибка загрузки файла.');
            setProgress(0);
        } finally {
            setIsLoading(false);
        }
    };

    const removeFile = () => {
        setFile(null);
        setProgress(0);
    };

    return (
        <div className="upload-page">
            <div className="upload-container">
                <h2 className="upload-title">Загрузка файла</h2>
                <p className="upload-description">
                    {accessType === 'private'
                        ? 'Файлы шифруются перед загрузкой. Расшифровать их можете только вы.'
                        : accessType === 'public'
                            ? 'Файлы будут доступны любому, у кого есть ссылка.'
                            : 'Файлы будут зашифрованы для конкретного получателя.'}
                </p>

                <div
                    {...getRootProps()}
                    className={`dropzone ${isDragActive ? 'active' : ''} ${file ? 'has-file' : ''}`}
                >
                    <input {...getInputProps()} />
                    {file ? (
                        <div className="file-preview">
                            <div className="file-icon">
                                {file.type.startsWith('image/') ? (
                                    <i className="icon-image"></i>
                                ) : (
                                    <i className="icon-file"></i>
                                )}
                            </div>
                            <div className="file-info">
                                <p className="file-name">{file.name}</p>
                                <p className="file-size">{(file.size / 1024 / 1024).toFixed(2)} MB</p>
                            </div>
                            <button
                                type="button"
                                className="remove-btn"
                                onClick={(e) => {
                                    e.stopPropagation();
                                    removeFile();
                                }}
                                disabled={isLoading}
                            >
                                &times;
                            </button>
                        </div>
                    ) : (
                        <>
                            <i className="icon-upload"></i>
                            <p className="dropzone-text">
                                {isDragActive ? 'Перетащите файл сюда': 'Перетащите файл или щелкните, чтобы выбрать'}
                            </p>
                            <p className="dropzone-hint">Поддерживает: PDF, JPG, PNG, TXT (макс. 50 МБ)</p>
                        </>
                    )}
                </div>

                <div className="access-controls">
                    <div className="access-options">
                        <label className="access-option">
                            <input
                                type="radio"
                                name="accessType"
                                value="private"
                                checked={accessType === 'private'}
                                onChange={() => setAccessType('private')}
                            />
                            <span className="access-label">
                <i className="icon-lock"></i> Личное
              </span>
                        </label>

                        <label className="access-option">
                            <input
                                type="radio"
                                name="accessType"
                                value="public"
                                checked={accessType === 'public'}
                                onChange={() => setAccessType('public')}
                            />
                            <span className="access-label">
                <i className="icon-globe"></i> Публичный (любой, у кого есть ссылка)
              </span>
                        </label>

                        <label className="access-option">
                            <input
                                type="radio"
                                name="accessType"
                                value="specific"
                                checked={accessType === 'specific'}
                                onChange={() => setAccessType('specific')}
                            />
                            <span className="access-label">
                <i className="icon-user"></i> Конкретный пользователь
              </span>
                        </label>
                    </div>

                    {accessType === 'specific' && (
                        <div className="recipient-section">
                            <label>Выберите получателя:</label>
                            <UserSelector
                                onUserSelect={(user) => {
                                    if (user) {
                                        setRecipientEmail(user.email);
                                        setRecipientPublicKey(user.public_key);
                                        setRecipientUuid(user.uuid);
                                    } else {
                                        setRecipientEmail('');
                                        setRecipientPublicKey(null);
                                        setRecipientUuid(null);
                                    }
                                }}
                            />
                        </div>
                    )}
                </div>

                {progress > 0 && progress < 100 && (
                    <div className="progress-bar">
                        <div
                            className="progress-fill"
                            style={{ width: `${progress}%` }}
                        ></div>
                        <span className="progress-text">{Math.round(progress)}%</span>
                    </div>
                )}

                <div className="actions">
                    <button
                        type="button"
                        className="cancel-btn"
                        onClick={() => navigate('/storage/my')}
                        disabled={isLoading}
                    >
                        Отмена
                    </button>
                    <button
                        type="submit"
                        className="upload-btn"
                        onClick={handleSubmit}
                        disabled={!file || isLoading || (accessType === 'specific' && !recipientEmail)}
                    >
                        {isLoading ? (
                            <>
                                <i className="icon-spinner"></i>
                                {accessType === 'private' ? 'Шифрование и загрузка...' : 'Загрузка...'}
                            </>
                        ) : (
                            accessType === 'private' ? 'Зашифровать и загрузить' : 'Загрузить файл'
                        )}
                    </button>
                </div>
            </div>
        </div>
    );
}