import {BrowserRouter as Router, Routes, Route, useLocation, Navigate} from 'react-router-dom';
import SignIn from './pages/SignIn';
import SignUp from './pages/SignUp';
import FileStorage from './pages/FileStorage';
import MyFiles from './pages/MyFiles';
import SharedFiles from './pages/SharedFiles';
import UploadFile from './pages/UploadFile';
import React from "react";
import useAuthStore from "./stores/authStore";
import {AppWrapper} from "./components/AppWrapper";
import {FileDownloadRedirect} from "./components/FileDownloadRedirect";

function ProtectedRoute({ children }) {
    const { isAuthenticated, loading } = useAuthStore();
    const location = useLocation();

    if (loading) {
        return <div>Loading...</div>;
    }

    if (!isAuthenticated) {
        return <Navigate to="/signin" state={{ from: location }} replace />;
    }

    return children;
}

function App() {
    return (
                <AppWrapper>
                    <Router>
                        <Routes>
                            <Route path="/signin" element={<SignIn />} />
                            <Route path="/signup" element={<SignUp />} />

                            <Route path="/storage" element={<ProtectedRoute><FileStorage /></ProtectedRoute>}>
                                <Route path="my" element={<MyFiles />} />
                                <Route path="shared" element={<SharedFiles />} />
                            </Route>

                            <Route path="/upload" element={<ProtectedRoute><UploadFile /></ProtectedRoute>} />
                            <Route path="/file/:id" element={<FileDownloadRedirect />} />
                            <Route path="*" element={<SignIn />} />
                        </Routes>
                    </Router>
                </AppWrapper>
    );
}

export default App;