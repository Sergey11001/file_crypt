import usePrivateKeyStore from "../stores/privateKeyStore";

export default function PrivateKeyUpload() {
    const { setPrivateKey, clearAll, keyFile, setKeyFile } = usePrivateKeyStore();

    const handleFileChange = (e) => {
        const file = e.target.files[0];
        if (file && file.name.endsWith(".pem")) {
            setKeyFile(file);

            const reader = new FileReader();
            reader.onload = (event) => {
                setPrivateKey(event.target.result);
            };
            reader.readAsText(file);

        } else {
            alert("Пожалуйста, выберите файл с расширением .pem");
        }
    };

    const handleRemoveKey = () => {
        setKeyFile(null);
        clearAll()
    };

    return (
        <div className="key-upload">
            <h3>Приватный ключ</h3>
            <div className="upload-area">
                <input
                    type="file"
                    id="key-upload"
                    accept=".pem"
                    onChange={handleFileChange}
                    style={{ display: "none" }}
                />
                <label htmlFor="key-upload" className="upload-btn">
                    {keyFile ? keyFile.name : "Выберите файл ключа (.pem)"}
                </label>
                {keyFile && (
                    <button className="remove-btn" onClick={handleRemoveKey}>
                        ×
                    </button>
                )}
            </div>
        </div>
    );
}