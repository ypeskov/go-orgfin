

async function getKey(key, salt) {
    const enc = new TextEncoder();
    const keyMaterial = await crypto.subtle.importKey(
        'raw',
        enc.encode(key),
        'PBKDF2',
        false,
        ['deriveKey']
    );
    return crypto.subtle.deriveKey(
        {
            name: 'PBKDF2',
            salt: salt,
            iterations: 100000,
            hash: 'SHA-256'
        },
        keyMaterial,
        {name: 'AES-GCM', length: 256},
        false,
        ['encrypt', 'decrypt']
    );
}

function base64ToArrayBuffer(base64) {
    const binary_string = window.atob(base64);
    const len = binary_string.length;
    const bytes = new Uint8Array(len);
    for (let i = 0; i < len; i++) {
        bytes[i] = binary_string.charCodeAt(i);
    }
    return bytes.buffer;
}

function arrayBufferToBase64(buffer) {
    let binary = '';
    const bytes = new Uint8Array(buffer);
    const len = bytes.byteLength;
    for (let i = 0; i < len; i++) {
        binary += String.fromCharCode(bytes[i]);
    }
    return window.btoa(binary);
}

async function encrypt(text, key) {
    const enc = new TextEncoder();
    const encoded = enc.encode(text);
    const iv = crypto.getRandomValues(new Uint8Array(12));
    const salt = crypto.getRandomValues(new Uint8Array(16));
    const cryptoKey = await getKey(key, salt);
    const ciphertext = await crypto.subtle.encrypt(
        {name: 'AES-GCM', iv: iv},
        cryptoKey,
        encoded
    );

    const encData = {
        data: arrayBufferToBase64(ciphertext),
        salt: arrayBufferToBase64(salt),
        iv: arrayBufferToBase64(iv)
    };

    return encData;
}

async function decrypt(encryptedData, key, salt, iv) {
    const cryptoKey = await getKey(key, new Uint8Array(salt));
    const decrypted = await crypto.subtle.decrypt(
        {name: 'AES-GCM', iv: new Uint8Array(iv)},
        cryptoKey,
        new Uint8Array(encryptedData)
    );
    const dec = new TextDecoder();
    return dec.decode(decrypted);
}

async function showPassword() {
    const salt = base64ToArrayBuffer(document.getElementById('password-salt').textContent);
    const iv = base64ToArrayBuffer(document.getElementById('password-iv').textContent);
    const encryptedPassword = base64ToArrayBuffer(document.getElementById('password').textContent);

    document.getElementById('decrypted-password').textContent
        = await decrypt(encryptedPassword, getMasterPassword(), salt, iv);
    document.getElementById('decrypted-password').classList.remove('hidden');
    document.getElementById('show-password-button').classList.add('hidden');
}



function getMasterPassword() {
    let masterPassword = sessionStorage.getItem('masterPassword');

    if (masterPassword === null) {
        masterPassword = prompt('Enter master password:');
        sessionStorage.setItem('masterPassword', masterPassword);
    }

    return masterPassword;
}
