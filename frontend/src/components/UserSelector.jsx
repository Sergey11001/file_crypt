import { useState, useEffect } from 'react';
import usersApi from '../api/users';

export default function UserSelector({ onUserSelect, fileUUID }) {
    const [searchTerm, setSearchTerm] = useState('');
    const [users, setUsers] = useState([]);
    const [filteredUsers, setFilteredUsers] = useState([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState(null);
    const [selectedUser, setSelectedUser] = useState(null);

    useEffect(() => {
        const fetchUsers = async () => {
            setIsLoading(true);
            try {
                let response;
                if (fileUUID) {
                     response = await usersApi.getUsersForShare(fileUUID);
                }else {
                     response = await usersApi.getUsers();
                }
                setUsers(response.data.users);
                setFilteredUsers(response.data.users);
            } catch (err) {
                setError(err.message || 'Failed to load users');
            } finally {
                setIsLoading(false);
            }
        };

        fetchUsers();
    }, []);

    useEffect(() => {
        if (searchTerm.trim() === '') {
            setFilteredUsers(users);
        } else {
            const filtered = users.filter(user =>
                user.email.toLowerCase().includes(searchTerm.toLowerCase()) ||
                (user.name && user.name.toLowerCase().includes(searchTerm.toLowerCase()))
            );
            setFilteredUsers(filtered);
        }
    }, [searchTerm, users]);

    const handleUserSelect = (user) => {
        setSelectedUser(user);
        onUserSelect({
            uuid: user.uuid,
            email: user.email,
            public_key: user.public_key
        });
    };

    const clearSelection = () => {
        setSelectedUser(null);
        onUserSelect(null);
    };

    if (isLoading) return <div className="loading">Loading users...</div>;
    if (error) return <div className="error">{error}</div>;

    return (
        <div className="user-selector">
            {selectedUser ? (
                <div className="selected-user">
                    <div className="user-info">
                        <span className="user-name">{selectedUser.name}</span>
                        <span className="user-email">{selectedUser.email}</span>
                    </div>
                    <button
                        className="clear-btn"
                        onClick={clearSelection}
                    >
                        ×
                    </button>
                </div>
            ) : (
                <>
                    <input
                        type="text"
                        placeholder="Поиск по электронной почте или имени..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        className="search-input"
                    />
                    <div className="users-list">
                        {filteredUsers.length > 0 ? (
                            filteredUsers.slice(0,5).map(user => (
                                <div
                                    key={user.uuid}
                                    className="user-item"
                                    onClick={() => handleUserSelect(user)}
                                >
                                    <div className="user-info">
                                        <span className="user-name">{user.name}</span>
                                        <span className="user-email">{user.email}</span>
                                    </div>
                                </div>
                            ))
                        ) : (
                            <div className="no-results">Пользователи не найдены</div>
                        )}
                    </div>
                </>
            )}
        </div>
    );
}