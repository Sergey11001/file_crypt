import {useState} from 'react';
import { Link, Outlet, useLocation, useNavigate } from 'react-router-dom';
import PrivateKeyUpload from '../components/PrivateKeyUpload';
import {toast} from "react-toastify";
import {exportPrivateKey, exportPublicKey, generateKeyPair} from "../utils/crypto";
import useAuthStore from "../stores/authStore"
import usersApi from "../api/users";
import usePrivateKeyStore from "../stores/privateKeyStore";

export default function FileStorage() {
    const {user, signOut} = useAuthStore()
    const {clearAll} = usePrivateKeyStore()
    const location = useLocation();
    const navigate = useNavigate();
    const [activeTab, setActiveTab] = useState(
        location.pathname.includes('shared') ? 'shared' : 'my'
    );
    const [isRegenerating, setIsRegenerating] = useState(false);
    const [showConfirmModal, setShowConfirmModal] = useState(false);

    const handleLogout = () => {
        clearAll()
        signOut()
        navigate('/signin');
    };

    const handleRegenerate = async () => {
        setIsRegenerating(true);
        try {
            await regenerateKeys();
            toast.success('Ключи успешно сгенерированы!');
            window.location.reload();
        } catch (error) {
            toast.error(`Не удалось повторно сгенерировать ключи: ${error.message}`);
        } finally {
            setIsRegenerating(false);
            setShowConfirmModal(false);
        }
    };

    const regenerateKeys = async () => {
        const keyPair = await generateKeyPair();
        const publicKey = await exportPublicKey(keyPair.publicKey);
        const privateKeyPem = await exportPrivateKey(keyPair.privateKey);

        await usersApi.updateUserKeys({
            public_key: publicKey,
        })

        localStorage.setItem('public_key', publicKey);

        const blob = new Blob([privateKeyPem], { type: 'application/x-pem-file' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `filecrypt_${user.email}_private_key.pem`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);

        clearAll()
    }

    return (
        <div className="file-storage-container">
            <div className="file-storage-top">
                <h2>Хранилище файлов</h2>
                <button
                    onClick={handleLogout}
                    className="logout-btn"
                >
                    Выйти
                </button>
            </div>

            <div className="key-upload-section">
                <PrivateKeyUpload />
                <button
                    className="regenerate-keys-button"
                    onClick={() => setShowConfirmModal(true)}
                    disabled={isRegenerating}
                >
                    {isRegenerating ? 'Перевыпуск...' : 'Перевыпуск ключей'}
                </button>

                {showConfirmModal && (
                    <div className="modal-overlay">
                        <div className="modal-content">
                            <button
                                className="modal-close"
                                onClick={() => setShowConfirmModal(false)}
                            >
                                &times;
                            </button>

                            <h3>Подтвердите перевыпуск ключей</h3>
                            <p>Это сделает ваши текущие ключи недействительными и удалит все ваши файлы. Вы уверены, что хотите продолжить?</p>
                            <div className="modal-actions">
                                <button
                                    type="submit"
                                    className="btn-confirm"
                                    onClick={handleRegenerate}
                                    disabled={isRegenerating}
                                >
                                    {isRegenerating ? 'Обработка...' : 'Подтвердить'}
                                </button>
                            </div>
                        </div>
                    </div>
                )}
            </div>

            <div className="tabs-container">
                <Link
                    to="/storage/my"
                    className={`tab ${activeTab === 'my' ? 'active' : ''}`}
                    onClick={() => setActiveTab('my')}
                >
                    Мои файлы
                </Link>
                <Link
                    to="/storage/shared"
                    className={`tab ${activeTab === 'shared' ? 'active' : ''}`}
                    onClick={() => setActiveTab('shared')}
                >
                    Доступные файлы
                </Link>
            </div>

            <div className="add-file-btn-container">
                <Link to="/upload" className="add-file-btn">
                    + Добавить файл
                </Link>
            </div>

            {/* Контент страницы */}
            <div className="files-content">
                <Outlet />
            </div>
        </div>
    );
}