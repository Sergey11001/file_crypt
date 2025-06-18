import { useState, useEffect } from 'react';
import Spinner from "./Spinner";
import UserSelector from "./UserSelector";
import filesApi from "../api/files";
import usersApi from "../api/users";
import {toast} from "react-toastify";
import {useCurrentOrigin} from "../hooks/useCurrentOrigin";

export default function FileItem({ file, onDownload, onShare, onDelete, isOwner, isDownloading }) {
    const [showShareModal, setShowShareModal] = useState(false);
    const [recipientEmail, setRecipientEmail] = useState('');
    const [recipientPublicKey, setRecipientPublicKey] = useState(null);
    const [recipientUuid, setRecipientUuid] = useState(null);
    const [isSharing, setIsSharing] = useState(false);
    const [sharedUsers, setSharedUsers] = useState([]);
    const [loadingSharedUsers, setLoadingSharedUsers] = useState(false);

    const origin = useCurrentOrigin()

    useEffect(() => {
        if (showShareModal && isOwner) {
            fetchSharedUsers();
        }
    }, [isOwner, showShareModal]);

    const fetchSharedUsers = async () => {
        setLoadingSharedUsers(true);
        try {
            const response = await usersApi.getSharedUsers(file.id);
            setSharedUsers(response.data.users);
        } catch (error) {
            toast.warn("Не удалось получить пользователей c доступом к файлу");
        } finally {
            setLoadingSharedUsers(false);
        }
    };

    const handleShareSubmit = async (e) => {
        e.preventDefault();
        if (!recipientEmail) return;

        setIsSharing(true);
        try {
            await onShare(file.id, file.encryptedKey, recipientPublicKey, recipientUuid);
            setRecipientEmail('');
            setRecipientPublicKey(null);
            setRecipientUuid(null);
            await fetchSharedUsers();
        } catch (error) {
            toast.warn("Не удалось поделиться файлом");
        } finally {
            setIsSharing(false);
        }
    };

    const handleRevokeAccess = async (userId) => {
        try {
            await filesApi.deleteFileAccess(file.id,  {
                recipient_uuid: userId
            });

            await fetchSharedUsers();
        } catch (error) {
            toast.warn("Не удалось отозвать доступ.");
        }
    };

    return (
        <>
            <div className="file-item">
                <div className="file-info">
                    <span className="file-name">
                      {file.encryptedKey ? (
                          file.name
                      ) : (
                          <span className="shareable-file">
                          {file.name} (Общий доступ:{' '}
                              <span
                                  className="shared-link"
                                  onClick={() => {
                                      navigator.clipboard.writeText(`${origin}/file/${file.id}`);
                                      toast.success('Ссылка скопирована!');
                                  }}
                              >
                            {origin}/file/{file.id}
                          </span>)
                        </span>
                      )}
                    </span>
                    <div className="file-meta">
                        <span className="file-size">{file.size}</span>
                        <span className="file-date">{file.date}</span>
                    </div>
                </div>
                <div className="file-actions">
                    <button
                        className="download-btn"
                        onClick={() => onDownload()}
                        disabled={isDownloading}
                    >
                        {isDownloading ? (
                            <>
                                <Spinner size="small" color="white" />
                                Скачивание...
                            </>
                        ) : 'Скачать'}
                    </button>
                    {isOwner && (
                        <>
                            {
                                file.encryptedKey ?
                                    <button
                                        className="share-btn"
                                        onClick={() => {
                                            setShowShareModal(true)
                                        }}
                                    >
                                        Поделиться
                                    </button>
                                    : null
                            }

                            <button
                                className="delete-btn"
                                onClick={() => onDelete(file.id)}
                            >
                                Удалить
                            </button>
                        </>
                    )}
                </div>
            </div>

            {showShareModal && (
                <div className="modal-overlay">
                    <div className="modal-content">
                        <button
                            className="modal-close"
                            onClick={() => setShowShareModal(false)}
                        >
                            &times;
                        </button>

                        <h3>Управление доступом для {file.name}"</h3>

                        <form onSubmit={handleShareSubmit} className="share-form">
                            <div className="form-group recipient-section">
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
                                    fileUUID={file.id}
                                />
                            </div>

                            <div className="modal-actions">
                                <button
                                    type="submit"
                                    className="btn-confirm"
                                    disabled={isSharing || !recipientEmail}
                                >
                                    {isSharing ? (
                                        <>
                                            <Spinner size="small" color="white" />
                                            Загрузка...
                                        </>
                                    ) : 'Share'}
                                </button>
                            </div>
                        </form>

                        <div className="shared-users-list">
                            <h4>Пользователи с доступом</h4>
                            {loadingSharedUsers ? (
                                <div className="loading-shared-users">
                                    <Spinner size="medium" />
                                </div>
                            ) : sharedUsers.length > 0 ? (
                                <ul>
                                    {sharedUsers.map(user => (
                                        <li key={user.uuid} className="shared-user-item">
                                            <div>
                                                <p>{user.email}</p>
                                                <p>{user.name}</p>
                                            </div>

                                            <button
                                                onClick={() => handleRevokeAccess(user.uuid)}
                                                className="btn-revoke"
                                                title="Revoke access"
                                            >
                                                &times;
                                            </button>
                                        </li>
                                    ))}
                                </ul>
                            ) : (
                                <p className="no-users-message">Пока нет пользователей, имеющих доступ</p>
                            )}
                        </div>
                    </div>
                </div>
            )}
        </>
    );
}
