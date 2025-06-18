import {toast} from "react-toastify";

export async function decryptFile(encryptedBuffer, encryptedAesKey, privateKeyPem) {
    try {
        if (!encryptedBuffer || !encryptedAesKey || !privateKeyPem) {
            throw new Error("Missing required parameters");
        }

        const symmetricKey = await decryptWithPrivateKey(encryptedAesKey, privateKeyPem);
        if (!(symmetricKey instanceof ArrayBuffer)) {
            throw new Error("Symmetric key must be an ArrayBuffer");
        }

        if (encryptedBuffer.byteLength < 12) {
            throw new Error("Encrypted data too short to contain IV");
        }

        const iv = encryptedBuffer.slice(0, 12);
        const encryptedData = encryptedBuffer.slice(12);

        const cryptoKey = await window.crypto.subtle.importKey(
            "raw",
            symmetricKey,
            { name: "AES-GCM" },
            false,
            ["decrypt"]
        );

        const decryptedData = await window.crypto.subtle.decrypt(
            {
                name: "AES-GCM",
                iv: new Uint8Array(iv),
                tagLength: 128
            },
            cryptoKey,
            encryptedData
        );

        return new Blob([decryptedData]);

    } catch (error) {
        throw new Error(`File decryption failed: ${error.message}`);
    }
}
export async function decryptWithPrivateKey(encryptedKeyBase64, privateKeyPem) {
    const privateKey = await importPrivateKey(privateKeyPem);

    const encryptedKeyBuffer = base64ToArrayBuffer(encryptedKeyBase64);

    return await window.crypto.subtle.decrypt(
        {name: "RSA-OAEP"},
        privateKey,
        encryptedKeyBuffer
    );
}

async function importPrivateKey(pem) {
    const pemHeader = "-----BEGIN PRIVATE KEY-----";
    const pemFooter = "-----END PRIVATE KEY-----";
    const pemContents = pem.replace(pemHeader, "").replace(pemFooter, "").replace(/\s+/g, "");

    const binaryDer = base64ToArrayBuffer(pemContents);

    return await window.crypto.subtle.importKey(
        "pkcs8",
        binaryDer,
        { name: "RSA-OAEP", hash: "SHA-256" },
        true,
        ["decrypt"]
    );
}

export async function generateKeyPair() {
    return await window.crypto.subtle.generateKey(
        {
            name: "RSA-OAEP",
            modulusLength: 2048,
            publicExponent: new Uint8Array([0x01, 0x00, 0x01]),
            hash: "SHA-256",
        },
        true,
        ["encrypt", "decrypt"]
    );
}

export async function exportPublicKey(publicKey) {
    const exported = await window.crypto.subtle.exportKey("spki", publicKey);
    return arrayBufferToBase64(exported);
}

export async function exportPrivateKey(privateKey) {
    const exported = await window.crypto.subtle.exportKey("pkcs8", privateKey);
    const base64 = arrayBufferToBase64(exported);
    return `-----BEGIN PRIVATE KEY-----\n${base64}\n-----END PRIVATE KEY-----`;
}

export async function generateAesKey() {
    return await window.crypto.subtle.generateKey(
        {
            name: "AES-GCM",
            length: 256,
        },
        true,
        ["encrypt", "decrypt"]
    );
}

export async function encryptFile(file, key) {
    const fileData = await file.arrayBuffer();
    const iv = window.crypto.getRandomValues(new Uint8Array(12));

    const encryptedData = await window.crypto.subtle.encrypt(
        {
            name: "AES-GCM",
            iv: iv
        },
        key,
        fileData
    );

    const result = new Uint8Array(iv.length + encryptedData.byteLength);
    result.set(iv, 0);
    result.set(new Uint8Array(encryptedData), iv.length);

    return new Blob([result], { type: file.type });
}

export async function encryptWithPublicKey(data, publicKeyPem) {
    try {
        if (data instanceof CryptoKey) {
            if (!data.extractable) {
                throw new Error("Key is not extractable");
            }

            const exportedKey = await window.crypto.subtle.exportKey("raw", data);
            const publicKey = await importPublicKey(publicKeyPem);

            const encrypted = await window.crypto.subtle.encrypt(
                { name: "RSA-OAEP" },
                publicKey,
                exportedKey
            );

            return arrayBufferToBase64(encrypted);
        } else {
            const publicKey = await importPublicKey(publicKeyPem);

            const encrypted = await window.crypto.subtle.encrypt(
                { name: "RSA-OAEP" },
                publicKey,
                data
            );

            return arrayBufferToBase64(encrypted);
        }
    } catch (error) {
        toast.error("Ошибка шифрования");
        throw new Error("Failed to encrypt with public key");
    }
}

async function importPublicKey(pem) {
    const pemHeader = "-----BEGIN PUBLIC KEY-----";
    const pemFooter = "-----END PUBLIC KEY-----";
    const pemContents = pem
        .replace(pemHeader, '')
        .replace(pemFooter, '')
        .replace(/\s+/g, '');

    const binaryDer = base64ToArrayBuffer(pemContents);

    return await window.crypto.subtle.importKey(
        "spki",
        binaryDer,
        { name: "RSA-OAEP", hash: "SHA-256" },
        true,
        ["encrypt"]
    );
}

function arrayBufferToBase64(buffer) {
    let binary = '';
    const bytes = new Uint8Array(buffer);
    for (let i = 0; i < bytes.byteLength; i++) {
        binary += String.fromCharCode(bytes[i]);
    }
    return window.btoa(binary);
}

function base64ToArrayBuffer(base64) {
    const binaryString = window.atob(base64);
    const bytes = new Uint8Array(binaryString.length);
    for (let i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i);
    }
    return bytes;
}